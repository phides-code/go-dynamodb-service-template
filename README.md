# go-dynamodb-service

A Go project template which provides CRUD services for a DynamoDB table, using AWS Lambda and API Gateway, deployed with AWS SAM and GitHub Actions.

-   Find and replace "Banana"/"Bananas" with the table name (upper and lowercase)
-   Update fields in database.go
-   Update OriginURL in constants.go
-   Add AWS secrets to GitHub Actions for this repo:
    -   gh secret set AWS_ACCESS_KEY_ID
    -   gh secret set AWS_SECRET_ACCESS_KEY
    -   gh run rerun
-   Get API ID with: aws apigateway get-rest-apis |grep \"id\" -A1
-   API URL: https://my-api-id.execute-api.ca-central-1.amazonaws.com/Prod/my-table-name
