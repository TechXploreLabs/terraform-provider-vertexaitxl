// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccModelGardenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccModelGardenResourceConfig("gemini-1.5-pro-002"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vertexaitxl_model_garden.test",
						tfjsonpath.New("model_name"),
						knownvalue.StringExact("gemini-1.5-pro-002"),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccModelGardenResourceConfig("gemini-1.5-pro-002"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"vertexaitxl_model_garden.test",
						tfjsonpath.New("model_name"),
						knownvalue.StringExact("gemini-1.5-pro-002"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccModelGardenResourceConfig(model_name string) string {
	return fmt.Sprintf(`
resource "vertexaitxl_model_garden" "test" {
  prompt     = <<EOF
		The week ahead brings a mix of weather conditions.
		Sunday is expected to be sunny with a temperature of 77°F and a humidity level of 50 percent. Winds will be light at around 10 km/h.
		Monday will see partly cloudy skies with a slightly cooler temperature of 72°F and humidity increasing to 55 percent. Winds will pick up slightly to around 15 km/h.
		Tuesday brings rain showers, with temperatures dropping to 64°F and humidity rising to 70 percent. Expect stronger winds at 20 km/h.
		Wednesday may see thunderstorms, with a temperature of 68°F and high humidity of 75 percent. Winds will be gusty at 25 km/h.
		Thursday will be cloudy with a temperature of 66°F and moderate humidity at 60 percent. Winds will ease slightly to 18 km/h.
		Friday returns to partly cloudy conditions, with a temperature of 73°F and lower humidity at 45 percent. Winds will be light at 12 km/h.
		Finally, Saturday rounds off the week with sunny skies, a temperature of 80°F, and a humidity level of 40 percent. Winds will be gentle at 8 km/h.
    EOF
  project_id = "my-project"
  location   = "us-central1"
  model_name = %[1]q

  response_schema = jsonencode(
            {
               properties = {
                   forecast = {
                       items = {
                           properties = {
                               Day           = {
                                   type = "string"
                                }
                               "Day of week" = {
                                   type = "integer"
                                }
                               Forecast      = {
                                   type = "string"
                                }
                               Humidity      = {
                                   type = "string"
                                }
                               Temperature   = {
                                   type = "integer"
                                }
                               Wind_Speed    = {
                                   type = "integer"
                                }
                            }
                           required   = [
                               "Day",
                               "Forecast",
                               "Humidity",
                               "Temperature",
                               "Wind_Speed",
                               "Day of week",
                            ]
                           type       = "object"
                        }
                       type  = "array"
                    }
                }
               required   = [
                   "forecast",
                ]
               type       = "object"
            }
        )


}
`, model_name)
}
