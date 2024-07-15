# Abort as soon as a command fails
set -e

export AWS_PAGER=""

cd scripts/env-file
npm ci
cd ../..
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