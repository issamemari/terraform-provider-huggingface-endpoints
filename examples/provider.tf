terraform {
  required_providers {
    huggingface = {
      source = "example.com/issamemari/huggingface"
    }
  }
}

provider "huggingface" {
  host      = "https://api.endpoints.huggingface.cloud/v2/endpoint"
  namespace = "issamemari"
  token     = ""
}
