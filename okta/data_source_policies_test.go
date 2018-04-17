package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourcePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_policies.default", "id"),
				),
			},
		},
	})
}

func testAccDataSourcePolicyConfig() string {
	return fmt.Sprintf(`
data "okta_policies" "default" {
  name = "Default Policy"
  type = "PASSWORD"
}
`)
}
