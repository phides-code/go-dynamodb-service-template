# go-dynamodb-service

A Go project template which provides CRUD services for a DynamoDB table, using AWS Lambda and API Gateway, deployed with AWS SAM and GitHub Actions.

-   Find and replace `Appname` with the name of your app (upper and lowercase A)
-   Find and replace `Banana`/`Bananas` with the table name (upper and lowercase B)
-   Find and replace `ca-central-1` with your AWS region
-   Update fields in `database.go`
-   Update OriginURL in `constants.go`
-   Add AWS secrets to GitHub Actions for this repo:
    -   `gh secret set AWS_ACCESS_KEY_ID`
    -   `gh secret set AWS_SECRET_ACCESS_KEY`
    -   `gh run rerun`
-   Run: `make deploy`
-   Get API ID with: `aws apigateway get-rest-apis |grep \"id\" -A1`
-   API URL: `https://my-api-id.execute-api.ca-central-1.amazonaws.com/Prod/my-table-name`

-   To run it locally: `make build && sam local start-api --port 8000`
