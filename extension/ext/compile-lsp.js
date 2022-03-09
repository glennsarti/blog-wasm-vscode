const fs = require('fs');
const path = require('path');
const cp = require('child_process');

const lspDir = path.join(__dirname, "..", "..", "lsp");
const outPath = path.join(lspDir, "pkg");

if (!fs.existsSync(outPath)) {
  console.log('Creating output directory');
  fs.mkdirSync(outPath, { recursive: true });
}

function getEnv() {
  // Dodgey but effective!
  let thisEnv = JSON.parse(JSON.stringify(process.env));
  return thisEnv;
}

console.log("Building for WASM...");
let execEnv = getEnv();
execEnv["GOOS"] = "js";
execEnv["GOARCH"] = "wasm";
let result = cp.execSync("go build -o \"" + path.join(outPath, "demo-ls.wasm") + "\"", {
  cwd: lspDir,
  env: execEnv
});
console.log(result.toString());

console.log("Building for Node...");
execEnv = getEnv();
let outFile = path.join(outPath, "demo-ls");
if (process.platform === "win32") {
  outFile = outFile + ".exe";
}
result = cp.execSync("go build -o \"" + outFile + "\"", {
  cwd: lspDir,
  env: execEnv
});
console.log(result.toString());

process.exit(0);
