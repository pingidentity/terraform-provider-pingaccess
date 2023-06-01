terraform {
  required_version = ">=1.1"
  required_providers {
    pingaccess = {
      version = "~> 0.0.1"
      source = "pingidentity/pingaccess"
    }
  }
}
provider "pingaccess" {
  username = "administrator"
  password = "2Access"
  https_host = "https://localhost:9000"
}

# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #
resource "pingaccess_proxy" "proxyExample" {
  host = "example"
  name = "example" 
  port = 1234
  password = {
    # This value will be stored into your state file 
    # and will not detect any configuration changes made in the UI
    value = "example"
  }
  requires_authentication = true
  username = "example"
}