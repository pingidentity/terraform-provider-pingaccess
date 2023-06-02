terraform {
  required_version = ">=1.1"
  required_providers {
    pingaccess = {
      version = "~>0.0.1"
      source = "pingidentity/pingaccess"
    }
  }
}

provider "pingaccess" {
  username = "administrator"
  password = "2Access"
  https_host = "https://localhost:9000"
}
# this resource does not support import
resource "pingaccess_certificates" "example" {
  alias     = "test"
  # this property needs to contain base64 encode value of your pem certificate
  file_data = ""
}