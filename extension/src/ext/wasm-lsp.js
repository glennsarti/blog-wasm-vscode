import wasm from 'go';

var golangReady = false;

onmessage = async function (e) {
  // DEBUG: console.log("Worker: Message received from main");
  // DEBUG: console.log(e.data);
  // The data is a JSON object, but for now we can convert it to a JSON string
  let msg = JSON.stringify(e.data);
  global.lsponmessage(msg);
};

// This needs to be in the global scope for GoLang
global.workerReady = () => {
  // TODO: Would be nice to have a better message but here we are!
  postMessage("Worker is ready!");
  console.log("Web Worker is ready");
  golangReady = true;
};

// This needs to be in the global scope for GoLang
global.goLangPostMessage = (message) => {
  // DEBUG: console.log("Worker: Message received from golang");
  // The message is a JSON string. Need to convert it back to a JSON object
  // TODO: What if the parsing fails?
  const jsonMsg = JSON.parse(message);
  postMessage(jsonMsg);
};

async function initWasm() {
  console.log("Creating Go instance...");
  go = new Go();
  console.log("Fetching WASM...");
  let result = await WebAssembly.instantiate(
    // wasmAsArrayBuffer() is injected as part of the webpacking
    wasmAsArrayBuffer(),
    go.importObject
  );
  console.log("Running WASM...");
  go.run(result.instance);
}
initWasm();

console.debug("wasm-lsp completed");

module.exports = {};
exports = {};
