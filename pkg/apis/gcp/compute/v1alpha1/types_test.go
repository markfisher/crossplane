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

package v1alpha1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/resource"
	"github.com/crossplaneio/crossplane/pkg/test"
)

const (
	namespace = "default"
	name      = "test-cluster"
)

var (
	c   client.Client
	ctx = context.TODO()
)

var _ resource.Managed = &GKECluster{}

func TestMain(m *testing.M) {
	t := test.NewEnv(namespace, SchemeBuilder.SchemeBuilder, test.CRDs())
	c = t.StartClient()
	t.StopAndExit(m.Run())
}

func TestGKECluster(t *testing.T) {
	key := types.NamespacedName{Name: name, Namespace: namespace}
	created := &GKECluster{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: GKEClusterSpec{
			ClusterVersion: "1.1.1",
			NumNodes:       int64(1),
			Zone:           "us-central1-a",
			MachineType:    "n1-standard-1",
			ResourceSpec: v1alpha1.ResourceSpec{
				ProviderReference: &core.ObjectReference{},
			},
		},
	}
	g := NewGomegaWithT(t)

	// Test Create
	fetched := &GKECluster{}
	g.Expect(c.Create(ctx, created)).NotTo(HaveOccurred())

	g.Expect(c.Get(ctx, key, fetched)).NotTo(HaveOccurred())
	g.Expect(fetched).To(Equal(created))

	// Test Updating the Labels
	updated := fetched.DeepCopy()
	updated.Labels = map[string]string{"hello": "world"}
	g.Expect(c.Update(ctx, updated)).NotTo(HaveOccurred())

	g.Expect(c.Get(ctx, key, fetched)).NotTo(HaveOccurred())
	g.Expect(fetched).To(Equal(updated))

	// Test Delete
	g.Expect(c.Delete(ctx, fetched)).NotTo(HaveOccurred())
	g.Expect(c.Get(ctx, key, fetched)).To(HaveOccurred())
}

func TestParseClusterSpec(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
		want *GKEClusterSpec
	}{
		{
			name: "NilProperties",
			args: nil,
			want: &GKEClusterSpec{
				ResourceSpec: v1alpha1.ResourceSpec{
					ReclaimPolicy: DefaultReclaimPolicy,
				},
				EnableIPAlias: false,
				Labels:        map[string]string{},
				NumNodes:      1,
				Scopes:        []string{},
			},
		},
		{
			name: "EmptyProperties",
			args: map[string]string{},
			want: &GKEClusterSpec{
				ResourceSpec: v1alpha1.ResourceSpec{
					ReclaimPolicy: DefaultReclaimPolicy,
				},
				EnableIPAlias: false,
				Labels:        map[string]string{},
				NumNodes:      1,
				Scopes:        []string{},
			},
		},
		{
			name: "ValidValues",
			args: map[string]string{
				"enableIPAlias":    "true",
				"machineType":      "test-machine",
				"numNodes":         "3",
				"scopes":           "foo,bar",
				"zone":             "test-zone",
				"createSubnetwork": "true",
				"clusterVersion":   "1.11",
			},
			want: &GKEClusterSpec{
				ResourceSpec: v1alpha1.ResourceSpec{
					ReclaimPolicy: DefaultReclaimPolicy,
				},
				EnableIPAlias:    true,
				CreateSubnetwork: true,
				Labels:           map[string]string{},
				MachineType:      "test-machine",
				NumNodes:         3,
				Scopes:           []string{"foo", "bar"},
				Zone:             "test-zone",
				ClusterVersion:   "1.11",
			},
		},
		{
			name: "InvalidValues",
			args: map[string]string{
				"enableIPAlias": "really",
				"machineType":   "test-machine",
				"numNodes":      "3.3",
				"scopes":        "foo,bar",
				"zone":          "test-zone",
			},
			want: &GKEClusterSpec{
				ResourceSpec: v1alpha1.ResourceSpec{
					ReclaimPolicy: DefaultReclaimPolicy,
				},
				Labels:        map[string]string{},
				EnableIPAlias: false,
				MachineType:   "test-machine",
				NumNodes:      1,
				Scopes:        []string{"foo", "bar"},
				Zone:          "test-zone",
			},
		},
		{
			name: "Defaults",
			args: map[string]string{
				"machineType": "test-machine",
				"zone":        "test-zone",
			},
			want: &GKEClusterSpec{
				ResourceSpec: v1alpha1.ResourceSpec{
					ReclaimPolicy: DefaultReclaimPolicy,
				},
				EnableIPAlias: false,
				Labels:        map[string]string{},
				MachineType:   "test-machine",
				NumNodes:      1,
				Scopes:        []string{},
				Zone:          "test-zone",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseClusterSpec(tt.args)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseClusterSpec(): -want, +got\n%s", diff)
			}
		})
	}
}

func Test_parseNodesNumber(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int64
	}{
		{name: "Empty", args: "", want: DefaultNumberOfNodes},
		{name: "Invalid", args: "foo", want: DefaultNumberOfNodes},
		{name: "0", args: "0", want: int64(0)},
		{name: "44", args: "44", want: int64(44)},
		{name: "-44", args: "-44", want: DefaultNumberOfNodes},
		{name: "1.2", args: "1.2", want: DefaultNumberOfNodes},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNodesNumber(tt.args); got != tt.want {
				t.Errorf("parseNodesNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
