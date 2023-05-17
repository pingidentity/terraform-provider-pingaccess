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

const thirdPartyServiceId = "f4884616-5ded-4c74-9dbb-b565d3e3a4fd"

// Attributes to test with. Add optional properties to test here if desired.
type thirdPartyServiceResourceModel struct {
	id                    string
	name                  string
	availabilityProfileId int64
	targets               []string
	stateId               string
}

func TestAccThirdPartyService(t *testing.T) {
	resourceName := "myThirdPartyService"
	initialResourceModel := thirdPartyServiceResourceModel{
		id:                    "f4884616-5ded-4c74-9dbb-b565d3e3a4fd",
		name:                  "example",
		availabilityProfileId: 1,
		targets:               []string{"localhost:80", "localhost:443"},
		stateId:               "f4884616-5ded-4c74-9dbb-b565d3e3a4fd",
	}
	updatedResourceModel := thirdPartyServiceResourceModel{
		id:                    "f4884616-5ded-4c74-9dbb-b565d3e3a4fd",
		name:                  "updatedexample",
		availabilityProfileId: 1,
		targets:               []string{"localhost:80"},
		stateId:               "f4884616-5ded-4c74-9dbb-b565d3e3a4fd",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckThirdPartyServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThirdPartyService(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedThirdPartyServiceAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccThirdPartyService(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedThirdPartyServiceAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccThirdPartyService(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_third_party_services." + resourceName,
				ImportStateId:           thirdPartyServiceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccThirdPartyService(resourceName string, resourceModel thirdPartyServiceResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_third_party_services" "%[1]s" {
  id                      = "%[2]s"
  name                    = "%[3]s"
  availability_profile_id = "%[4]d"
  targets                 = %[5]s
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.availabilityProfileId,
		acctest.StringSliceToTerraformString(resourceModel.targets))
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedThirdPartyServiceAttributes(config thirdPartyServiceResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ThirdPartyService"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ThirdPartyServicesApi.GetThirdPartyService(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.stateId, "targets",
			config.targets, response.Targets)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckThirdPartyServiceDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ThirdPartyServicesApi.GetThirdPartyService(ctx, thirdPartyServiceId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("ThirdPartyService", thirdPartyServiceId)
	}
	return nil
}
