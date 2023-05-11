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

const virtualhostId = "3"

// Attributes to test with. Add optional properties to test here if desired.
type virtualhostResourceModel struct {
	id                        int64
	agentResourceCacheTTL     int64
	host                      string
	keyPairId                 int64
	port                      int64
	trustedCertificateGroupId int64
	stateId                   string
}

func TestAccVirtualHost(t *testing.T) {
	resourceName := "myVirtualHost"
	initialResourceModel := virtualhostResourceModel{
		id:                        3,
		agentResourceCacheTTL:     0,
		host:                      "test",
		keyPairId:                 3,
		port:                      1234,
		trustedCertificateGroupId: 0,
		stateId:                   "3",
	}
	updatedResourceModel := virtualhostResourceModel{
		id:                        3,
		agentResourceCacheTTL:     0,
		host:                      "updatedhostname",
		keyPairId:                 3,
		port:                      123,
		trustedCertificateGroupId: 0,
		stateId:                   "3",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckVirtualHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualHost(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedVirtualHostAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccVirtualHost(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedVirtualHostAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccVirtualHost(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_virtualhosts." + resourceName,
				ImportStateId:           virtualhostId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccVirtualHost(resourceName string, resourceModel virtualhostResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_virtualhosts" "%[1]s" {
  id                           = "%[2]d"
  agent_resource_cache_ttl     = %[3]d
  host                         = "%[4]s"
  keypair_id                   = %[5]d
  port                         = %[6]d
  trusted_certificate_group_id = %[7]d
}`, resourceName,
		resourceModel.id,
		resourceModel.agentResourceCacheTTL,
		resourceModel.host,
		resourceModel.keyPairId,
		resourceModel.port,
		resourceModel.trustedCertificateGroupId,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedVirtualHostAttributes(config virtualhostResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "VirtualHost"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VirtualhostsApi.GetVirtualHost(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "host",
			config.host, response.Host)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.stateId, "port",
			config.port, response.Port)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckVirtualHostDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.VirtualhostsApi.GetVirtualHost(ctx, virtualhostId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("VirtualHost", virtualhostId)
	}
	return nil
}
