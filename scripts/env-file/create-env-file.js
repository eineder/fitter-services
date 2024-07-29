const process = require("process");
const fs = require("fs");
const cf = require("@aws-sdk/client-cloudformation");

execute();

async function execute() {
  let stage = "TEST";
  if (process.argv.length > 2) stage = process.argv[2];
  const fileName = `.${stage}.env`;
  console.log(`Creating ${fileName} file...`);

  const client = new cf.CloudFormationClient();
  const lsc = new cf.ListStacksCommand();
  const response = await client.send(lsc);
  stackNames = response.StackSummaries.filter((stack) =>
    stack.StackName.startsWith(`${stage}-fitter-services`)
  ).map((stack) => stack.StackName);
  const outputs = await getOutputs(stackNames);

  const envs = outputs.map((output) => {
    const key = camelToSnakeCase(output.OutputKey);
    const value = output.OutputValue;
    return { key, value };
  });
  envs.push({ key: "AWS_SDK_LOAD_CONFIG", value: "1" });
  const envContent = envs
    .sort((a, b) => a.key.localeCompare(b.key))
    .reduce((acc, env) => {
      return `${acc}${env.key}=${env.value}\n`;
    }, "");

  fs.writeFileSync(fileName, envContent);
  console.log(`${fileName} file created.`);
}
async function getOutputs(testStageStackNames) {
  const client = new cf.CloudFormationClient();
  const promises = [];
  for (const stackName of testStageStackNames) {
    const cmd = new cf.DescribeStacksCommand({
      StackName: stackName,
    });
    const promise = client.send(cmd);
    promises.push(promise);
  }
  const responses = await Promise.all(promises);
  const outputs = responses.flatMap((response) => response.Stacks[0].Outputs);
  return outputs.filter((output) => output !== undefined);
}

function camelToSnakeCase(str) {
  return str
    .replace(/[A-Z]/g, (letter, index) => `${index > 0 ? "_" : ""}${letter}`)
    .toUpperCase();
}
