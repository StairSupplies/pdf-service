variable "app_name" {
  description = "Name of the deployment and service"
  type        = string
}

variable "app_namespace" {
  description = "Kubernetes namespace to deploy into"
  type        = string
}

variable "image_repo" {
  description = "Docker image repository (e.g. stairsupplies/pdf-service)"
  type        = string
}

variable "image_tag" {
  description = "Docker image tag"
  type        = string
}

variable "image_pull_secret" {
  description = "Kubernetes image pull secret name"
  type        = string
  default     = "dockerhub"
}

variable "replica_count" {
  description = "Number of pod replicas"
  type        = number
  default     = 1
}

variable "requests_memory" {
  description = "Memory request per pod — validate against real load before increasing"
  type        = string
  default     = "128Mi"
}

variable "requests_cpu" {
  description = "CPU request per pod — validate against real load before increasing"
  type        = string
  default     = "250m"
}

variable "limits_memory" {
  description = "Memory limit per pod — covers Go process (~30MB) + pdflatex peak (~200MB); raise if OOMKilled"
  type        = string
  default     = "512Mi"
}

variable "limits_cpu" {
  description = "CPU limit per pod — pdflatex is single-threaded; validate against real load"
  type        = string
  default     = "500m"
}

variable "app_env" {
  description = "APP_ENV passed to the container (e.g. production)"
  type        = string
  default     = "production"
}

variable "app_port" {
  description = "Port the service listens on"
  type        = number
  default     = 8080
}

variable "api_key" {
  description = "Bearer token for PDF_SERVICE_API_KEY"
  type        = string
  sensitive   = true
}

variable "pdflatex_timeout" {
  description = "Per-request pdflatex timeout in seconds (PDF_SERVICE_PDFLATEX_TIMEOUT)"
  type        = string
  default     = "55"
}

variable "write_timeout" {
  description = "HTTP server write timeout in seconds (WRITE_TIMEOUT); must exceed pdflatex_timeout"
  type        = string
  default     = "120"
}

variable "image_pull_policy" {
  description = "Image pull policy — use IfNotPresent for local k3d testing (imported images)"
  type        = string
  default     = "Always"
}
