package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingaccess/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingaccess/internal/provider"
)

const trustedCertificateGroupId = "3"

// Attributes to test with. Add optional properties to test here if desired.
type trustedCertificateGroupResourceModel struct {
	id                int64
	stateId           string
	useJavaTrustStore bool
	ocsp              bool
	name              string
}

func TestAccTrustedCertificateGroup(t *testing.T) {
	resourceName := "myTrustedCertificateGroup"
	initialResourceModel := trustedCertificateGroupResourceModel{
		id:                3,
		stateId:           "3",
		useJavaTrustStore: true,
		ocsp:              true,
		name:              "example",
	}
	updatedResourceModel := trustedCertificateGroupResourceModel{
		id:                3,
		stateId:           "3",
		useJavaTrustStore: true,
		ocsp:              false,
		name:              "example2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTrustedCertificateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedCertificateGroup(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTrustedCertificateGroupAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccTrustedCertificateGroup(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTrustedCertificateGroupAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccTrustedCertificateGroup(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_trusted_certificate_groups." + resourceName,
				ImportStateId:           trustedCertificateGroupId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccTrustedCertificateGroup(resourceName string, resourceModel trustedCertificateGroupResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_trusted_certificate_groups" "%[1]s" {
  id                   = %[2]d
  name                 = "%[3]s"
  use_java_trust_store = true
  revocation_checking = {
    ocsp = %[4]t
  }
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.ocsp,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedTrustedCertificateGroupAttributes(config trustedCertificateGroupResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "TrustedCertificateGroup"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TrustedCertificateGroupsApi.GetTrustedCertificateGroup(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, &config.stateId, "ocsp",
			config.ocsp, *response.RevocationChecking.Ocsp)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTrustedCertificateGroupDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.TrustedCertificateGroupsApi.GetTrustedCertificateGroup(ctx, trustedCertificateGroupId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("TrustedCertificateGroup", trustedCertificateGroupId)
	}
	return nil
}
