terraform {
  required_version = ">=1.1"
  required_providers {
    pingaccess = {
      source = "pingidentity/pingaccess"
    }
  }
}

provider "pingaccess" {
  username = "administrator"
  password = "2Access"
  https_host = "https://localhost:9000"
}
resource "pingaccess_access_token_validator" "test2" {
  classname = "com.pingidentity.pa.accesstokenvalidators.JwksEndpoint"
  name       = "example"
  configuration = {
    path                 = "/example"
  }
}