resource "huggingface_endpoint" "endpoint2" {
  name = "test-endpoint-issa2"

  compute = {
    accelerator   = "cpu"
    instance_size = "x8"
    instance_type = "intel-icl"
    scaling = {
      min_replica           = 0
      max_replica           = 2
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

  cloud = {
    region = "us-east-1"
    vendor = "aws"
  }

  type = "protected"
}

output "endpoint2" {
  value = huggingface_endpoint.endpoint2
}
