
# variable "bastion" {}
# variable "bastion_default" {
#   type = object({
#     ssh_private_key = string
#     ssh_public_key = string
#     host = string
#     user = string
#     port = number
#   }) 
#   default = {
#     ssh_private_key = "~/.ssh/id_rsa"
#     ssh_public_key = "~/.ssh/id_rsa.pub"
#     host = ""
#     user = ""
#     port = 22
#   }
#   description = "Default values for using a ssh bastion"
# }

# variable "rke_name" {
#   type = string
# }

# variable "domain_name" {
#   type = string
# }

# variable "api_domain" {
#   type = string
# }


variable "maas" {}
variable "rancher" {}
variable "harvester" {}

variable "harvester_bootstrap_passwd" {}
# variable "rancher_bootstrap" {}

# variable "ca_cert_path" {
#   type = string
#   description = "CA certificate path from Terraform root"
# }

# variable "ca_key_path" {
#   type = string
#   description = "CA key path from Terraform root"
# }