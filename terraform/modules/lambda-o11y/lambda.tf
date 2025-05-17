data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

####################################################################################################
# Lambda用S3
####################################################################################################

resource "aws_s3_bucket" "this" {
  bucket = var.bucket
}

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  bucket = aws_s3_bucket.this.bucket
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket                  = aws_s3_bucket.this.bucket
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

####################################################################################################
# Lambda
####################################################################################################

locals {
  this_s3_zip_path      = "archive/bootstrap.zip"
  this_function_dir     = "${path.module}/../../../app"
  this_function_cmd_dir = "${path.module}/../../../app/cmd/lambda-o11y"
  this_zip_path         = "${local.this_function_cmd_dir}/build/bootstrap.zip"
  this_function_code = join("", [
    for file in fileset(local.this_function_dir, "**")
    : filebase64sha256("${local.this_function_dir}/${file}") if endswith(file, ".go") || endswith(file, "go.mod") || endswith(file, "go.sum")
  ])
}

resource "null_resource" "this" {
  triggers = {
    code_diff = local.this_function_code
  }
  provisioner "local-exec" {
    command = "cd ${path.module}/../../.. && make zip-lambda-o11y"
  }
}

resource "aws_s3_object" "zip" {
  bucket      = aws_s3_bucket.this.bucket
  key         = local.this_s3_zip_path
  source      = local.this_zip_path
  source_hash = local.this_function_code
  depends_on  = [null_resource.this]
}

resource "aws_lambda_function" "this" {
  function_name    = "${var.stage}-softorwhite-lambda-o11y"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = aws_s3_bucket.this.bucket
  s3_key           = aws_s3_object.zip.key
  source_code_hash = local.this_function_code
  role             = aws_iam_role.this.arn
  timeout          = 900
  tracing_config {
    mode = "Active"
  }

  layers           = ["arn:aws:lambda:${data.aws_region.current.name}:901920570463:layer:aws-otel-collector-amd64-ver-0-115-0:3"]
  environment {
    variables = {
      STAGE = var.stage
    }
  }
}


####################################################################################################
# IAM Role
####################################################################################################

resource "aws_iam_role" "this" {
  name = "${var.stage}-lambda-o11y-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = ["sts:AssumeRole"]
        Principal = {
          Service = ["lambda.amazonaws.com"]
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "vpc_access_execution_this" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}


resource "aws_iam_role_policy_attachment" "test_xray" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy_attachment" "hello-lambda-cloudwatch-insights" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}