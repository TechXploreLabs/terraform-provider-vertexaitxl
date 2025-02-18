terraform {
  required_providers {
    vertexaitxl = {
      source = "TechXploreLabs/vertexaitxl"
    }
  }
}

provider "vertexaitxl" {
  credentials = "path/to/serviceaccount.json"
}
