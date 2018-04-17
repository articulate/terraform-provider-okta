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
