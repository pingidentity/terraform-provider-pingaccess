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

resource "pingaccess_hsm_providers" "hsmProviderExample" {
  classname = "com.pingidentity.pa.hsm.pkcs11.plugin.PKCS11HsmProvider"
  name       = "example"
  configuration = {
    slot_id = "2"
    password = "example"
    library = "example"
  }
}