resource "kubernetes_secret" "pdf_service_api_key" {
  metadata {
    name      = "${var.app_name}-secret"
    namespace = var.app_namespace
  }
  data = {
    api_key = var.api_key
  }
}

resource "kubernetes_manifest" "pdf_service_deployment" {
  manifest = yamldecode(templatefile("${path.module}/manifests/pdf-service-deployment.yaml.tftpl", {
    app_name          = var.app_name
    app_namespace     = var.app_namespace
    replica_count     = var.replica_count
    image_repo        = var.image_repo
    image_tag         = var.image_tag
    image_pull_secret  = var.image_pull_secret
    image_pull_policy  = var.image_pull_policy
    requests_memory   = var.requests_memory
    requests_cpu      = var.requests_cpu
    limits_memory     = var.limits_memory
    limits_cpu        = var.limits_cpu
    app_env           = var.app_env
    app_port          = var.app_port
    secret_name       = "${var.app_name}-secret"
    pdflatex_timeout  = var.pdflatex_timeout
    write_timeout     = var.write_timeout
  }))

  field_manager {
    name            = "pdf-service-deployment-field-manager"
    force_conflicts = true
  }
}

resource "kubernetes_manifest" "pdf_service_cluster_ip" {
  manifest = yamldecode(templatefile("${path.module}/manifests/pdf-service-service.yaml.tftpl", {
    service_name = var.app_name
    namespace    = var.app_namespace
    http_port    = var.app_port
    app_label    = var.app_name
  }))

  field_manager {
    name            = "pdf-service-cluster-ip-field-manager"
    force_conflicts = true
  }
}
