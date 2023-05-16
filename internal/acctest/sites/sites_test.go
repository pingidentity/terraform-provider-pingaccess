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

const siteId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type siteResourceModel struct {
	id      int64
	name    string
	targets []string
	stateId string
}

func TestAccSite(t *testing.T) {
	resourceName := "mySite"
	initialResourceModel := siteResourceModel{
		id:      2,
		name:    "example",
		targets: []string{"localhost:80", "localhost:443"},
		stateId: "2",
	}
	updatedResourceModel := siteResourceModel{
		id:      2,
		name:    "updatedexample",
		targets: []string{"localhost:80"},
		stateId: "2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSite(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSiteAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSite(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSiteAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSite(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_sites." + resourceName,
				ImportStateId:           siteId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccSite(resourceName string, resourceModel siteResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_sites" "%[1]s" {
  id      = "%[2]d"
  name    = "%[3]s"
  targets = %[4]s
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		acctest.StringSliceToTerraformString(resourceModel.targets))
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedSiteAttributes(config siteResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Site"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SitesApi.GetSite(ctx, config.stateId).Execute()

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
func testAccCheckSiteDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.SitesApi.GetSite(ctx, siteId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Site", siteId)
	}
	return nil
}
