terraform {
  required_providers {
    pingaccess = {
      source = "pingidentity.com/terraform/pingaccess"
    }
  }
}

provider "pingaccess" {}

data "AccessTokenValidator" "example" {}