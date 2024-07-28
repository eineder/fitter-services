# Make sure that no page is used when invoking AWS CLI commands
export AWS_PAGER=""

echo "**** Installing AWS CDK... ****"
npm install -g aws-cdk || exit 1
echo "**** ✅ Successfully installed AWS CDK ****"

echo "**** Synthesize... ****"
cdk synth
echo "**** ✅ Successfully synthesized ****"

echo "**** Deploying TEST stage... ****"
cdk deploy "TEST/*" --require-approval never || exit 1
echo "**** ✅ Successfully deployed TEST stage ****"

echo "**** Creating .env file for TEST... ****"
cd scripts/env-file
npm ci || exit 1
cd ../..
node ./scripts/env-file/create-env-file.js TEST|| exit 1
echo "**** ✅ Successfully created .env file for TEST ****"

echo "**** Loading TEST environment variables to be available in the following commands... ****"
source .TEST.env || exit 1
echo "**** ✅ Successfully loaded TEST environment variables ****"

echo "**** Invoking lambda function to seed TEST database... ****"
aws lambda invoke \
    --function-name $PRIME_SWEARWORDS_FUNCTION_NAME \
    --payload '{}' \
    --invocation-type RequestResponse \
    lambda-out.json || exit 1
echo "**** ✅ Successfully invoked lambda function to seed TEST database... ****"

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

echo "**** Deploying PROD stage... ****"
cdk deploy "PROD/*" --require-approval never || exit 1
echo "**** ✅ Successfully deployed PROD stage ****"

echo "**** Creating .env file for PROD... ****"
cd scripts/env-file
npm ci || exit 1
cd ../..
node ./scripts/env-file/create-env-file.js PROD|| exit 1
echo "**** ✅ Successfully created .env file for PROD ****"

echo "**** Loading PROD environment variables to be available in the following commands... ****"
source .PROD.env || exit 1
echo "**** ✅ Successfully loaded PROD environment variables ****"

echo "**** Invoking lambda function to seed PROD database... ****"
aws lambda invoke \
    --function-name $PRIME_SWEARWORDS_FUNCTION_NAME \
    --payload '{}' \
    --invocation-type RequestResponse \
    lambda-out.json || exit 1
echo "**** ✅ Successfully invoked lambda function to seed PROD database... ****"

# Check if the output is null - exit if not
if ! grep -q "null" lambda-out.json; then
  echo "**** ❌ Lambda function output is not null as expected: ****"
  cat lambda-out.json
  echo "**** Exiting... ****"
  exit 1
fi
echo "**** ✅ Successfully seeded table $SWEARWORDS_TABLE_NAME ****"
