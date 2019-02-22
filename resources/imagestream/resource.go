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

package imagestream

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	a "github.com/openshift/api/image/v1"
	c "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"

	m "k8s.io/apimachinery/pkg/apis/meta/v1"
	r "k8s.io/client-go/rest"
)

// Resource returns resource for OpenShift imagestream.
func Resource() *schema.Resource {
	return &schema.Resource{
		Create:   create,
		Read:     read,
		Update:   update,
		Delete:   delete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"docker_image_repository": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"lookup_policy": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: false,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func create(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating image client")
	}

	object := m.ObjectMeta{
		Name: data.Get("name").(string),
	}

	lookupPolicy := a.ImageLookupPolicy{
		Local: data.Get("spec.0.lookup_policy.0.local").(bool),
	}

	spec := a.ImageStreamSpec{
		DockerImageRepository: data.Get("spec.0.docker_image_repository").(string),
		LookupPolicy:          lookupPolicy,
	}

	stream := a.ImageStream{
		ObjectMeta: object,
		Spec:       spec,
	}

	imagestream, err := client.ImageStreams(data.Get("project").(string)).Create(&stream)
	if err != nil {
		return errors.Wrapf(err, "error creating imagestream %s", data.Get("name").(string))
	}

	data.SetId(string(imagestream.Name))

	return read(data, meta)
}

func read(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating imagestream client")
	}

	imagestream, err := client.ImageStreams(data.Get("project").(string)).Get(data.Id(), m.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "error reading imagestream %s", data.Get("name").(string))
	}

	data.SetId(imagestream.Name)
	data.Set("project", imagestream.Namespace)
	data.Set("name", imagestream.Name)
	data.Set("docker_image_repository", imagestream.Spec.DockerImageRepository)
	data.Set("lookup_policy", imagestream.Spec.LookupPolicy.Local)

	return nil
}

func delete(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating imagestream client")
	}

	err = client.ImageStreams(data.Get("project").(string)).Delete(data.Id(), &m.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deleting imagestream %s", data.Get("name").(string))
	}

	data.SetId("")

	return nil
}

func update(data *schema.ResourceData, meta interface{}) error {
	client, err := c.NewForConfig(meta.(*r.Config))
	if err != nil {
		return errors.Wrapf(err, "error raised creating image client")
	}

	object := m.ObjectMeta{
		Name: data.Get("name").(string),
	}

	lookupPolicy := a.ImageLookupPolicy{
		Local: data.Get("spec.0.lookup_policy.0.local").(bool),
	}

	spec := a.ImageStreamSpec{
		DockerImageRepository: data.Get("spec.0.docker_image_repository").(string),
		LookupPolicy:          lookupPolicy,
	}

	stream := a.ImageStream{
		ObjectMeta: object,
		Spec:       spec,
	}

	imagestream, err := client.ImageStreams(data.Get("project").(string)).Create(&stream)
	if err != nil {
		return errors.Wrapf(err, "error creating imagestream %s", data.Get("name").(string))
	}

	data.SetId(string(imagestream.Name))

	return read(data, meta)
}
