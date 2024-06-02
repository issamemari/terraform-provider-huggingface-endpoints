terraform {
  required_providers {
    huggingface = {
      source = "hashicorp.com/edu/huggingface"
    }
  }
}

provider "huggingface" {
  host      = "https://api.endpoints.huggingface.cloud/v2/endpoint"
  namespace = "issamemari"
  token     = ""
}
