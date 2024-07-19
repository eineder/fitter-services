const fs = require("fs");
const path = require("path");

function getStackOutputs(templateFile) {
  const template = JSON.parse(fs.readFileSync(templateFile, "utf-8"));
  const outputs = template.Outputs || null;
  return outputs;
}

function traverseDirectory(dir, fileCallback) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const filePath = path.join(dir, file);
    if (fs.lstatSync(filePath).isDirectory()) {
      traverseDirectory(filePath, fileCallback);
    } else if (path.extname(file) === ".json") {
      fileCallback(filePath);
    }
  }
}

function camelToSnakeCase(str) {
  return str
    .replace(/[A-Z]/g, (letter, index) => `${index > 0 ? "_" : ""}${letter}`)
    .toUpperCase();
}

function main() {
  const envFileName = ".env";
  console.log(`Creating ${envFileName} file...`);

  const cdkOutDir = "./cdk.out"; // Directory where CDK synthesized templates are stored
  const allOutputs = {};

  traverseDirectory(cdkOutDir, (filePath) => {
    if (!filePath.includes("TEST")) return;
    const stackName = path.relative(cdkOutDir, filePath).replace(".json", "");
    const stackOutputs = getStackOutputs(filePath);
    if (stackOutputs) allOutputs[stackName] = stackOutputs;
  });

  const outputs = [];
  for (let stackName in allOutputs) {
    const stackOutputs = allOutputs[stackName];
    for (let outputName in stackOutputs) {
      const output = stackOutputs[outputName].Value;
      const outputValue = typeof output === "string" ? output : output.Ref;
      outputs.push({
        OutputKey: outputName,
        OutputValue: outputValue,
      });
    }
  }
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

  fs.writeFileSync(envFileName, envContent);
  console.log(`${envFileName} file created.`);
}

main();
