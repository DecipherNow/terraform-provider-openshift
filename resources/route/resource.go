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

package route

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	a "github.com/openshift/api/route/v1"
	c "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	m "k8s.io/apimachinery/pkg/apis/meta/v1"
	r "k8s.io/client-go/rest"
)

// Resource returns resource for OpenShift routes.
func Resource() *schema.Resource {
	return &schema.Resource{
		Create:   create,
		Read:     read,
		Update:   update,
		Delete:   delete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
		},
	}
}

func create(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating route client")
	}

	object := m.ObjectMeta{
		Name:        data.Get("name").(string),
		Namespace:   data.Get("namespace").(string),
		Annotations: map[string]string{
			// "openshift.io/description":  data.Get("description").(string),
			// "openshift.io/display-name": data.Get("display_name").(string),
		},
	}

	route, err := client.Routes(object.Namespace).Create(&a.Route{ObjectMeta: object})
	if err != nil {
		return errors.Wrapf(err, "error creating route %s", object.Name)
	}

	data.SetId(fmt.Sprintf("%s/%s", route.Namespace, route.Name))

	return read(data, meta)
}

func delete(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating route client")
	}

	err = client.Routes(data.Get("namespace").(string)).Delete(data.Id(), &m.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deleting route %s", data.Id())
	}

	data.SetId("")

	return nil
}

func read(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	route, err := client.Routes(data.Get("namespace").(string)).Get(data.Id(), m.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "error reading route %s", data.Id())
	}

	data.SetId(fmt.Sprintf("%s/%s", route.Namespace, route.Name))
	data.Set("name", route.Name)
	data.Set("namespace", route.Namespace)

	return nil
}

func update(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	object := m.ObjectMeta{
		Name:        data.Get("name").(string),
		Annotations: map[string]string{
			// "openshift.io/description":  data.Get("description").(string),
			// "openshift.io/display-name": data.Get("display_name").(string),
		},
	}

	project, err := client.Projects().Update(&a.Project{ObjectMeta: object})
	if err != nil {
		return errors.Wrapf(err, "error reading project %s", data.Get("name").(string))
	}

	data.SetId(project.Name)

	return read(data, meta)
}
