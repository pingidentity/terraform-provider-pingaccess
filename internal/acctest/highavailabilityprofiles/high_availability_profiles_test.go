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

const highAvailabilityProfileId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type highAvailabilityProfileResourceModel struct {
	id                     int64
	stateId                string
	name                   string
	failedRetryTimeout     float64
	failureHttpStatusCodes []string
}

func TestAccHighAvailabilityProfile(t *testing.T) {
	resourceName := "myhighAvailabilityProfile"
	initialResourceModel := highAvailabilityProfileResourceModel{
		id:                     2,
		stateId:                "2",
		failedRetryTimeout:     60,
		failureHttpStatusCodes: []string{"208", "209"},
		name:                   "name",
	}
	updatedResourceModel := highAvailabilityProfileResourceModel{
		id:                     2,
		stateId:                "2",
		failedRetryTimeout:     65,
		failureHttpStatusCodes: []string{"208"},
		name:                   "updated name",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckHighAvailabilityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHighAvailabilityProfile(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedHighAvailabilityProfileAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccHighAvailabilityProfile(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedHighAvailabilityProfileAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccHighAvailabilityProfile(resourceName, updatedResourceModel),
				ResourceName:      "pingaccess_high_availability_profile." + resourceName,
				ImportStateId:     highAvailabilityProfileId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHighAvailabilityProfile(resourceName string, resourceModel highAvailabilityProfileResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_high_availability_profile" "%[1]s" {
  id        = "%[2]d"
  classname = "com.pingidentity.pa.ha.availability.ondemand.OnDemandAvailabilityPlugin"
  name      = "%[3]s"
  configuration = {
    failed_retry_timeout      = %[4]f
    failure_http_status_codes = %[5]s
  }
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.failedRetryTimeout,
		acctest.StringSliceToTerraformString(resourceModel.failureHttpStatusCodes))
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedHighAvailabilityProfileAttributes(config highAvailabilityProfileResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "HighAvailabilityProfile"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.HighAvailabilityApi.GetAvailabilityProfile(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.name, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}

		configValues := response.GetConfiguration()
		configResponseFailedRetryTimeout := configValues["failedRetryTimeout"].(float64)
		configResponseFailureHttpStatusCodes := configValues["failureHttpStatusCodes"].([]interface{})
		err = acctest.TestAttributesMatchFloat(resourceType, &config.stateId, "failed_retry_timeout",
			config.failedRetryTimeout, configResponseFailedRetryTimeout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.stateId, "failure_http_status_codes",
			config.failureHttpStatusCodes, acctest.InterfaceSliceToStringSlice(configResponseFailureHttpStatusCodes))
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckHighAvailabilityProfileDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.HighAvailabilityApi.DeleteAvailabilityProfile(ctx, highAvailabilityProfileId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("HighAvailabilityProfile", highAvailabilityProfileId)
	}
	return nil
}
