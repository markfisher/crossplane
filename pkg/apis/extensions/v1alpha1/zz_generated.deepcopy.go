// +build !ignore_autogenerated

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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppMetadataSpec) DeepCopyInto(out *AppMetadataSpec) {
	*out = *in
	if in.Icons != nil {
		in, out := &in.Icons, &out.Icons
		*out = make([]IconSpec, len(*in))
		copy(*out, *in)
	}
	if in.Maintainers != nil {
		in, out := &in.Maintainers, &out.Maintainers
		*out = make([]ContributorSpec, len(*in))
		copy(*out, *in)
	}
	if in.Owners != nil {
		in, out := &in.Owners, &out.Owners
		*out = make([]ContributorSpec, len(*in))
		copy(*out, *in)
	}
	if in.Keywords != nil {
		in, out := &in.Keywords, &out.Keywords
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]LinkSpec, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppMetadataSpec.
func (in *AppMetadataSpec) DeepCopy() *AppMetadataSpec {
	if in == nil {
		return nil
	}
	out := new(AppMetadataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CRDList) DeepCopyInto(out *CRDList) {
	*out = *in
	if in.Owned != nil {
		in, out := &in.Owned, &out.Owned
		*out = make([]v1.TypeMeta, len(*in))
		copy(*out, *in)
	}
	if in.DependsOn != nil {
		in, out := &in.DependsOn, &out.DependsOn
		*out = make([]v1.TypeMeta, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CRDList.
func (in *CRDList) DeepCopy() *CRDList {
	if in == nil {
		return nil
	}
	out := new(CRDList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContributorSpec) DeepCopyInto(out *ContributorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContributorSpec.
func (in *ContributorSpec) DeepCopy() *ContributorSpec {
	if in == nil {
		return nil
	}
	out := new(ContributorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerDeployment) DeepCopyInto(out *ControllerDeployment) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerDeployment.
func (in *ControllerDeployment) DeepCopy() *ControllerDeployment {
	if in == nil {
		return nil
	}
	out := new(ControllerDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerJob) DeepCopyInto(out *ControllerJob) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerJob.
func (in *ControllerJob) DeepCopy() *ControllerJob {
	if in == nil {
		return nil
	}
	out := new(ControllerJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerSpec) DeepCopyInto(out *ControllerSpec) {
	*out = *in
	if in.Deployment != nil {
		in, out := &in.Deployment, &out.Deployment
		*out = new(ControllerDeployment)
		(*in).DeepCopyInto(*out)
	}
	if in.Job != nil {
		in, out := &in.Job, &out.Job
		*out = new(ControllerJob)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerSpec.
func (in *ControllerSpec) DeepCopy() *ControllerSpec {
	if in == nil {
		return nil
	}
	out := new(ControllerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Extension) DeepCopyInto(out *Extension) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extension.
func (in *Extension) DeepCopy() *Extension {
	if in == nil {
		return nil
	}
	out := new(Extension)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Extension) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionList) DeepCopyInto(out *ExtensionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Extension, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionList.
func (in *ExtensionList) DeepCopy() *ExtensionList {
	if in == nil {
		return nil
	}
	out := new(ExtensionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExtensionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequest) DeepCopyInto(out *ExtensionRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequest.
func (in *ExtensionRequest) DeepCopy() *ExtensionRequest {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExtensionRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequestList) DeepCopyInto(out *ExtensionRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ExtensionRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequestList.
func (in *ExtensionRequestList) DeepCopy() *ExtensionRequestList {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ExtensionRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequestSpec) DeepCopyInto(out *ExtensionRequestSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequestSpec.
func (in *ExtensionRequestSpec) DeepCopy() *ExtensionRequestSpec {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionRequestStatus) DeepCopyInto(out *ExtensionRequestStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.InstallJob != nil {
		in, out := &in.InstallJob, &out.InstallJob
		*out = new(corev1.ObjectReference)
		**out = **in
	}
	if in.ExtensionRecord != nil {
		in, out := &in.ExtensionRecord, &out.ExtensionRecord
		*out = new(corev1.ObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionRequestStatus.
func (in *ExtensionRequestStatus) DeepCopy() *ExtensionRequestStatus {
	if in == nil {
		return nil
	}
	out := new(ExtensionRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionSpec) DeepCopyInto(out *ExtensionSpec) {
	*out = *in
	in.AppMetadataSpec.DeepCopyInto(&out.AppMetadataSpec)
	in.CRDs.DeepCopyInto(&out.CRDs)
	in.Controller.DeepCopyInto(&out.Controller)
	in.Permissions.DeepCopyInto(&out.Permissions)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionSpec.
func (in *ExtensionSpec) DeepCopy() *ExtensionSpec {
	if in == nil {
		return nil
	}
	out := new(ExtensionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtensionStatus) DeepCopyInto(out *ExtensionStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.ControllerRef != nil {
		in, out := &in.ControllerRef, &out.ControllerRef
		*out = new(corev1.ObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtensionStatus.
func (in *ExtensionStatus) DeepCopy() *ExtensionStatus {
	if in == nil {
		return nil
	}
	out := new(ExtensionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IconSpec) DeepCopyInto(out *IconSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IconSpec.
func (in *IconSpec) DeepCopy() *IconSpec {
	if in == nil {
		return nil
	}
	out := new(IconSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkSpec) DeepCopyInto(out *LinkSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkSpec.
func (in *LinkSpec) DeepCopy() *LinkSpec {
	if in == nil {
		return nil
	}
	out := new(LinkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PermissionsSpec) DeepCopyInto(out *PermissionsSpec) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]rbacv1.PolicyRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PermissionsSpec.
func (in *PermissionsSpec) DeepCopy() *PermissionsSpec {
	if in == nil {
		return nil
	}
	out := new(PermissionsSpec)
	in.DeepCopyInto(out)
	return out
}
