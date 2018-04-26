package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Witiz1932@teleworm.us is a fake email address created at fakemailgenerator.com
// view inbox: http://www.fakemailgenerator.com/inbox/teleworm.us/witiz1932/

func TestAccOktaUsers(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_users.test-" + strconv.Itoa(ri)
	config := testOktaUsers(ri)
	updatedConfig := testOktaUsers_updated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUsersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firstname", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "lastname", strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "email", "Witiz1932@teleworm.us"),
					resource.TestCheckResourceAttr(resourceName, "role", "SUPER_ADMIN"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUsersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firstname", "testAcc_updated"),
					resource.TestCheckResourceAttr(resourceName, "lastname", strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "email", "Witiz1932@teleworm.us"),
					resource.TestCheckResourceAttr(resourceName, "role", "READ_ONLY_ADMIN"),
				),
			},
		},
	})
}

func TestAccOktaUsersRole_delete(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_users.test-" + strconv.Itoa(ri)
	config := testOktaUsers(ri)
	updatedConfig := testOktaUsersRole_delete(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUsersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firstname", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "lastname", strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "email", "Witiz1932@teleworm.us"),
					resource.TestCheckResourceAttr(resourceName, "role", "SUPER_ADMIN"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUsersExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "firstname", "testAcc_role_delete"),
					resource.TestCheckResourceAttr(resourceName, "lastname", strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "email", "Witiz1932@teleworm.us"),
					resource.TestCheckResourceAttr(resourceName, "role", ""),
				),
			},
		},
	})
}

func testOktaUsersExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		userID, hasID := rs.Primary.Attributes["id"]
		if !hasID {
			return fmt.Errorf("[ERROR] No id found in state")
		}
		firstName, hasFirstName := rs.Primary.Attributes["firstname"]
		if !hasFirstName {
			return fmt.Errorf("Error: no firstname found in state")
		}
		lastName, hasLastName := rs.Primary.Attributes["lastname"]
		if !hasLastName {
			return fmt.Errorf("Error: no lastname found in state")
		}

		err := testUserExists(true, userID, firstName, lastName)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*Config).oktaClient
		userRoles, _, err := client.Users.ListRoles(userID)
		if err != nil {
			return fmt.Errorf("Error: listing user role: %v", err)
		}
		role, hasRole := rs.Primary.Attributes["role"]
		if userRoles != nil {
			if !hasRole {
				return fmt.Errorf("Error: Okta role %v exists but Terraform state role does not", userRoles.Role[0].Type)
			}
			if role != userRoles.Role[0].Type {
				return fmt.Errorf("Error: Okta role %v does not match Terraform state role %v", userRoles.Role[0].Type, role)
			}
		}
		return nil
	}
	return nil
}

func testOktaUsersDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_users" {
			continue
		}

		userID, hasID := rs.Primary.Attributes["id"]
		if !hasID {
			return fmt.Errorf("[ERROR] No id found in state")
		}
		firstName, hasFirstName := rs.Primary.Attributes["firstname"]
		if !hasFirstName {
			return fmt.Errorf("Error: no firstname found in state")
		}
		lastName, hasLastName := rs.Primary.Attributes["lastname"]
		if !hasLastName {
			return fmt.Errorf("Error: no lastname found in state")
		}

		err := testUserExists(false, userID, firstName, lastName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testUserExists(expected bool, userID string, firstName string, lastName string) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	exists := false
	_, _, err := client.Users.GetByID(userID)
	if err != nil {
		if client.OktaErrorCode != "E0000007" {
			return fmt.Errorf("[ERROR] Error Listing User in Okta: %v", err)
		}
	} else {
		exists = true
	}

	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] User %v %v not found in Okta", firstName, lastName)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] User %v %v still exists in Okta", firstName, lastName)
	}
	return nil
}

func testOktaUsers(rInt int) string {
	return fmt.Sprintf(`
resource "okta_users" "test-%d" {
  firstname = "testAcc"
  lastname  = "%d"
  email     = "Witiz1932@teleworm.us"
  role      = "SUPER_ADMIN"
}
`, rInt, rInt)
}

func testOktaUsers_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_users" "test-%d" {
  firstname = "testAcc_updated"
  lastname  = "%d"
  email     = "Witiz1932@teleworm.us"
  role      = "READ_ONLY_ADMIN"
}
`, rInt, rInt)
}

func testOktaUsersRole_delete(rInt int) string {
	return fmt.Sprintf(`
resource "okta_users" "test-%d" {
  firstname = "testAcc_role_delete"
  lastname  = "%d"
  email     = "Witiz1932@teleworm.us"
}
`, rInt, rInt)
}
