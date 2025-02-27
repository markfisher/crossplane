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

package gcp

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/crossplaneio/crossplane/pkg/controller/gcp/cache"
	"github.com/crossplaneio/crossplane/pkg/controller/gcp/compute"
	"github.com/crossplaneio/crossplane/pkg/controller/gcp/database"
	"github.com/crossplaneio/crossplane/pkg/controller/gcp/storage"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs,
		cache.Add,
		cache.AddClaim,
		compute.Add,
		compute.AddClaim,
		database.Add,
		database.AddPostgreSQLClaim,
		database.AddMySQLClaim,
		storage.Add,
		storage.AddClaim,
	)
}

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
