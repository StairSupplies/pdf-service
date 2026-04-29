output "app_name" {
  description = "The name of the pdf-service deployment"
  value       = var.app_name
}

output "app_namespace" {
  description = "The Kubernetes namespace of the pdf-service deployment"
  value       = var.app_namespace
}
