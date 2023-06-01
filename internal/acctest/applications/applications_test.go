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

const applicationId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type applicationResourceModel struct {
	id int64	
	accessValidatorId int64	
	agentCacheInvalidatedExpiration int64	
	agentCacheInvalidatedResponseDuration int64	
	agentId int64	
	allowEmptyPathSegments bool	
	applicationType ApplicationTypeView didn't match any expected values go to Web API to find what's needed	
	authenticationChallengePolicyId string	
	caseSensitivePath bool	
	contextRoot string	
	defaultAuthType DefaultAuthTypeView didn't match any expected values go to Web API to find what's needed	
	description string	
	destination DestinationView didn't match any expected values go to Web API to find what's needed	
	enabled bool	
	fallbackPostEncoding string	
	identityMappingIds Map[string,int] didn't match any expected values go to Web API to find what's needed	
	issuer string	
	lastModified int64	
	manualOrderingEnabled bool	
	name string	
	policy Map[string,List[PolicyItem]] didn't match any expected values go to Web API to find what's needed	
	realm string	
	requireHTTPS bool	
	resourceOrder []string	
	sidebandClientId string	
	siteId int64	
	spaSupportEnabled bool	
	virtualHostIds []string	
	webSessionId int64
}

func TestAccApplication(t *testing.T) {
	resourceName := "myApplication"
	initialResourceModel := applicationResourceModel{
		id: fill in test value,	
		accessValidatorId: fill in test value,	
		agentCacheInvalidatedExpiration: fill in test value,	
		agentCacheInvalidatedResponseDuration: fill in test value,	
		agentId: fill in test value,	
		allowEmptyPathSegments: fill in test value,	
		applicationType: fill in test value,	
		authenticationChallengePolicyId: fill in test value,	
		caseSensitivePath: fill in test value,	
		contextRoot: fill in test value,	
		defaultAuthType: fill in test value,	
		description: fill in test value,	
		destination: fill in test value,	
		enabled: fill in test value,	
		fallbackPostEncoding: fill in test value,	
		identityMappingIds: fill in test value,	
		issuer: fill in test value,	
		lastModified: fill in test value,	
		manualOrderingEnabled: fill in test value,	
		name: fill in test value,	
		policy: fill in test value,	
		realm: fill in test value,	
		requireHTTPS: fill in test value,	
		resourceOrder: fill in test value,	
		sidebandClientId: fill in test value,	
		siteId: fill in test value,	
		spaSupportEnabled: fill in test value,	
		virtualHostIds: fill in test value,	
		webSessionId: fill in test value,
	}
	updatedResourceModel := applicationResourceModel{
		id: fill in test value,	
		accessValidatorId: fill in test value,	
		agentCacheInvalidatedExpiration: fill in test value,	
		agentCacheInvalidatedResponseDuration: fill in test value,	
		agentId: fill in test value,	
		allowEmptyPathSegments: fill in test value,	
		applicationType: fill in test value,	
		authenticationChallengePolicyId: fill in test value,	
		caseSensitivePath: fill in test value,	
		contextRoot: fill in test value,	
		defaultAuthType: fill in test value,	
		description: fill in test value,	
		destination: fill in test value,	
		enabled: fill in test value,	
		fallbackPostEncoding: fill in test value,	
		identityMappingIds: fill in test value,	
		issuer: fill in test value,	
		lastModified: fill in test value,	
		manualOrderingEnabled: fill in test value,	
		name: fill in test value,	
		policy: fill in test value,	
		realm: fill in test value,	
		requireHTTPS: fill in test value,	
		resourceOrder: fill in test value,	
		sidebandClientId: fill in test value,	
		siteId: fill in test value,	
		spaSupportEnabled: fill in test value,	
		virtualHostIds: fill in test value,	
		webSessionId: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplication(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedApplicationAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccApplication(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedApplicationAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccApplication(resourceName, updatedResourceModel),
				ResourceName:            "pingaccess_applications." + resourceName,
				ImportStateId:           applicationId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"items"},
			},
		},
	})
}

func testAccApplication(resourceName string, resourceModel applicationResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_applications" "%[1]s" {
	FILL THIS IN
}`, resourceName,
	resourceModel.agentId,		
	resourceModel.authenticationChallengePolicyId,		
	resourceModel.contextRoot,		
	resourceModel.defaultAuthType,		
	resourceModel.name,		
	resourceModel.sidebandClientId,		
	resourceModel.siteId,		
	resourceModel.spaSupportEnabled,		
	resourceModel.virtualHostIds,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedApplicationAttributes(config applicationResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Application"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ApplicationsApi.GetApplication(ctx, config.stateId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckApplicationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.ApplicationsApi.GetApplication(ctx, applicationId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("Application", applicationId)
	}
	return nil
}
