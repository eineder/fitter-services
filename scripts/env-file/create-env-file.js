const fs = require("fs");
const cf = require("@aws-sdk/client-cloudformation");
const cp = require("@aws-sdk/client-codepipeline");
const { execSync } = require("child_process");

execute();

async function execute() {
  console.log("Creating .env file...");

  const pipelineName = `appsyncmasterclass_${getBranchName()}_pipeline`;
  const stages = await getPipelineStages(pipelineName);
  const testStage = stages.find((stage) => stage.name === "TEST");
  // Get unique stack names of the test stage
  const testStageStackNames = new Set(
    testStage.actions.map((ts) => ts.configuration.StackName)
  );

  const outputs = await getOutputs(testStageStackNames);

  let envContent = "";
  outputs.forEach((output) => {
    const key = camelToSnakeCase(output.OutputKey);
    const value = output.OutputValue;
    envContent += `${key}=${value}\n`;
  });
  fs.writeFileSync(".env", envContent);
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

function getBranchName() {
  try {
    const branch = execSync("git rev-parse --abbrev-ref HEAD")
      .toString()
      .trim();
    return branch;
  } catch (error) {
    console.error("Error getting git branch:", error);
    throw error;
  }
}
async function getPipelineStages(pipelineName) {
  const codePipelineClient = new cp.CodePipelineClient();
  const command = new cp.GetPipelineCommand({ name: pipelineName });
  const response = await codePipelineClient.send(command);
  return response.pipeline.stages;
}

function camelToSnakeCase(str) {
  return str
    .replace(/[A-Z]/g, (letter, index) => `${index > 0 ? "_" : ""}${letter}`)
    .toUpperCase();
}
