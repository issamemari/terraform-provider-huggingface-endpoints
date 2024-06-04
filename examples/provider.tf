terraform {
  required_providers {
    huggingface = {
      source = "issamemari/huggingface-endpoints"
    }
  }
}

provider "huggingface" {
  host      = "https://api.endpoints.huggingface.cloud/v2/endpoint"
  namespace = "issamemari"
  token     = ""
}
