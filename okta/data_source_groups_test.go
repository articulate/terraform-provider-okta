
// this data source test is a temporary "fix" for the policy tests to lookup the Everyone group ID
// this data source test needs to be deleted after the groups resource is added to the provider

package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_groups.everyone", "id"),
				),
			},
		},
	})
}

func testAccDataSourceGroupConfig() string {
	return fmt.Sprintf(`
data "okta_groups" "everyone" {
  name = "Everyone"
}
`)
}
