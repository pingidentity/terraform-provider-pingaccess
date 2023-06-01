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

resource "pingaccess_applications" "applicationExample" {
	agent_id = 0
  context_root = "/root"
  default_auth_type = "Web"
	name = "example"
  site_id = 1
	spa_support_enabled = false
	virtual_host_ids = [1,2]
}