# hello-do

This repo deploys a simple web service as a Digital Ocean App.

## Terraform

The required resources are created by the terraform code under `terraform/`.

Once "gotcha" here is that we want to create both the DOCR (Digital Ocean Container Registry) to host our image, as well as the DO App that runs the image. Creating the App requires an image to be present in the container registry, so the resources must be created in two separate applies.

Note the outputs.

- `docr_repo` is where we will push our images
- `app_live_url` is our app's endpoint

## Dagger

The Dagger pipeline in `ci/main.go` builds our image and pushes it to DOCR.

1. Authenticate docker to DOCR: `docker login registry.digitalocean.com`
2. Run the Dagger pipeline: `go run ci/main.go`
3. Once the pipeline is complete, the DO App will detect a new image at `:latest` and automatically deploy the new version of our image
4. Send a request to the `app_live_url` from above: `curl $app_live_url`. You should get a successful response `Hello`
