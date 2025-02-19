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

## Example

### Scenario for input config generation:

The following example demonstrates how to use the vertexaitxl_model_garden resource to generate a Google Cloud VPC configuration dynamically.

Scenario
Prompt: Request AI to generate a GCP VPC with five subnets across three regions.
Response Schema: Defines a structured output ensuring:
    VPC name is returned.
    Subnet list with names, CIDR ranges, and associated regions.

```terraform
terraform {
  required_providers {
    vertexaitxl = {
      source = "TechXploreLabs/vertexaitxl"
    }
  }
}

provider "vertexaitxl" {
    credentials = "/path/to/serviceaccount.json"  # Alternatively, "gcloud auth application-default login" can be used 
}


resource "vertexaitxl_model_garden" "vpc-config" {
  prompt     = <<EOF
      Create me a vpc in gcp with five subnet in us-central1 , europe-west1 and europe-west2 region with the range of 10.1.0.0/24 , 10.2.0.0/24 and 10.3.0.0/24 respectively.
      cidr ranges should be non overlapping. 
    EOF
  project_id = "my-project-id"
  location   = "us-central1"
  model_name = "gemini-1.5-pro-002"

  response_schema = jsonencode(
    {
      type     = "object"
      required = ["vpc"]
      properties = {
        vpc = {
          type = "object"
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
        }
      }
    }
  )
}


resource "google_compute_network" "vpc" {
  name                    = jsondecode(vertexaitxl_model_garden.vpc-config.response).vpc["name"]
  auto_create_subnetworks = false
  project                 = "my-project-id"
}


resource "google_compute_subnetwork" "subnet" {
  name          = jsondecode(vertexaitxl_model_garden.vpc-config.response).vpc.subnet[0].subnet_name
  ip_cidr_range = jsondecode(vertexaitxl_model_garden.vpc-config.response).vpc.subnet[0].cidr
  region        = jsondecode(vertexaitxl_model_garden.vpc-config.response).vpc.subnet[0].region
  network       = google_compute_network.vpc.id
  project       = "my-project-id"
}

output "vpc-config" {
  value = jsondecode(vertexaitxl_model_garden.vpc-config.response)
}
```

## Expected output

After applying the Terraform configuration, the AI model generates a structured response based on the provided schema:

```json
{
  "vpc" = {
    "name" = "my-vpc"
    "subnet" = [
      {
        "cidr" = "10.1.0.0/24"
        "region" = "us-central1"
        "subnet_name" = "subnet-us-central1-01"
      },
      {
        "cidr" = "10.2.0.0/24"
        "region" = "europe-west1"
        "subnet_name" = "subnet-europe-west1-01"
      },
      {
        "cidr" = "10.3.0.0/24"
        "region" = "europe-west2"
        "subnet_name" = "subnet-europe-west2-01"
      },
      {
        "cidr" = "10.1.1.0/24"
        "region" = "us-central1"
        "subnet_name" = "subnet-us-central1-02"
      },
      {
        "cidr" = "10.2.1.0/24"
        "region" = "europe-west1"
        "subnet_name" = "subnet-europe-west1-02"
      },
    ]
  }
}

```

### Scenario for gaurdrails:

The following example demonstrates how to use the vertexaitxl_model_garden resource to check gaurdrails.

Scenario
Prompt: Request AI to check whether count of label variable is 4.
Response Schema: Defines a structured output ensuring:
    label_count for getting the count.
    label_match for boolean output.
  
```terraform
variable "label" {
  type = list(map(string))
  default = [{
    "env"        = "dev"
    "billing"    = "org"
    "owner"      = "forbar"
    "department" = "finance"
    }, {
    "env"        = "test"
    "billing"    = "org"
    "department" = "finance"
  }]
}

resource "vertexaitxl_model_garden" "check-config" {
  count = length(var.label)
  prompt = <<EOF
        Check below variable as 4 labels
        ${jsonencode(var.label[count.index]
)}
    EOF
project_id = "my-project-id"
location   = "us-central1"
model_name = "gemini-1.5-pro-002"

response_schema = jsonencode(
  {
    type = "object"
    properties = {
      label_match = {
        type = "bool"
      }
      label_count = {
        type        = "integer"
        description = "count of the keys"
      }
    }
  }
)
}


output "response" {
  value = [jsondecode(vertexaitxl_model_garden.check-config[0].response),
  jsondecode(vertexaitxl_model_garden.check-config[1].response)]
}
```

## Expected output

After applying the Terraform configuration, the AI model generates a structured response based on the provided schema:

```json
response = [
  {
    "label_count" = 4
    "label_match" = true
  },
  {
    "label_count" = 3
    "label_match" = false
  },
]
```

### Required

- `project_id` (String) google cloud platform project id
- `location` (String) location
- `model_name` (String) model name 
- `prompt` (String) input prompt
- `response_schema` (String) response schema for controlled output generation https://cloud.google.com/vertex-ai/docs/reference/rest/v1/projects.locations.cachedContents#Schema

### Read-Only

- `response` (String) Example identifier
