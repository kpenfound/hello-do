terraform {
  cloud {
    organization = "kpenfound"

    workspaces {
      name = "hello-do"
    }
  }

  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "2.23.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

# Create a new container registry
resource "digitalocean_container_registry" "reg" {
  name                   = var.name
  subscription_tier_slug = "starter"
  region = var.do_region
}

# Create DO app
resource "digitalocean_app" "hello" {
  spec {
    name   = var.name
    region = var.do_region

    service {
      name               = var.name
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "basic-xxs"

      image {
        registry_type = "DOCR"
        # This is a bit weird.  The image and tag have to exist to create the app
        repository = var.image_name
        deploy_on_push {
          enabled = true
        }
      }
    }
  }
  depends_on = [
    digitalocean_container_registry.reg
  ]
}
