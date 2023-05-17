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

resource "pingaccess_third_party_services" "thirdPartyServiceExample" {
	name = "example"
  availability_profile_id = 1
	targets = ["localhost:80","localhost:443"]
}