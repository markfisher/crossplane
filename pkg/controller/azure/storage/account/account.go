/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package account

import (
	"context"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/crossplaneio/crossplane/pkg/apis/azure/storage/v1alpha1"
	azurev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/azure/v1alpha1"
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/clients/azure"
	azurestorage "github.com/crossplaneio/crossplane/pkg/clients/azure/storage"
	"github.com/crossplaneio/crossplane/pkg/logging"
	"github.com/crossplaneio/crossplane/pkg/meta"
	"github.com/crossplaneio/crossplane/pkg/resource"
)

const (
	controllerName = "account.storage.azure.crossplane.io"
	finalizer      = "finalizer." + controllerName

	reconcileTimeout      = 2 * time.Minute
	requeueAfterOnSuccess = 1 * time.Minute
	requeueAfterOnWait    = 30 * time.Second
)

var (
	resultRequeue    = reconcile.Result{Requeue: true}
	requeueOnSuccess = reconcile.Result{RequeueAfter: requeueAfterOnSuccess}
	requeueOnWait    = reconcile.Result{RequeueAfter: requeueAfterOnWait}

	log = logging.Logger.WithName("controller." + controllerName)
)

// Reconciler reconciles an Azure storage account
type Reconciler struct {
	client.Client
	syncdeleterMaker
}

// Add creates a newSyncdeleter Controller and adds it to the Manager with default RBAC.
// The Manager will set fields on the Controller and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r := &Reconciler{
		Client:           mgr.GetClient(),
		syncdeleterMaker: &accountSyncdeleterMaker{mgr.GetClient()},
	}

	// Create a newSyncdeleter controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to storage account
	if err := c.Watch(&source.Kind{Type: &v1alpha1.Account{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to account Secret
	return c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Account{},
	})
}

// Reconcile reads that state of the cluster for a Provider acct and makes changes based on the state read
// and what is in the Provider.Spec
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.V(logging.Debug).Info("reconciling", "kind", v1alpha1.AccountKindAPIVersion, "request", request)

	ctx, cancel := context.WithTimeout(context.Background(), reconcileTimeout)
	defer cancel()

	b := &v1alpha1.Account{}
	if err := r.Get(ctx, request.NamespacedName, b); err != nil {
		if kerrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	bh, err := r.newSyncdeleter(ctx, b)
	if err != nil {
		b.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return resultRequeue, r.Status().Update(ctx, b)
	}

	// Check for deletion
	if b.DeletionTimestamp != nil {
		return bh.delete(ctx)
	}

	return bh.sync(ctx)
}

type syncdeleterMaker interface {
	newSyncdeleter(context.Context, *v1alpha1.Account) (syncdeleter, error)
}

type accountSyncdeleterMaker struct {
	client.Client
}

func (m *accountSyncdeleterMaker) newSyncdeleter(ctx context.Context, b *v1alpha1.Account) (syncdeleter, error) {
	p := &azurev1alpha1.Provider{}
	n := meta.NamespacedNameOf(b.Spec.ProviderReference)
	if err := m.Get(ctx, n, p); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider %s", n)
	}

	s := &corev1.Secret{}
	n = types.NamespacedName{Namespace: p.GetNamespace(), Name: p.Spec.Secret.Name}
	if err := m.Get(ctx, n, s); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider's secret %s", n)
	}

	storageClient, err := azurestorage.NewStorageAccountClient(s.Data[p.Spec.Secret.Key])
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create storageClient from json")
	}

	return newAccountSyncDeleter(
		azurestorage.NewAccountHandle(storageClient, b.Spec.ResourceGroupName, b.Spec.StorageAccountName),
		m.Client, b), nil
}

type deleter interface {
	delete(context.Context) (reconcile.Result, error)
}

type syncer interface {
	sync(context.Context) (reconcile.Result, error)
}

type creater interface {
	create(context.Context) (reconcile.Result, error)
}

type updater interface {
	update(context.Context, *storage.Account) (reconcile.Result, error)
}

type syncbacker interface {
	syncback(context.Context, *storage.Account) (reconcile.Result, error)
}

type secretupdater interface {
	updatesecret(ctx context.Context, acct *storage.Account) error
}

type syncdeleter interface {
	deleter
	syncer
}

type accountSyncDeleter struct {
	createupdater
	azurestorage.AccountOperations
	kube client.Client
	acct *v1alpha1.Account
}

func newAccountSyncDeleter(ao azurestorage.AccountOperations, kube client.Client, b *v1alpha1.Account) *accountSyncDeleter {
	return &accountSyncDeleter{
		createupdater:     newAccountCreateUpdater(ao, kube, b),
		AccountOperations: ao,
		kube:              kube,
		acct:              b,
	}
}

func (asd *accountSyncDeleter) delete(ctx context.Context) (reconcile.Result, error) {
	asd.acct.Status.SetConditions(corev1alpha1.Deleting())
	if asd.acct.Spec.ReclaimPolicy == corev1alpha1.ReclaimDelete {
		if err := asd.Delete(ctx); err != nil && !azure.IsNotFound(err) {
			asd.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
			return resultRequeue, asd.kube.Status().Update(ctx, asd.acct)
		}
	}

	// NOTE(negz): We don't update the conditioned status here because assuming
	// no other finalizers need to be cleaned up the object should cease to
	// exist after we update it.
	meta.RemoveFinalizer(asd.acct, finalizer)
	return reconcile.Result{}, asd.kube.Update(ctx, asd.acct)
}

const uidTag = "UID"

// sync - synchronizes the state of the storage account resource with the state of the
// account Kubernetes acct
func (asd *accountSyncDeleter) sync(ctx context.Context) (reconcile.Result, error) {
	account, err := asd.Get(ctx)
	if err != nil && !azure.IsNotFound(err) {
		asd.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return resultRequeue, asd.kube.Status().Update(ctx, asd.acct)
	}

	if account == nil {
		return asd.create(ctx)
	}

	// for existing account check UID tag
	if uid := to.String(account.Tags[uidTag]); uid != "" && uid != string(asd.acct.GetUID()) {
		err := errors.Errorf("storage account: %s already exists and owned by: %s", to.String(account.Name), uid)
		asd.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return reconcile.Result{}, asd.kube.Status().Update(ctx, asd.acct)
	}

	return asd.update(ctx, account)
}

// createupdater interface defining create and update operations on/for storage account resource
type createupdater interface {
	creater
	updater
}

// accountCreateUpdater implementation of createupdater interface
type accountCreateUpdater struct {
	syncbacker
	azurestorage.AccountOperations
	kube      client.Client
	acct      *v1alpha1.Account
	projectID string
}

// newAccountCreateUpdater new instance of accountCreateUpdater
func newAccountCreateUpdater(ao azurestorage.AccountOperations, kube client.Client, acct *v1alpha1.Account) *accountCreateUpdater {
	return &accountCreateUpdater{
		syncbacker:        newAccountSyncBacker(ao, kube, acct),
		AccountOperations: ao,
		kube:              kube,
		acct:              acct,
	}
}

// create new storage account resource and save changes back to account specs
func (acu *accountCreateUpdater) create(ctx context.Context) (reconcile.Result, error) {
	acu.acct.Status.SetConditions(corev1alpha1.Creating())
	meta.AddFinalizer(acu.acct, finalizer)

	// Set UID to the account storage spec
	// TODO(illya) - this eventually needs to be in Defaulter Mutating web hook
	if tags := acu.acct.Spec.StorageAccountSpec.Tags; tags == nil {
		acu.acct.Spec.StorageAccountSpec.Tags = make(map[string]string)
	}
	acu.acct.Spec.StorageAccountSpec.Tags[uidTag] = string(acu.acct.GetUID())

	accountSpec := v1alpha1.ToStorageAccountCreate(acu.acct.Spec.StorageAccountSpec)

	a, err := acu.Create(ctx, accountSpec)
	if err != nil {
		acu.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return resultRequeue, acu.kube.Status().Update(ctx, acu.acct)
	}

	return acu.syncback(ctx, a)
}

// update storage account resource if needed
func (acu *accountCreateUpdater) update(ctx context.Context, account *storage.Account) (reconcile.Result, error) {
	if account.ProvisioningState == storage.Succeeded {
		acu.acct.Status.SetConditions(corev1alpha1.Available())
		resource.SetBindable(acu.acct)

		current := v1alpha1.NewStorageAccountSpec(account)
		if reflect.DeepEqual(current, acu.acct.Spec.StorageAccountSpec) {
			acu.acct.Status.SetConditions(corev1alpha1.ReconcileSuccess())
			return requeueOnSuccess, acu.kube.Status().Update(ctx, acu.acct)
		}

		a, err := acu.Update(ctx, v1alpha1.ToStorageAccountUpdate(acu.acct.Spec.StorageAccountSpec))
		if err != nil {
			acu.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
			return resultRequeue, acu.kube.Status().Update(ctx, acu.acct)
		}
		account = a
	}

	return acu.syncback(ctx, account)
}

type accountSyncbacker struct {
	secretupdater
	acct *v1alpha1.Account
	kube client.Client
}

func newAccountSyncBacker(ao azurestorage.AccountOperations, kube client.Client, acct *v1alpha1.Account) *accountSyncbacker {
	return &accountSyncbacker{
		secretupdater: newAccountSecretUpdater(ao, kube, acct),
		kube:          kube,
		acct:          acct,
	}
}

func (asb *accountSyncbacker) syncback(ctx context.Context, acct *storage.Account) (reconcile.Result, error) {
	asb.acct.Spec.StorageAccountSpec = v1alpha1.NewStorageAccountSpec(acct)
	if err := asb.kube.Update(ctx, asb.acct); err != nil {
		return resultRequeue, err
	}

	asb.acct.Status.StorageAccountStatus = v1alpha1.NewStorageAccountStatus(acct)

	if acct.ProvisioningState != storage.Succeeded {
		asb.acct.Status.SetConditions(corev1alpha1.ReconcileSuccess())
		return requeueOnWait, asb.kube.Status().Update(ctx, asb.acct)
	}

	if err := asb.updatesecret(ctx, acct); err != nil {
		asb.acct.Status.SetConditions(corev1alpha1.ReconcileError(err))
		return resultRequeue, asb.kube.Status().Update(ctx, asb.acct)
	}

	asb.acct.Status.SetConditions(corev1alpha1.ReconcileSuccess())
	return requeueOnSuccess, asb.kube.Status().Update(ctx, asb.acct)
}

type accountSecretUpdater struct {
	azurestorage.AccountOperations
	acct *v1alpha1.Account
	kube client.Client
}

func newAccountSecretUpdater(ao azurestorage.AccountOperations, kube client.Client, acct *v1alpha1.Account) *accountSecretUpdater {
	return &accountSecretUpdater{
		AccountOperations: ao,
		acct:              acct,
		kube:              kube,
	}
}

func (asu *accountSecretUpdater) updatesecret(ctx context.Context, acct *storage.Account) error {
	secret := resource.ConnectionSecretFor(asu.acct, v1alpha1.AccountGroupVersionKind)
	key := types.NamespacedName{Namespace: secret.Namespace, Name: secret.Name}

	if acct.PrimaryEndpoints != nil {
		secret.Data[corev1alpha1.ResourceCredentialsSecretEndpointKey] = []byte(to.String(acct.PrimaryEndpoints.Blob))
	}

	keys, err := asu.ListKeys(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to list account keys")
	}
	if len(keys) == 0 {
		return errors.New("account keys are empty")
	}

	secret.Data[corev1alpha1.ResourceCredentialsSecretUserKey] = []byte(asu.acct.Spec.StorageAccountName)
	secret.Data[corev1alpha1.ResourceCredentialsSecretPasswordKey] = []byte(to.String(keys[0].Value))

	if err := asu.kube.Create(ctx, secret); err != nil {
		if kerrors.IsAlreadyExists(err) {
			return errors.Wrapf(asu.kube.Update(ctx, secret), "failed to update secret: %s", key)
		}
		return errors.Wrapf(err, "failed to create secret: %s", key)
	}

	return nil
}
