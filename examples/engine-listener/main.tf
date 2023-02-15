terraform {
  required_providers {
    pingdirectory = {
      source = "pingidentity/pingaccess"
    }
  }
}

provider "pingdirectory" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
}

resource "pingaccess_engine_listener" "example" {
  name   = "example"
  port   = 443
  secure = true
}