export AWS_PAGER=""

npm install scripts/env-file
node ./scripts/env-file/create-env-file.js
echo "Created .env file"

source .env

aws lambda invoke \
    --function-name $PRIME_SWEARWORDS_FUNCTION_NAME \
    --payload '{}' \
    --invocation-type Event \
    /dev/null
echo "Invoked lambda function to seed database: '$PRIME_SWEARWORDS_FUNCTION_NAME'"

go test ./...
echo "Ran tests"