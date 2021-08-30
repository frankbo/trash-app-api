provider "aws" {
  region = "eu-central-1"
}

terraform {
  backend "s3" {
    bucket = "blb-trash-api-terraform-state"
    key    = "default-infrastructure"
    region = "eu-central-1"
  }
}

resource "aws_s3_bucket" "terraform_state" {
  bucket = "blb-trash-api-terraform-state"

  versioning {
    enabled = true
  }
}

variable "lambda_function_name" {
  default = "trash-api"
}

resource "aws_cloudwatch_log_group" "trash_logs" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 14
}

resource "aws_lambda_function" "trash_api_lambda" {
  filename      = "function.zip"
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "trash-app-api"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("function.zip")

  runtime = "go1.x"

  depends_on = [
    aws_iam_role_policy_attachment.lambda_logs,
    aws_cloudwatch_log_group.trash_logs,
  ]
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ],
          "Resource" : "arn:aws:logs:*:*:*",
          "Effect" : "Allow"
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name  
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_lambda_permission" "apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.trash_api_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/* portion grants access from any method on any resource
  # within the API Gateway "REST API".
  source_arn = "${aws_api_gateway_rest_api.trash_api_gateway.execution_arn}/*/*"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.trash_api_gateway.id
  parent_id   = aws_api_gateway_rest_api.trash_api_gateway.root_resource_id
  path_part   = "events"
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = aws_api_gateway_rest_api.trash_api_gateway.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.querystring.locationId" = true
    "method.request.querystring.streetId"   = false
  }
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.trash_api_gateway.id
  resource_id = aws_api_gateway_method.proxy.resource_id
  http_method = aws_api_gateway_method.proxy.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.trash_api_lambda.invoke_arn
}

resource "aws_api_gateway_deployment" "trash_api_gw_deploy" {
  depends_on = [
    aws_api_gateway_integration.lambda,
  ]

  rest_api_id = aws_api_gateway_rest_api.trash_api_gateway.id
  stage_name  = "v1"
}

resource "aws_api_gateway_rest_api" "trash_api_gateway" {
  name        = "trash-api"
  description = "Trash App Backend"
}

output "base_url" {
  value = aws_api_gateway_deployment.trash_api_gw_deploy.invoke_url
}
