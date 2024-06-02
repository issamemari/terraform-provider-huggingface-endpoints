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

resource "huggingface_endpoint" "edu" {
  name = "test-endpoint-issa"

  compute = {
    accelerator  = "cpu"
    instance_size = "x8"
    instance_type = "intel-icl"
    scaling = {
      min_replica = 0
      max_replica = 1
      scale_to_zero_timeout = 15
    }
  }

  model = {
    framework = "pytorch"
    image = {
      huggingface = {
        env = {}
      }
    }
    repository = "sentence-transformers/all-MiniLM-L6-v2"
    revision   = "main"
    task       = "sentence-embeddings"
  }

  provider_details = {
    region = "us-east-1"
    vendor = "aws"
  }

  type = "protected"
}

output "edu_endpoint" {
  value = huggingface_endpoint.edu
}
