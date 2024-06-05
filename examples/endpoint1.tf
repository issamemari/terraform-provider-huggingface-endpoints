resource "huggingface_endpoint" "endpoint1" {
  name = "test-endpoint-issa1"

  compute = {
    accelerator   = "cpu"
    instance_size = "x8"
    instance_type = "intel-icl"
    scaling = {
      min_replica           = 0
      max_replica           = 1
      scale_to_zero_timeout = 15
    }
  }

  model = {
    framework = "pytorch"
    image = {
      # custom = {
      #   url          = "ghcr.io/huggingface/text-embeddings-inference:cpu-0.6.0"
      #   health_route = "/health"
      #   env          = {
      #     MAX_BATCH_TOKENS        = 1000000
      #     MAX_CONCURRENT_REQUESTS = 512
      #     MODEL_ID                = "/repository"
      #   }
      # }
      huggingface = {
        env = {}
      }
    }
    repository = "sentence-transformers/all-MiniLM-L6-v2"
    task       = "sentence-embeddings"
  }

  cloud = {
    region = "us-east-1"
    vendor = "aws"
  }

  type = "protected"
}

output "endpoint1" {
  value = huggingface_endpoint.endpoint1
}
