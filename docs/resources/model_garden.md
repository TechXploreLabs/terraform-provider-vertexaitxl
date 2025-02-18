---
page_title: "vertexaitxl_model_garden Resource - terraform-provider-verteaitxl"
description: |-
  vertexaitxl_model_garden resource generate controlled response for the input prompt w.r.t the response schema. 
---

# vertexaitxl_model_garden Terraform resource

## Overview

The vertexaitxl_model_garden Terraform resource enables controlled AI-generated responses based on a given prompt and a predefined response schema. This allows structured output generation using Vertex AI's foundation models.

This resource is particularly useful when integrating Generative AI models into Infrastructure as Code (IaC) workflows, ensuring deterministic and structured AI-generated configurations.

https://cloud.google.com/vertex-ai/generative-ai/docs/learn/models

## Example Usage

The following example demonstrates how to use the vertexaitxl_model_garden resource to generate a Google Cloud VPC configuration dynamically.

Scenario
Prompt: Request AI to generate a GCP VPC with five subnets across three regions.
Response Schema: Defines a structured output ensuring:
    VPC name is returned.
    Subnet list with names, CIDR ranges, and associated regions.

```terraform
resource "vertexaitxl_model_garden" "gcp-vpc" {
  prompt     = <<EOF
      Create me a vpc in gcp with five subnet in us-central1 , europe-west1 and europe-west2 region with the range of 10.1.0.0/24 , 10.2.0.0/24 and 10.3.0.0/24 respectively.
      cidr ranges should be non overlapping. 
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

output "vpc" {
  value = jsondecode(vertexaitxl_model_garden.vpc.response)
}
```

## Expected output

After applying the Terraform configuration, the AI model generates a structured response based on the provided schema:

```json
{
  "vpc": [
    {
      "name": "vpc-name",
      "subnet": [
        {
          "subnet_name": "subnet-us-central1-0",
          "cidr": "10.1.0.0/24",
          "region": "us-central1"
        },
        {
          "subnet_name": "subnet-europe-west1-0",
          "cidr": "10.2.0.0/24",
          "region": "europe-west1"
        },
        {
          "subnet_name": "subnet-europe-west2-0",
          "cidr": "10.3.0.0/24",
          "region": "europe-west2"
        },
        {
          "subnet_name": "subnet-us-central1-1",
          "cidr": "10.1.1.0/24",
          "region": "us-central1"
        },
        {
          "subnet_name": "subnet-europe-west1-1",
          "cidr": "10.2.1.0/24",
          "region": "europe-west1"
        }
      ]
    }
  ]
}

```

### Required

- `project_id` (String) google cloud platform project id
- `location` (String) location
- `model_name` (String) model name 
- `prompt` (String) input prompt
- `response_schema` (String) response schema for controlled output generation https://cloud.google.com/vertex-ai/docs/reference/rest/v1/projects.locations.cachedContents#Schema

### Read-Only

- `response` (String) Example identifier
