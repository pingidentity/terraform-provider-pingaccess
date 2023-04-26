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

const Id = "3"
const className = "com.pingidentity.pa.accesstokenvalidators.JwksEndpoint"

// Attributes to test with. Add optional properties to test here if desired.
type accessTokenValidatorResourceModel struct {
	id        int64
	classname string
	name      string
	path      string
	stateId   string
}

func TestAccAccessTokenValidator(t *testing.T) {
	resourceName := "myAccessTokenValidator"
	initialResourceModel := accessTokenValidatorResourceModel{
		classname: className,
		name:      "example",
		path:      "/example",
		id:        3,
		stateId:   "3",
	}
	updatedResourceModel := accessTokenValidatorResourceModel{
		classname: className,
		name:      "updated test name",
		path:      "/updatedexamplepath",
		id:        3,
		stateId:   "3",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAccessTokenValidatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessTokenValidator(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAccessTokenValidatorAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAccessTokenValidator(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAccessTokenValidatorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAccessTokenValidator(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_access_token_validator." + resourceName,
				ImportStateId:           Id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccAccessTokenValidator(resourceName string, resourceModel accessTokenValidatorResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_access_token_validator" "%[1]s" {
	classname = "%2s"
  name = "%3s"
  configuration = {
    path = "%4s"
  }
	id = %5d
}`, resourceName,
		resourceModel.classname,
		resourceModel.name,
		resourceModel.path,
		resourceModel.id,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedAccessTokenValidatorAttributes(config accessTokenValidatorResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Access Token Validator"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AccessTokenValidatorsApi.GetAccessTokenValidator(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}

		configResponse := response.GetConfiguration()
		configFromResponse := internaltypes.StringValueOrNull(configResponse["path"])
		err = acctest.TestAttributesMatchString(resourceType, &config.path, "path",
			config.path, configFromResponse.ValueString())
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAccessTokenValidatorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AccessTokenValidatorsApi.GetAccessTokenValidator(ctx, Id).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Access Token Validator", Id)
	}
	return nil
}
