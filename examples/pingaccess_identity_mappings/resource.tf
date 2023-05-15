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

# resource "pingaccess_identity_mappings" "webAccessTokenIdentityMappingExample" {
# 	classname = "com.pingidentity.pa.identitymappings.WebSessionAccessTokenIdentityMapping"	
# 	configuration = {}
# 	name = "webAccessTokenExample"
# }

# resource "pingaccess_identity_mappings" "jwtIdentityMappingExample" {
# 	classname = "com.pingidentity.pa.identitymappings.JwtIdentityMapping"	
# 	name = "jwtExample"
# 	configuration = {
#     map_as_bearer_token = false
#     header_name = "header_name"
#     audience = "audience"
#     exclusion_list = false
#     exclusion_list_attributes = ["active"]
#     exclusion_list_subject = "auth_time"
#     attribute_mappings = [
#       {
#         subject = true,
#         user_attribute_name = "active",
#         jwt_claim_name = "jwt_claim_name"
#       },
#       {
#         subject = false
#         user_attribute_name = "pi\\.pa\\.attr_exp"
#         jwt_claim_name = "pa_attr_exp"
#       },
#       {
#         subject = false
#         user_attribute_name = "auth_time"
#         jwt_claim_name = "auth_time"
#       }
#     ]
#     cache_jwt = true
#     client_certificate_jwt_claim_name = "cert_chain_jwt_claim"
#     max_depth = 1
#   }
# }

# resource "pingaccess_identity_mappings" "headerIdentityMappingExample" {
# 	classname = "com.pingidentity.pa.identitymappings.HeaderIdentityMapping"	
# 	configuration = {
#     header_name_prefix = "prefix-",
#     exclusion_list = true,
#     exclusion_list_attributes = [
#       "active"
#     ],
#     exclusion_list_subject = "amr",
#     attribute_header_mappings = [
#       {
#         subject = true,
#         attribute_name = "realm",
#         header_name = "realm"
#       }
#     ],
#     header_client_certificate_mappings = [
#       {
#         header_name = "header"
#       },
#       {
#         header_name = "tehe"
#       }
#     ]
#   }
# 	name = "example"
# }