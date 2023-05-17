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

resource "pingaccess_authn_req_lists" "authnReqListExample" {
  authn_reqs = ["example1","example2"]	
	name = "example"
}