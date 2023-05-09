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
	internaltypes "github.com/pingidentity/terraform-provider-pingaccess/internal/types"
)

const hsmProviderId = "2"
const className = "com.pingidentity.pa.hsm.pkcs11.plugin.PKCS11HsmProvider"

// Attributes to test with. Add optional properties to test here if desired.
type hsmProviderResourceModel struct {
	id        int64
	classname string
	slot_id   string
	password  string
	library   string
	name      string
	stateId   string
}

func TestAccHsmProvider(t *testing.T) {
	resourceName := "myHsmProvider"
	initialResourceModel := hsmProviderResourceModel{
		classname: className,
		name:      "example",
		id:        2,
		password:  "password",
		stateId:   "2",
		slot_id:   "2",
		library:   "example2",
	}
	updatedResourceModel := hsmProviderResourceModel{
		classname: className,
		name:      "updated test name",
		id:        2,
		password:  "password",
		stateId:   "2",
		slot_id:   "3",
		library:   "/updatedexample3",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckHsmProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHsmProvider(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedHsmProviderAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccHsmProvider(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedHsmProviderAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccHsmProvider(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_hsm_providers." + resourceName,
				ImportStateId:           hsmProviderId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccHsmProvider(resourceName string, resourceModel hsmProviderResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_hsm_providers" "%[1]s" {
  id        = %[2]d
  classname = "%[3]s"
  name      = "%[4]s"
  configuration = {
    slot_id  = "%[5]s"
    password = "%[6]s"
    library  = "%[7]s"
  }
}`, resourceName,
		resourceModel.id,
		resourceModel.classname,
		resourceModel.name,
		resourceModel.slot_id,
		resourceModel.password,
		resourceModel.library,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedHsmProviderAttributes(config hsmProviderResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "HsmProvider"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.HsmProvidersApi.GetHsmProvider(ctx, config.stateId).Execute()

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
		configFromResponse := internaltypes.StringValueOrNull(configResponse["slotId"])
		err = acctest.TestAttributesMatchString(resourceType, &config.slot_id, "slot_id",
			config.slot_id, configFromResponse.ValueString())
		if err != nil {
			return err
		}
		configFromResponselib := internaltypes.StringValueOrNull(configResponse["library"])
		err = acctest.TestAttributesMatchString(resourceType, &config.library, "library",
			config.library, configFromResponselib.ValueString())
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckHsmProviderDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.HsmProvidersApi.GetHsmProvider(ctx, hsmProviderId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("HsmProvider", hsmProviderId)
	}
	return nil
}
