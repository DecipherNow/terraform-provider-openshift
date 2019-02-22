// Copyright 2019 Decipher Technology Studios LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	a "github.com/openshift/api/route/v1"
	c "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	m "k8s.io/apimachinery/pkg/apis/meta/v1"
	r "k8s.io/client-go/rest"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: create,
		Read: read,
		Update: update,
		Delete: delete, 
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "A unique name for the route within the project.",
			},
			"host": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Public hostname for the route. If not specified, a hostname is generated.",
			},
			"path": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: false,
				Description: "Path that the router watches to route traffic to the service.",
			},
			"service": {
				Type: schema.TypeString,
				Required: true,
				ForceNew, false,
				Description: "Service to route to.",
			},
			"port": {
				Type: schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Description: "Target port for traffic.",
			},
			"tls": {
				Type: schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Description: "Currently, routes can either be configured for TLS passthrough or left unsecured",
			},
			"insecure_redirect": {
				Type: schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Description: "Boolean value to determine if unsecured traffic should be redirected to TLS",
				Default: false,
			},
		}
	}
}