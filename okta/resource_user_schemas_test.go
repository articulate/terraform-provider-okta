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

func TestAccOktaUserSchemas_baseCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas_baseCheck(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Editing a base user subschema not supported in this terraform provider at this time"),
				PlanOnly:    true,
			},
		},
	})
}

// uncomment this test when there's support to edit base subschemas
//func TestAccOktaUserSchemas_subschemaCheck(t *testing.T) {
//	ri := acctest.RandInt()
//	config := testOktaUserSchemas(ri)
//	updatedConfig := testOktaUserSchemas_subschemaCheck(ri)
//	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testOktaUserSchemasDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: config,
//				Check: resource.ComposeTestCheckFunc(
//					testOktaUserSchemasExists(resourceName),
//				),
//			},
//			{
//				Config:      updatedConfig,
//				ExpectError: regexp.MustCompile("You cannot change the subschema field for an existing User SubSchema"),
//				PlanOnly:    true,
//				Check: resource.ComposeTestCheckFunc(
//					testOktaUserSchemasExists(resourceName),
//				),
//			},
//		},
//	})
//}

func TestAccOktaUserSchemas_indexCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_indexCheck(ri)
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the index field for an existing User SubSchema"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas_typeCheck(t *testing.T) {
	ri := acctest.RandInt()
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_typeCheck(ri)
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
			{
				Config:      updatedConfig,
				ExpectError: regexp.MustCompile("You cannot change the type field for an existing User SubSchema"),
				PlanOnly:    true,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
				),
			},
		},
	})
}

func TestAccOktaUserSchemas(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_user_schemas.test-" + strconv.Itoa(ri)
	config := testOktaUserSchemas(ri)
	updatedConfig := testOktaUserSchemas_updated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testOktaUserSchemasDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testOktaUserSchemasExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subschema", "custom"),
					resource.TestCheckResourceAttr(resourceName, "index", "testAcc"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "title", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "string"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test updated"),
					resource.TestCheckResourceAttr(resourceName, "required", "true"),
					resource.TestCheckResourceAttr(resourceName, "minlength", "1"),
					resource.TestCheckResourceAttr(resourceName, "maxlength", "50"),
					resource.TestCheckResourceAttr(resourceName, "permissions", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceName, "master", "OKTA"),
				),
			},
		},
	})
}

// type tests -> boolean, number, interger, & array
// test enum & oneof

func testOktaUserSchemasExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		subschema, hasSchema := rs.Primary.Attributes["subschema"]
		if !hasSchema {
			return fmt.Errorf("Error: No subschema found in state")
		}
		index, hasIndex := rs.Primary.Attributes["index"]
		if !hasIndex {
			return fmt.Errorf("Error: no index found in state")
		}
		_, hasTitle := rs.Primary.Attributes["title"]
		if !hasTitle {
			return fmt.Errorf("Error: no title found in state")
		}
		_, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("Error: no type found in state")
		}

		err := testUserSchemaExists(true, subschema, index)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func testOktaUserSchemasDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "okta_user_schemas" {
			continue
		}

		subschema, hasSchema := rs.Primary.Attributes["subschema"]
		if !hasSchema {
			return fmt.Errorf("Error: No subschema found in state")
		}
		index, hasIndex := rs.Primary.Attributes["index"]
		if !hasIndex {
			return fmt.Errorf("Error: no index found in state")
		}
		_, hasTitle := rs.Primary.Attributes["title"]
		if !hasTitle {
			return fmt.Errorf("Error: no title found in state")
		}
		_, hasType := rs.Primary.Attributes["type"]
		if !hasType {
			return fmt.Errorf("Error: no type found in state")
		}

		err := testUserSchemaExists(false, subschema, index)
		if err != nil {
			return err
		}
	}
	return nil
}

func testUserSchemaExists(expected bool, scope string, index string) error {
	client := testAccProvider.Meta().(*Config).oktaClient

	exists := false
	subschemas, _, err := client.Schemas.GetUserSubSchemaIndex(scope)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing User Subschemas in Okta: %v", err)
	}
	for _, key := range subschemas {
		if key == index {
			exists = true
			break
		}
	}

	if expected == true && exists == false {
		return fmt.Errorf("[ERROR] User Schema %v not found in Okta", index)
	} else if expected == false && exists == true {
		return fmt.Errorf("[ERROR] User Schema %v still exists in Okta", index)
	}
	return nil
}

func testOktaUserSchemas(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema   = "custom"
  index       = "testAcc%d"
  title       = "terraform acceptance test"
  type        = "string"
  description = "terraform acceptance test"
}
`, rInt, rInt)
}

func testOktaUserSchemas_updated(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema   = "custom"
  index       = "testAcc%d"
  title       = "terraform acceptance test updated"
  type        = "string"
  description = "terraform acceptance test updated"
  required    = true
  minlength   = 1
  maxlength   = 50
  permissions = "READ_WRITE"
  master      = "OKTA"
}
`, rInt, rInt)
}

func testOktaUserSchemas_baseCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "base"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "string"
}
`, rInt, rInt)
}

// uncomment this test when there's support to edit base subschemas
//func testOktaUserSchemas_subschemaCheck(rInt int) string {
//	return fmt.Sprintf(`
//resource "okta_user_schemas" "test-%d" {
//  subschema = "base"
//  index     = "testAcc%d"
//  title     = "terraform acceptance test"
//  type      = "string"
//}
//`, rInt, rInt)
//}

func testOktaUserSchemas_indexCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAccChanged%d"
  title     = "terraform acceptance test"
  type      = "string"
}
`, rInt, rInt)
}

func testOktaUserSchemas_typeCheck(rInt int) string {
	return fmt.Sprintf(`
resource "okta_user_schemas" "test-%d" {
  subschema = "custom"
  index     = "testAcc%d"
  title     = "terraform acceptance test"
  type      = "array"
}
`, rInt, rInt)
}
