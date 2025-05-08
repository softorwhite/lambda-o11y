
output "lambda_iam_role_arn" {
  description = "IAM Role ARN for the Lambda function"
  value       = aws_iam_role.this.arn
}