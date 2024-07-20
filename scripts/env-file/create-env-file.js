const fs = require("fs");
const process = require("process");
const cf = require("@aws-sdk/client-cloudformation");
const cp = require("@aws-sdk/client-codepipeline");

execute();

async function execute() {
  const fileName = ".env";
  console.log(`Creating ${fileName} file...`);

  const pipelineName = `appsyncmasterclass_pipeline`;
  const stages = await getPipelineStages(pipelineName);
  const testStage = stages.find((stage) => stage.name === "TEST");
  if (!testStage) {
    throw new Error("TEST stage not found in the pipeline.");
  }
  // Get unique stack names of the test stage
  const testStageStackNames = new Set(
    testStage.actions.map((ts) => ts.configuration.StackName)
  );

  const outputs = await getOutputs(testStageStackNames);

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

async function getPipelineStages(pipelineName) {
  const codePipelineClient = new cp.CodePipelineClient({
    region: process.env.AWS_DEFAULT_REGION,
  });
  const command = new cp.GetPipelineCommand({ name: pipelineName });
  const response = await codePipelineClient.send(command);
  return response.pipeline.stages;
}

function camelToSnakeCase(str) {
  return str
    .replace(/[A-Z]/g, (letter, index) => `${index > 0 ? "_" : ""}${letter}`)
    .toUpperCase();
}
