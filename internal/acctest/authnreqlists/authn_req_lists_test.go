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

const authnReqListId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type authnReqListResourceModel struct {
	id        int64
	authnReqs []string
	stateId   string
	name      string
}

func TestAccAuthnReqList(t *testing.T) {
	resourceName := "myAuthnReqList"
	initialResourceModel := authnReqListResourceModel{
		id:        2,
		stateId:   "2",
		authnReqs: []string{"example1", "example2"},
		name:      "name",
	}
	updatedResourceModel := authnReqListResourceModel{
		id:        2,
		stateId:   "2",
		authnReqs: []string{"example1"},
		name:      "updated name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAuthnReqListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthnReqList(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthnReqListAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthnReqList(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthnReqListAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAuthnReqList(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_authn_req_lists." + resourceName,
				ImportStateId:           authnReqListId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccAuthnReqList(resourceName string, resourceModel authnReqListResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_authn_req_lists" "%[1]s" {
  id         = "%[2]d"
  name       = "%[3]s"
  authn_reqs = %[4]s
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		acctest.StringSliceToTerraformString(resourceModel.authnReqs))
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedAuthnReqListAttributes(config authnReqListResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthnReqList"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthnReqListsApi.GetAuthnReqList(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.stateId, "authn_reqs",
			config.authnReqs, response.AuthnReqs)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthnReqListDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AuthnReqListsApi.GetAuthnReqList(ctx, authnReqListId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthnReqList", authnReqListId)
	}
	return nil
}
