# Make sure that no page is used when invoking AWS CLI commands
export AWS_PAGER=""

echo "**** Installing AWS CDK... ****"
npm install -g aws-cdk || exit 1
echo "**** ✅ Successfully installed AWS CDK ****"

echo "**** Synthesize... ****"
cdk synth
echo "**** ✅ Successfully synthesized ****"

echo "**** Deploying dev stage... ****"
cdk deploy "dev/*" --require-approval never || exit 1
echo "**** ✅ Successfully deployed dev stage ****"

echo "**** Creating .env file for dev... ****"
cd scripts/env-file
npm ci || exit 1
cd ../..
node ./scripts/env-file/create-env-file.js dev|| exit 1
echo "**** ✅ Successfully created .env file for dev ****"

echo "**** Loading dev environment variables to be available in the following commands... ****"
source .dev.env || exit 1
echo "**** ✅ Successfully loaded dev environment variables ****"

echo "**** Invoking lambda function to seed dev database... ****"
aws lambda invoke \
    --function-name $PRIME_SWEARWORDS_FUNCTION_NAME \
    --payload '{}' \
    --invocation-type RequestResponse \
    lambda-out.json || exit 1
echo "**** ✅ Successfully invoked lambda function to seed dev database... ****"

# Check if the output is null - exit if not
if ! grep -q "null" lambda-out.json; then
  echo "**** ❌ Lambda function output is not null as expected: ****"
  cat lambda-out.json
  echo
  echo "**** Exiting... ****"
  exit 1
fi
echo "**** ✅ Successfully seeded table $SWEARWORDS_TABLE_NAME ****"


echo "**** Running tests... ****"
go test -count=1 ./... || exit 1
# -count=1 to avoid caching
echo "**** ✅ Successfully ran tests ****"

echo "**** Deploying prod stage... ****"
cdk deploy "prod/*" --require-approval never || exit 1
echo "**** ✅ Successfully deployed prod stage ****"

echo "**** Creating .env file for prod... ****"
cd scripts/env-file
npm ci || exit 1
cd ../..
node ./scripts/env-file/create-env-file.js prod|| exit 1
echo "**** ✅ Successfully created .env file for prod ****"

echo "**** Loading prod environment variables to be available in the following commands... ****"
source .prod.env || exit 1
echo "**** ✅ Successfully loaded prod environment variables ****"

echo "**** Invoking lambda function to seed prod database... ****"
aws lambda invoke \
    --function-name $PRIME_SWEARWORDS_FUNCTION_NAME \
    --payload '{}' \
    --invocation-type RequestResponse \
    lambda-out.json || exit 1
echo "**** ✅ Successfully invoked lambda function to seed prod database... ****"

# Check if the output is null - exit if not
if ! grep -q "null" lambda-out.json; then
  echo "**** ❌ Lambda function output is not null as expected: ****"
  cat lambda-out.json
  echo "**** Exiting... ****"
  exit 1
fi
echo "**** ✅ Successfully seeded table $SWEARWORDS_TABLE_NAME ****"
