terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.93.0"
    }
  }

  backend "s3" {
    bucket  = "softorwhite-lambda-o11y"
    region  = "ap-northeast-1"
    key     = "lambda-o11y.tfstate"
    encrypt = true
  }
}


provider "aws" {
  region = "ap-northeast-1"
  default_tags {
    tags = {
      "service"    = "softorwhite-lambda-o11y"
      "repo"       = "lambda-o11y"
      "managed_by" = "terraform"
    }
  }
}


module "lambda_o11y" {
  source                  = "../../modules/lambda-o11y"
  bucket                  = "softorwhite-development-lambda-o11y"
  stage                   = "development"
}


