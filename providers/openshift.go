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

package providers

import (
	"os"
	"path/filepath"

	"github.com/deciphernow/terraform-provider-openshift/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// OpenShift returns the provider for OpenShift.
func OpenShift() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: configure,
		ResourcesMap:  resources.Resources(),
		Schema: map[string]*schema.Schema{
			"mode": {
				Default:     "kubeconfig",
				Description: "The configuration mode of the provider.",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func configure(data *schema.ResourceData) (interface{}, error) {
	var err error
	var config *rest.Config

	switch data.Get("mode").(string) {
	case "cluster":
		config, err = rest.InClusterConfig()
	case "kubeconfig":
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}

	if err != nil {
		return nil, errors.Wrap(err, "error raised creating configuration")
	}

	return config, nil
}
