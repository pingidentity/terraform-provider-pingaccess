package acctest_test

import (
	"fmt"
	"testing"

	"github.com/pingidentity/terraform-provider-pingaccess/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingaccess/internal/provider"
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const identityMappingId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type identityMappingResourceModel struct {
	id int64	
	className string	
	configuration determine this value manually	
	name string
}

func TestAccIdentityMapping(t *testing.T) {
	resourceName := "myIdentityMapping"
	initialResourceModel := identityMappingResourceModel{
		id: fill in test value,	
		className: fill in test value,	
		configuration: fill in test value,	
		name: fill in test value,
	}
	updatedResourceModel := identityMappingResourceModel{
		id: fill in test value,	
		className: fill in test value,	
		configuration: fill in test value,	
		name: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckIdentityMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdentityMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdentityMapping(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdentityMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccIdentityMapping(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_identity_mappings." + resourceName,
				ImportStateId:           Id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccIdentityMapping(resourceName string, resourceModel identityMappingResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_identity_mappings" "%[1]s" {
	FILL THIS IN
}`, resourceName,
	resourceModel.className,		
	resourceModel.configuration,		
	resourceModel.name,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedIdentityMappingAttributes(config identityMappingResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IdentityMapping"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.IdentityMappingsApi.GetIdentityMapping(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIdentityMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.IdentityMappingsApi.GetIdentityMapping(ctx, Id).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdentityMapping", Id)
	}
	return nil
}
