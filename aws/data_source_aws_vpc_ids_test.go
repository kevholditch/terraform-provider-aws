package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAwsVpcIDs_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpcIDsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrGreaterThanValue("data.aws_vpc_ids.all", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsVpcIDs_tags(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpcIDsConfig_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_vpc_ids.selected", "ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsVpcIDs_filters(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpcIDsConfig_filters(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_vpc_ids.selected", "ids.#", "1"),
				),
			},
		},
	})
}

func testCheckResourceAttrGreaterThanValue(name, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", name, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", name, ms.Path)
		}

		if v, ok := is.Attributes[key]; !ok || v == value {
			if !ok {
				return fmt.Errorf("%s: Attribute '%s' not found", name, key)
			}

			return fmt.Errorf(
				"%s: Attribute '%s' is not greater than %#v, got %#v",
				name,
				key,
				value,
				v)
		}
		return nil

	}
}

func testAccDataSourceAwsVpcIDsConfig() string {
	return fmt.Sprintf(`
	resource "aws_vpc" "test-vpc" {
  		cidr_block = "10.0.0.0/24"
	}

	data "aws_vpc_ids" "all" {
		depends_on = ["aws_vpc.test-vpc"]
	}
	`)
}

func testAccDataSourceAwsVpcIDsConfig_tags(rName string) string {
	return fmt.Sprintf(`
	resource "aws_vpc" "test-vpc" {
  		cidr_block = "10.0.0.0/24"

  		tags {
  			Name = "testacc-vpc-%s"
  			Service = "testacc-test"
  		}
	}

	data "aws_vpc_ids" "selected" {
		tags {
			Name = "testacc-vpc-%s"
			Service = "testacc-test"
		}
		depends_on = ["aws_vpc.test-vpc"]
	}
	`, rName, rName)
}

func testAccDataSourceAwsVpcIDsConfig_filters(rName string) string {
	return fmt.Sprintf(`
	resource "aws_vpc" "test-vpc" {
  		cidr_block = "192.168.0.0/25"
  		tags {
  			Name = "testacc-vpc-%s"
  		}
	}

	data "aws_vpc_ids" "selected" {
		filter {
			name = "cidr-block-association.cidr-block"
    		values = ["192.168.0.0/25"]
		}
		filter {
   			name = "tag:Name"
    		values = ["testacc-vpc-%s"]
  		}
		depends_on = ["aws_vpc.test-vpc"]
	}
	`, rName, rName)
}
