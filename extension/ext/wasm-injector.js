const path = require("path");
const fs = require("fs");
const cp = require("child_process");

function wasmAsBase64() {
  let wasmServerPath = "C:\\Source\\blog-wasm-vscode\\lsp\\pkg\\demo-ls.wasm";
  console.log(`Reading ${wasmServerPath} as base64 ...`);
  return fs.readFileSync(wasmServerPath, "base64");
}

module.exports = function (source) {
  const magicCommand = "import wasm from 'go';";
  // Check if we need to inject WASM
  if (!source.startsWith(magicCommand)) {
    return source;
  }

  source = source.slice(magicCommand.length);

  console.log("Finding GOROOT");

  let result = cp.execSync("go env GOROOT");
  const goRoot = result.toString().trim();
  const wasmPath = path.join(goRoot, "misc", "wasm", "wasm_exec.js");
  console.log(`Reading ${wasmPath} ...`);
  const wasmSource = fs.readFileSync(wasmPath);

  wasmDecoder = `
const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';

// Use a lookup table to find the index.
const lookup = typeof Uint8Array === 'undefined' ? [] : new Uint8Array(256);
for (let i = 0; i < chars.length; i++) {
    lookup[chars.charCodeAt(i)] = i;
}

function wasmAsArrayBuffer() {
  let bufferLength = wasmBase64Content.length * 0.75,
    len = wasmBase64Content.length,
    i,
    p = 0,
    encoded1,
    encoded2,
    encoded3,
    encoded4;

  if (wasmBase64Content[wasmBase64Content.length - 1] === '=') {
    bufferLength--;
    if (wasmBase64Content[wasmBase64Content.length - 2] === '=') {
      bufferLength--;
    }
  }

  const arraybuffer = new ArrayBuffer(bufferLength),
      bytes = new Uint8Array(arraybuffer);

  for (i = 0; i < len; i += 4) {
    encoded1 = lookup[wasmBase64Content.charCodeAt(i)];
    encoded2 = lookup[wasmBase64Content.charCodeAt(i + 1)];
    encoded3 = lookup[wasmBase64Content.charCodeAt(i + 2)];
    encoded4 = lookup[wasmBase64Content.charCodeAt(i + 3)];

    bytes[p++] = (encoded1 << 2) | (encoded2 >> 4);
    bytes[p++] = ((encoded2 & 15) << 4) | (encoded3 >> 2);
    bytes[p++] = ((encoded3 & 3) << 6) | (encoded4 & 63);
  }

  return arraybuffer;
};
`;

  return (
    "// region: wasm_exec.js\n" +
    wasmSource +
    "\n// endregion: wasm_exec.js\n\n" +
    "// region: base64_wasm\n" +
    "wasmBase64Content = '" +
    wasmAsBase64() +
    "';\n" +
    wasmDecoder +
    "// endregion: base64_wasm\n\n" +
    source
  );
};
