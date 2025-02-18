resource "vertexaitxl_model_garden" "gemini" {
  prompt     = <<EOF
      Create me a vpc in gcp with five subnet in us-central1 , europe-west1 and europe-west2 region with the range of 10.1.0.0/24 , 10.2.0.0/24 and 10.3.0.0/24 respectively.
      cidr rangese should be non overlapping. 
    EOF
  project_id = "my-project"
  location   = "us-central1"
  model_name = "gemini-1.5-pro-002"

  response_schema = jsonencode(
    {
      properties = {
        vpc = {
          items = {
            properties = {
              name = {
                type        = "string"
                description = "Name of the vpc"
              }
              subnet = {
                type = "array"
                items = {
                  type = "object"
                  properties = {
                    subnet_name = {
                      type        = "string"
                      description = "name of the subnet"
                    }
                    cidr = {
                      type = "string"
                    }
                    region = {
                      type = "string"
                    }
                  }
                  required = ["subnet_name", "cidr", "region"]
                }
              }
            }
            required = [
              "name",
              "subnet",
            ]
            type = "object"
          }
          type = "array"
        }
      }
      required = [
        "vpc",
      ]
      type = "object"
    }
  )
}
