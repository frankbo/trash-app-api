provider "aws" {
  region = "eu-central-1"
}

resource "aws_lambda_function" "trash_api_lambda" {
  filename      = "${path.cwd}/../function.zip"
  function_name = "trash-api"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "trash-app-api"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("${path.cwd}/../function.zip")

  runtime = "go1.x"
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

  # request_parameters = {
  #   "integration.request.querystring.locationId" = "method.request.querystring.locationId"
  #   "integration.request.querystring.streetId"   = "method.request.querystring.streetId"
  # }
}

# resource "aws_api_gateway_method" "proxy_root" {
#   rest_api_id   = aws_api_gateway_rest_api.trash_api_gateway.id
#   resource_id   = aws_api_gateway_rest_api.trash_api_gateway.root_resource_id
#   http_method   = "GET"
#   authorization = "NONE"

#   request_parameters = {
#     "method.request.querystring.locationId" = true
#     "method.request.querystring.streetId"   = false
#   }
# }

# resource "aws_api_gateway_integration" "lambda_root" {
#   rest_api_id = aws_api_gateway_rest_api.trash_api_gateway.id
#   resource_id = aws_api_gateway_method.proxy_root.resource_id
#   http_method = aws_api_gateway_method.proxy_root.http_method

#   integration_http_method = "POST"
#   type                    = "AWS_PROXY"
#   uri                     = aws_lambda_function.trash_api_lambda.invoke_arn

#   request_parameters = {
#     "integration.request.querystring.locationId" = "method.request.querystring.locationId"
#     "integration.request.querystring.streetId"   = "method.request.querystring.streetId"
#   }
# }

resource "aws_api_gateway_deployment" "trash_api_gw_deploy" {
  depends_on = [
    aws_api_gateway_integration.lambda,
  ]

  rest_api_id = aws_api_gateway_rest_api.trash_api_gateway.id
  stage_name  = "test"
}

resource "aws_api_gateway_rest_api" "trash_api_gateway" {
  name        = "trash-api"
  description = "Trash App Backend"
}

output "base_url" {
  value = aws_api_gateway_deployment.trash_api_gw_deploy.invoke_url
}
