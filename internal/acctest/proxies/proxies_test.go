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

const proxyId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type proxyResourceModel struct {
	id       int64
	host     string
	name     string
	port     int64
	username string
	stateId  string
}

func TestAccProxy(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := proxyResourceModel{
		id:       2,
		host:     "example",
		name:     "example",
		port:     1234,
		username: "example",
		stateId:  "2",
	}
	updatedResourceModel := proxyResourceModel{
		id:       2,
		host:     "updatedexample",
		name:     "updated example",
		port:     1235,
		username: "updatedexample",
		stateId:  "2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProxy(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedProxyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccProxy(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedProxyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccProxy(resourceName, updatedResourceModel),
				ResourceName:      "pingaccess_proxy." + resourceName,
				ImportStateId:     proxyId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccProxy(resourceName string, resourceModel proxyResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_proxy" "%[1]s" {
  id   = %[2]d
  host = "%[3]s"
  name = "%[4]s"
  port = %[5]d
  password = {
    value = "examplePassword"
  }
  requires_authentication = true
  username                = "%[6]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.host,
		resourceModel.name,
		resourceModel.port,
		resourceModel.username,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedProxyAttributes(config proxyResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Proxy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ProxiesApi.GetProxy(ctx, config.stateId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "host",
			config.host, response.Host)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.stateId, "port",
			config.port, response.Port)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "username",
			config.username, *response.Username)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckProxyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ProxiesApi.GetProxy(ctx, proxyId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Proxy", proxyId)
	}
	return nil
}
