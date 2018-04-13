package okta

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOktaProvider_nameErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderSignOn(ri)
	updatedConfig := testOktaProviderSignOn_nameErrors(ri)
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
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the name field or type field of an existing Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaProvider_typeErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderSignOn(ri)
	updatedConfig := testOktaProviderSignOn_typeErrors(ri)
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
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the name field or type field of an existing Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
				),
			},
		},
	})
}
func TestAccOktaProviderSignOn(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderSignOn(ri)
	updatedConfig := testOktaProviderSignOn_updated(ri)
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
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "OKTA_SIGN_ON"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test SignOn Policy Updated"),
				),
			},
		},
	})
}

func TestAccOktaProviderSignOn_passErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderSignOn(ri)
	updatedConfig := testOktaProviderSignOn_passErrors(ri)
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
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("password settings options not supported in the Okta SignOn Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaProviderSignOn_authpErrors(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderSignOn(ri)
	updatedConfig := testOktaProviderSignOn_authpErrors(ri)
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
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("authprovider condition options not supported in the Okta SignOn Policy"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaProviderPassword(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaProviderPassword(ri)
	updatedConfig := testOktaProviderPassword_updated(ri)
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
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc-"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "priority", "999"),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform Acceptance Test Password Policy Updated"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minlength", "12"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minlowercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minuppercase", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minnumber", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minsymbol", "0"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeusername", "false"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeattributes.0", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.excludeattributes.1", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.dictionarylookup", "true"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.maxagedays", "60"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.expirewarndays", "15"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.minageminutes", "60"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.historycount", "5"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.maxlockoutattempts", "3"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.autounlockminutes", "2"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.showlockoutfailures", "true"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.questionminlength", "10"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.recoveryemailtoken", "20160"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.smsrecovery", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "settings.0.password.0.skipunlock", "true"),
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

func testOktaProviderSignOn(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy"
}
`, rInt, rInt)
}

func testOktaProviderSignOn_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "INACTIVE"
  priority    = 999
  description = "Terraform Acceptance Test SignOn Policy Updated"
}
`, rInt, rInt)
}

func testOktaProviderSignOn_nameErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAccChanged-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
}
`, rInt, rInt)
}
func testOktaProviderSignOn_typeErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
}
`, rInt, rInt)
}

func testOktaProviderSignOn_passErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  settings {
    password {
      minlength = 12
    }
  }
}
`, rInt, rInt)
}

func testOktaProviderSignOn_authpErrors(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "OKTA_SIGN_ON"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test SignOn Policy Error Check"
  conditions {
    authprovider {
      provider = "ACTIVE_DIRECTORY"
    }
  }
}
`, rInt, rInt)
}

func testOktaProviderPassword(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test Password Policy"
}
`, rInt, rInt)
}

func testOktaProviderPassword_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_policies" "test-%d" {
  type        = "PASSWORD"
  name        = "testAcc-%d"
  status      = "INACTIVE"
  priority    = 999
  description = "Terraform Acceptance Test Password Policy Updated"
  settings {
    password {
      minlength = 12
      minlowercase = 0
      minuppercase = 0
      minnumber = 0
      minsymbol = 0
      excludeusername = false
      excludeattributes = [ "firstName", "lastName" ]
      dictionarylookup = true
      maxagedays = 60
      expirewarndays = 15
      minageminutes = 60
      historycount = 5
      maxlockoutattempts = 3
      autounlockminutes = 2
      showlockoutfailures = true
      questionminlength = 10
      recoveryemailtoken = 20160
      smsrecovery = "ACTIVE"
      skipunlock = true
    }
  }
}
`, rInt, rInt)
}
