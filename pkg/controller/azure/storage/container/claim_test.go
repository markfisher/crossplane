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

package container

import (
	"context"
	"testing"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/crossplaneio/crossplane/pkg/apis/azure/storage/v1alpha1"
	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	storagev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/storage/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/resource"
	"github.com/crossplaneio/crossplane/pkg/test"
)

var _ resource.ManagedConfigurator = resource.ManagedConfiguratorFn(ConfigureContainer)

func TestConfigureContainer(t *testing.T) {
	type args struct {
		ctx context.Context
		cm  resource.Claim
		cs  *corev1alpha1.ResourceClass
		mg  resource.Managed
	}

	type want struct {
		mg  resource.Managed
		err error
	}

	claimUID := types.UID("definitely-a-uuid")
	providerName := "coolprovider"
	bucketName := "coolbucket"
	bucketPrivate := storagev1alpha1.ACLPrivate

	cases := map[string]struct {
		args args
		want want
	}{
		"Successful": {
			args: args{
				cm: &storagev1alpha1.Bucket{
					ObjectMeta: metav1.ObjectMeta{UID: claimUID},
					Spec: storagev1alpha1.BucketSpec{
						Name:          bucketName,
						PredefinedACL: &bucketPrivate,
					},
				},
				cs: &corev1alpha1.ResourceClass{
					ProviderReference: &corev1.ObjectReference{Name: providerName},
					ReclaimPolicy:     corev1alpha1.ReclaimDelete,
				},
				mg: &v1alpha1.Container{},
			},
			want: want{
				mg: &v1alpha1.Container{
					Spec: v1alpha1.ContainerSpec{
						ReclaimPolicy:    corev1alpha1.ReclaimDelete,
						AccountReference: corev1.LocalObjectReference{Name: providerName},
						Metadata:         azblob.Metadata{},
					},
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := ConfigureContainer(tc.args.ctx, tc.args.cm, tc.args.cs, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("ConfigureContainer(...): -want error, +got error:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.mg, tc.args.mg, test.EquateConditions()); diff != "" {
				t.Errorf("ConfigureContainer(...) Managed: -want, +got:\n%s", diff)
			}
		})
	}
}
