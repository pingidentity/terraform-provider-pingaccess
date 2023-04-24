// package acctest_test

// import (
// 	"fmt"
// 	"testing"

// 	acctest "terraform-provider-pingaccess/internal/acctest"
// 	provider "terraform-provider-pingaccess/internal/provider"

// 	"github.com/hashicorp/terraform-plugin-framework/providerserver"
// 	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// )

// const engineListenerId = 2

// // Attributes to test with. Add optional properties to test here if desired.
// type engineListenerResourceModel struct {
// 	id                        int64
// 	name                      string
// 	port                      int64
// 	secure                    bool
// 	trustedCertificateGroupId int64
// }

// func TestAccEngineListener(t *testing.T) {
// 	resourceName := "myresource"
// 	initialResourceModel := engineListenerResourceModel{
// 		id:                        engineListenerId,
// 		name:                      "test",
// 		port:                      1234,
// 		secure:                    true,
// 		trustedCertificateGroupId: 0,
// 	}
// 	updatedResourceModel := engineListenerResourceModel{
// 		id:                        engineListenerId,
// 		name:                      "updated test name",
// 		port:                      123,
// 		secure:                    false,
// 		trustedCertificateGroupId: 0,
// 	}

// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
// 		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
// 			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
// 		},
// 		CheckDestroy: testAccCheckEngineListenerDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccEngineListener(resourceName, initialResourceModel),
// 				Check:  testAccCheckExpectedEngineListenerAttributes(initialResourceModel),
// 			},
// 			{
// 				// Test updating some fields
// 				Config: testAccEngineListener(resourceName, updatedResourceModel),
// 				Check:  testAccCheckExpectedEngineListenerAttributes(updatedResourceModel),
// 			},
// 			{
// 				// Test importing the resource
// 				Config:                  testAccEngineListener(resourceName, updatedResourceModel),
// 				ResourceName:            "pingaccess_engine_listener." + resourceName,
// 				ImportStateId:           updatedResourceModel.name,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"last_updated"},
// 			},
// 		},
// 	})
// }

// func testAccEngineListener(resourceName string, resourceModel engineListenerResourceModel) string {
// 	return fmt.Sprintf(`
// resource "pingaccess_engine_listener" "%[1]s" {
//   id            = %[2]d
//   name               = "%[3]s"
//   port      = %[4]d
//   secure              = %[5]t
// 	trusted_certificate_group_id = %[6]d
// }`, resourceName,
// 		resourceModel.id,
// 		resourceModel.name,
// 		resourceModel.port,
// 		resourceModel.secure,
// 		resourceModel.trustedCertificateGroupId,
// 	)
// }

// // Test that the expected attributes are set on the PingDirectory server
// func testAccCheckExpectedEngineListenerAttributes(config engineListenerResourceModel) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		testClient := acctest.TestClient()
// 		ctx := acctest.TestBasicAuthContext()
// 		apiReadListener, httpResp, err := apiClient.DefaultApi.EngineListenersIdGet(config.ProviderBasicAuthContext(ctx, providerConfig), internaltypes.Int64ToString(state.Id)).Execute()
// 		if err != nil {
// 			return err
// 		}
// 		// Verify that attributes have expected values
// 		resourceType := "Engine Listener"
// 		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "id",
// 			config.id, response.EngineListenerResponse.Id)
// 		if err != nil {
// 			return err
// 		}
// 		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name",
// 			config.name, response.EngineListenerResponse.Name)
// 		if err != nil {
// 			return err
// 		}
// 		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "port",
// 			config.port, response.EngineListenerResponse.Port)
// 		if err != nil {
// 			return err
// 		}
// 		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "secure",
// 			config.secure, response.EngineListenerResponse.Secure)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "trusted-certificate-group-id",
// 			config.trustedCertificateGroupId, response.EngineListenerResponse.TrustedCertificateGroupId)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// // Test that any objects created by the test are destroyed
// func testAccCheckEngineListenerDestroy(s *terraform.State) error {
// 	testClient := acctest.TestClient()
// 	ctx := acctest.TestBasicAuthContext()
// 	_, _, err := testClient.RequestCriteriaApi.GetRequestCriteria(ctx, engineListenerId).Execute()
// 	if err == nil {
// 		return acctest.ExpectedDestroyError("Engine Listener", engineListenerId)
// 	}
// 	return nil
// }
