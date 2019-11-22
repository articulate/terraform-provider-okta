package okta

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaGroupMembership(t *testing.T) {

	resourceName := fmt.Sprintf("%s.test", oktaUser)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "user_id", "TestAcc"),
					resource.TestCheckResourceAttr(resourceName, "group_id", "Smith"),
				),
			},
		},
	})
}

func testAccCheckMembershipDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*Config).oktaClient

	for _, r := range s.RootModule().Resources {
		if _, resp, err := client.User.ListUserGroups(r.Primary.ID, nil); err != nil {
			if strings.Contains(resp.Response.Status, "404") {
				continue
			}
			return fmt.Errorf("[ERROR] Error Getting Groups for User in Okta: %v", err)
		}
		return fmt.Errorf("User is still assigned to group")
	}

	return nil
}
