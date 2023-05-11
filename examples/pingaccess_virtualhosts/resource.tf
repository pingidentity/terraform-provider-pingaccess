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

resource "pingaccess_virtualhosts" "virtualhostExample" {
	host = "test"	
	port = 1234
}