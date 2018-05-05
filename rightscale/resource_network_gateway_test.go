package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	networkGatewayDescription = "Terraform RightScale provider test Network Gateway"
	networkGatewayType        = "internet"
)

func TestAccRightScaleNetworkGateway(t *testing.T) {
	t.Parallel()

	var (
		NetworkGatewayName = "terraform-test-" + testString + acctest.RandString(10)
		depl               cm15.NetworkGateway
		cloudHref          = getTestCloudFromEnv()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkGateway(NetworkGatewayName, networkGatewayDescription, networkGatewayType, cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkGatewayExists("rightscale_network_gateway.test_network_gateway", &depl),
					testAccCheckNetworkGatewayDescription(&depl, networkGatewayDescription),
				),
			},
		},
	})
}

func testAccCheckNetworkGatewayExists(n string, depl *cm15.NetworkGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().NetworkGatewayLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show()
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckNetworkGatewayDescription(depl *cm15.NetworkGateway, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckNetworkGatewayDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_network_gateway" {
			continue
		}

		loc := c.NetworkGatewayLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of Network: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, depl := range depls {
			if string(depl.Locator(c).Href) == self {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccNetworkGateway(name string, desc string, typ string, cloud string) string {
	return fmt.Sprintf(`
		resource "rightscale_network_gateway" "test_network_gateway" {
		   name = %q
		   description = %q
		   type = %q
		   cloud_href = %q
		 }
`, name, desc, typ, cloud)
}
