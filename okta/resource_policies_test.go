package okta

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaPoviderSignOn_create(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPoviderSignOn(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
		},
	})
}

func TestAccOktaPoviderSignOn_update(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaPoviderSignOn(ri)
	updatedConfig := testOktaPoviderSignOn_updated(ri)
	resourceName := "okta_policies.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
		},
	})
}

func testOktaProviderExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[ERROR] Resource Not found: %s", name)
		}

		policyType, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("[ERROR] No type found in state for Policy")
		}
		policyName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy")
		}

		err := testPolicyExists(true, policyType, policyName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testOktaProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_policies" {
			continue
		}

		policyType, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("[ERROR] No type found in state for Policy")
		}
		policyName, hasName := rs.Primary.Attributes["name"]
		if !hasName {
			return fmt.Errorf("[ERROR] No name found in state for Policy")
		}

		err := testPolicyExists(false, policyType, policyName)
		if err != nil {
			return err
		}
	}
	return nil
}

func testPolicyExists(expected bool, policyType string, policyName string) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	exists := false
	policies, _, err := client.Policies.GetPoliciesByType(policyType)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policies in Okta: %v", err)
	}
	if policies != nil {
		for _, policy := range policies.Policies {
			if policy.Name == policyName {
				exists = true
			}
		}
	}
	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] Policy %v not found in Okta", policyName)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] Policy %v still exists in Okta", policyName)
	}
	return nil
}

func testOktaPoviderSignOn(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, rInt)
}

func testOktaPoviderSignOn_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc"
  status      = "INACTIVE"
  priority    = 999
  description = "Terraform Acceptance Test SignOn Policy Updated"
}
`, rInt)
}
