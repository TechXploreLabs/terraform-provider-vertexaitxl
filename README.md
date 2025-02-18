# Terraform Provider Vertexaitxl

To use provider, declare the provider as a required provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    vertexaitxl = {
      source = "TechXploreLabs/vertexaitxl"
    }
  }
}

provider "vertexaitxl" {
    credentials = "/path/to/serviceaccount.json"  # Optional
}

```
## Resource description

1. vertexaitxl_model_garden - use input prompt and response schema to generate controlled output. Model that has controlled generation feature are recommended to use.
