package licensemanager_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/envvar"
)

func TestAccLicenseManagerReceivedLicensesDataSource_basic(t *testing.T) {
	datasourceName := "data.aws_licensemanager_received_licenses.test"
	licenseARN := envvar.SkipIfEmpty(t, licenseARNKey, envVarLicenseARNKeyError)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReceivedLicensesDataSourceConfig_arns(licenseARN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "arns.#", "1"),
				),
			},
		},
	})
}

func TestAccLicenseManagerReceivedLicensesDataSource_empty(t *testing.T) {
	datasourceName := "data.aws_licensemanager_received_licenses.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReceivedLicensesDataSourceConfig_empty(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "arns.#", "0"),
				),
			},
		},
	})
}

func testAccReceivedLicensesDataSourceConfig_arns(licenseARN string) string {
	return fmt.Sprintf(`
data "aws_licensemanager_received_licenses" "test" {
  filter {
    name = "ProductSKU"
    values = [
      data.aws_licensemanager_received_license.test.product_sku
    ]
  }
}

data "aws_licensemanager_received_license" "test" {
  license_arn = %[1]q
}
`, licenseARN)
}

func testAccReceivedLicensesDataSourceConfig_empty() string {
	return `
data "aws_licensemanager_received_licenses" "test" {
  filter {
    name = "IssuerName"
    values = [
      "This Is Fake"
    ]
  }
}
`
}
