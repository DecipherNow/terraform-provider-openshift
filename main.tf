terraform {}

provider "openshift" {}

resource "openshift_project" "openshift" {
  description  = "Terraform created project"
  display_name = "Terraform Project"
  name         = "terraform-project"
}

# resource "openshift_imagestream" "nginx" {
#   project = "terraform-project"
#   name    = "nginx"


#   spec = {
#     docker_image_repository = "docker.io/nginx"


#     lookup_policy = {
#       local = false
#     }
#   }
# }

