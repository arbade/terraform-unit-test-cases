variable "aws_region" {
  default = "eu-central-1"
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "table_name" {
  description = "The name to set for the dynamoDB table."
  type        = string
  default     = "terratest-example"
}

variable "main_vpc_cidr" {
  description = "The CIDR of the main VPC"
  type        = string
}

variable "public_subnet_cidr" {
  description = "The CIDR of public subnet"
  type        = string
}

variable "private_subnet_cidr" {
  description = "The CIDR of the private subnet"
  type        = string
}

variable "tag_name" {
  description = "A name used to tag the resource"
  type        = string
  default     = "terraform-network-example"
}