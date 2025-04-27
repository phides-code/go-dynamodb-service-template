curl -H "x-api-key: d963be82-49f5-4298-8cb0-93a7c9cc6b64" https://9qji5sgtcj.execute-api.ca-central-1.amazonaws.com/Prod/bananas |jq .

curl -X POST -H "Content-Type: application/json" -d @post-data.json -H "X-API-KEY: d963be82-49f5-4298-8cb0-93a7c9cc6b64" https://9qji5sgtcj.execute-api.ca-central-1.amazonaws.com/Prod/bananas |jq .

curl -X PUT -H "Content-Type: application/json" -d @put-data.json -H "X-API-KEY: d963be82-49f5-4298-8cb0-93a7c9cc6b64" https://9qji5sgtcj.execute-api.ca-central-1.amazonaws.com/Prod/bananas/bde8a0d4-af73-4e84-a21a-a9a1b98b5236 |jq .

curl -H "x-api-key: d963be82-49f5-4298-8cb0-93a7c9cc6b64" https://9qji5sgtcj.execute-api.ca-central-1.amazonaws.com/Prod/bananas/bde8a0d4-af73-4e84-a21a-a9a1b98b5236 |jq .

curl -X DELETE -H "x-api-key: d963be82-49f5-4298-8cb0-93a7c9cc6b64" https://9qji5sgtcj.execute-api.ca-central-1.amazonaws.com/Prod/bananas/bde8a0d4-af73-4e84-a21a-a9a1b98b5236 |jq .
