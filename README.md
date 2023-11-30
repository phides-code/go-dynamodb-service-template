# go-dynamodb-service

A Go project template which provides CRUD services for a DynamoDB table, using AWS Lambda and API Gateway, deployed with AWS SAM and GitHub Actions.

-   Find and replace "Banana" with the table name (upper and lowercase)
-   Update fields in database.go
-   Update AmplifyAppId in constants.go
-   Get API ID with: aws apigateway get-rest-apis |grep \"id\" -A1
-   API URL: https://my-api-id.execute-api.ca-central-1.amazonaws.com/Prod/my-table-name
