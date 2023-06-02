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

const certificateId = "1"
const fileData = "MIIDmjCCAoICCQDncp3LMAO6YjANBgkqhkiG9w0BAQsFADCBjjELMAkGA1UEBhMCVVMxDDAKBgNVBAgMA1NKQzERMA8GA1UEBwwIc2FuIEpvc2UxFjAUBgNVBAoMDXBpbmcgSWRlbnRpdHkxDzANBgNVBAsMBkRldm9wczEWMBQGA1UEAwwNdGVycmFmb3JtdGVzdDEdMBsGCSqGSIb3DQEJARYOdGVzdEBnbWFpbC5jb20wHhcNMjMwNTMwMTU1OTE5WhcNMjQwNTI5MTU1OTE5WjCBjjELMAkGA1UEBhMCVVMxDDAKBgNVBAgMA1NKQzERMA8GA1UEBwwIc2FuIEpvc2UxFjAUBgNVBAoMDXBpbmcgSWRlbnRpdHkxDzANBgNVBAsMBkRldm9wczEWMBQGA1UEAwwNdGVycmFmb3JtdGVzdDEdMBsGCSqGSIb3DQEJARYOdGVzdEBnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDB7u+oHHQgGrZdCk74A4XJzjzhMT9MN1MJIqar+96rKogDmt3LnCh+oN5hxy0QPjrW9SiRHPZME+e6YWtBNfg21KDws2nLoH/eGmb45ObM/nApX4oFZD06ccW4zWjxuxEdKzKAMWMP60UxCZwnK99cIRMYs0x85lHhcLfTuA3VAwg95X+2FxQDk8sAdNdl1zhWaR2YS+nrmP/iheG2fT8cVLTGdklPqL9nrUDAwwUyX5I8PLsLPzJzMoXV+on4zjypNxfXt2MmuLHOGxwgxvUVRiVeCTSMo1y763OUAnds1L+uJNq1vvsD0iFwyA78I3EzaX9c5Vxhbk+3JKFD1gY1AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAGqlkRIgsAFE6/WBayYlsITtnxJooTJyZ8CHFulRMskMYdoETYUeN5FqmJ05PGUHgXX0/3fQ9RYD3Mfuupm1Vqgx8q/v5cIrBefU7zW3bjy/BMAONkPAr617NkbHAj2XC1t5YFr6Vnnx9JQoIl70slBGABPwSkahrReE5f87qkkWqVI8aiuAzu0GRkMHbv1XzGfXfVF/iK9Lq6x80tyiqL987Krw6hHPlxS4GXjwvWWO0f0GfNwENxSv6uwxvCFIp01x7LHbkPHJvMH2Z5wSZges5ZDv/rciunSZ2xYh/jGzM1gIz29DBpmayl4AwKi5/ix7p3ujCA1jdlT+nlBZ/js="

// Attributes to test with. Add optional properties to test here if desired.
type engineListenerResourceModel struct {
	id      int64
	alias   string
	stateId string
}

func TestAccCertificate(t *testing.T) {
	resourceName := "myresource"
	initialResourceModel := engineListenerResourceModel{
		id:      1,
		alias:   "test",
		stateId: "1",
	}
	updatedResourceModel := engineListenerResourceModel{
		id:      1,
		alias:   "updated test name",
		stateId: "1",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingaccess": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCertificateAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccCertificate(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCertificateAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccCertificate(resourceName, updatedResourceModel),
				ResourceName:      "pingaccess_certificates." + resourceName,
				ImportStateId:     certificateId,
				ImportState:       false,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCertificate(resourceName string, resourceModel engineListenerResourceModel) string {
	return fmt.Sprintf(`
resource "pingaccess_certificates" "%[1]s" {
  alias     = "%[2]s"
  file_data = "%[3]s"
}`, resourceName,
		resourceModel.alias,
		fileData,
	)
}

// Test that the expected attributes are set on the PingAccess server
func testAccCheckExpectedCertificateAttributes(config engineListenerResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "Certificate"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, httpResp, err := testClient.CertificatesApi.GetTrustedCert(ctx, config.stateId).Execute()
		if httpResp.StatusCode != 200 {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "alias",
			config.alias, response.Alias)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckCertificateDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	httpResp, _ := testClient.CertificatesApi.DeleteTrustedCert(ctx, certificateId).Execute()
	if httpResp.StatusCode == 200 {
		return acctest.ExpectedDestroyError("Certificate", certificateId)
	}
	return nil
}
