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

	a "github.com/openshift/api/project/v1"
	c "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"

	m "k8s.io/apimachinery/pkg/apis/meta/v1"
	r "k8s.io/client-go/rest"
)

// Resource returns resource for OpenShift projects.
func Resource() *schema.Resource {
	return &schema.Resource{
		Create:   create,
		Read:     read,
		Update:   update,
		Delete:   delete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func create(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	object := m.ObjectMeta{
		Name: data.Get("name").(string),
		Annotations: map[string]string{
			"openshift.io/description":  data.Get("description").(string),
			"openshift.io/display-name": data.Get("display_name").(string),
		},
	}

	project, err := client.Projects().Create(&a.Project{ObjectMeta: object})
	if err != nil {
		return errors.Wrapf(err, "error creating project %s", data.Get("name").(string))
	}

	data.SetId(string(project.Name))

	return read(data, meta)
}

func delete(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	err = client.Projects().Delete(data.Id(), &m.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deleting project %s", data.Get("name").(string))
	}

	data.SetId("")

	return nil
}

func read(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	project, err := client.Projects().Get(data.Id(), m.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "error reading project %s", data.Get("name").(string))
	}

	data.SetId(project.Name)
	data.Set("name", project.Name)
	data.Set("description", project.Annotations["openshift.io/description"])
	data.Set("display_name", project.Annotations["openshift.io/display-name"])

	return nil
}

func update(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating project client")
	}

	object := m.ObjectMeta{
		Name: data.Get("name").(string),
		Annotations: map[string]string{
			"openshift.io/description":  data.Get("description").(string),
			"openshift.io/display-name": data.Get("display_name").(string),
		},
	}

	project, err := client.Projects().Update(&a.Project{ObjectMeta: object})
	if err != nil {
		return errors.Wrapf(err, "error reading project %s", data.Get("name").(string))
	}

	data.SetId(project.Name)

	return read(data, meta)
}
