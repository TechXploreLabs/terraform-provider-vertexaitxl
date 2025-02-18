terraform {
  required_providers {
    vertexaitxl = {
      source = "TechXploreLabs/vertexaitxl"
    }
  }
}

provider "vertexaitxl" {
  credentials = file("path/to/serviceaccount.json")
}
