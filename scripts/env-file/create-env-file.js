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
  const testStageStackNames = new Set(
    testStage.actions.map((ts) => ts.configuration.StackName)
  );

  const cfClient = new cf.CloudFormationClient();
  for (const stackName of testStageStackNames) {
    const cmd = new cf.DescribeStacksCommand({
      StackName: stackName,
    });
    const response = await cfClient.send(cmd);
    return response.Stacks[0];
  }

  console.log("Test stage:", testStage);
  // for (const action of testStage.actions) {
  //   // Assuming action configuration contains the stack name
  //   // This part depends on how your pipeline is set up
  //   const stackName = action.configuration.StackName;
  //   if (stackName) {
  //     const resources = await getStackResources(stackName);
  //     console.log(`Resources in stack ${stackName}:`, resources);
  //     // If you need to filter for sub-stacks, do so here by checking resource type
  //   }
  // }

  // const subStacks = (await listSubStacks(pipelineName)).filter((s) =>
  //   s.StackName.includes("AppSync")
  // );
  // const client = new cf.CloudFormationClient();

  // const subStackDescs = subStacks.map(async (stack) => {
  //   const cmd = new cf.DescribeStacksCommand({
  //     StackName: stack.StackName,
  //   });
  //   const output = await client.send(cmd);
  //   return output.Stacks[0];
  // });

  // for (const desc of subStackDescs) {
  //   const outputs = desc.Outputs;
  //   let envContent = "";
  //   outputs.forEach((output) => {
  //     envContent += `${desc.StackName}_${output.OutputKey}=${output.OutputValue}\n`;
  //   });
  //   fs.writeFileSync(".env", envContent);
  // }

  // client.send(cmd).then((data) => {
  //   const pipelineStack = data.Stacks[0];
  //   listSubStacks(pipelineStack.StackName).then((subStacks) => {
  //     for (let stack of subStacks) {
  //       const outputs = stack.Outputs;
  //       let envContent = "";
  //       outputs.forEach((output) => {
  //         envContent += `${stack.}_${camelToSnakeCase(output.OutputKey)}=${
  //           output.OutputValue
  //         }\n`;
  //       });
  //       fs.writeFileSync(".env", envContent);
  //     }
  //   });
  // });
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

async function listSubStacks(stackName) {
  const client = new cf.CloudFormationClient();
  const cmd = new cf.ListStackResourcesCommand({
    StackName: stackName,
  });

  try {
    const data = await client.send(cmd);
    const subStacks = data.StackResourceSummaries.filter(
      (resource) => resource.ResourceType === "AWS::CloudFormation::Stack"
    );
    return subStacks;
  } catch (error) {
    console.error("Error listing sub-stacks:", error);
    throw error;
  }
}

// function camelToSnakeCase(str) {
//   return str
//     .replace(/[A-Z]/g, (letter, index) => `${index > 0 ? "_" : ""}${letter}`)
//     .toUpperCase();
// }
