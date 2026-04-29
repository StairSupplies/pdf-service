terraform {
  required_version = "~> 1.0"

  backend "s3" {
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

locals {
  cluster_name      = "terminal-eks-cluster"
  region            = var.region
  app_namespace     = var.app_namespace
  app_name          = var.app_name
  image_repo        = var.image_repo
  image_tag         = var.image_tag
  image_pull_secret = var.image_pull_secret
  replica_count     = var.replica_count
  requests_memory   = var.requests_memory
  requests_cpu      = var.requests_cpu
  limits_memory     = var.limits_memory
  limits_cpu        = var.limits_cpu
  app_env           = var.app_env
  app_port          = var.app_port
  api_key           = var.api_key
  pdflatex_timeout  = var.pdflatex_timeout
  write_timeout     = var.write_timeout
}

provider "aws" {
  region = local.region
}

data "aws_eks_cluster" "eks" {
  name = local.cluster_name
}

data "aws_eks_cluster_auth" "eks" {
  name = local.cluster_name
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.eks.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.eks.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.eks.token
}

# ===============================
# PDF Service
# ===============================

module "pdf_service" {
  source = "../../modules/pdf-service"

  app_name          = local.app_name
  app_namespace     = local.app_namespace
  image_repo        = local.image_repo
  image_tag         = local.image_tag
  image_pull_secret = local.image_pull_secret
  replica_count     = local.replica_count
  requests_memory   = local.requests_memory
  requests_cpu      = local.requests_cpu
  limits_memory     = local.limits_memory
  limits_cpu        = local.limits_cpu
  app_env           = local.app_env
  app_port          = local.app_port
  api_key           = local.api_key
  pdflatex_timeout  = local.pdflatex_timeout
  write_timeout     = local.write_timeout
}

output "app_name" {
  value = module.pdf_service.app_name
}

output "app_namespace" {
  value = module.pdf_service.app_namespace
}
