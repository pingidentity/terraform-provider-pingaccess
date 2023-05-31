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

resource "pingaccess_high_availability_profile" "highAvailabilityProfileExample" {
  classname = "com.pingidentity.pa.ha.availability.ondemand.OnDemandAvailabilityPlugin"
  name       = "example"
  configuration = {
    failed_retry_timeout = 60
    failure_http_status_codes = ["208","209"]
  }
}