output "docr_repo" {
  value = digitalocean_container_registry.reg.endpoint
}

output "app_live_url" {
  value = digitalocean_app.hello.live_url
}
