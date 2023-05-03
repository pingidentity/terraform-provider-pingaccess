package acctest_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingaccess/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingaccess/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const Id = "f9ee7432-01c0-46ed-8887-edae88ddba45"

// Attributes to test with. Add optional properties to test here if desired.
type acmeserversResourceModel struct {
	id      string
	stateId string
	name    string
	url     string
}

func TestAccAcmeServer(t *testing.T) {
	resourceName := "myAcmeServer"
	initialResourceModel := acmeserversResourceModel{
		id:      Id,
		stateId: "f9ee7432-01c0-46ed-8887-edae88ddba45",
		name:    "example",
		url:     "https://thisisanexample.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAcmeServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAcmeServer(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAcmeServerAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAcmeServer(resourceName, initialResourceModel),
				ResourceName:            "pingaccess_acme_servers." + resourceName,
				ImportStateId:           Id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccAcmeServer(resourceName string, resourceModel acmeserversResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_acme_servers" "%[1]s" {
  id            = "%[2]s"
  name               = "%[3]s"
  url      = "%[4]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.url,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedAcmeServerAttributes(config acmeserversResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AcmeServer"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AcmeApi.GetAcmeServer(ctx, Id).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "url",
			config.url, response.Url)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAcmeServerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AcmeApi.GetAcmeServer(ctx, Id).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AcmeServer", Id)
	}
	return nil
}
