const express = require('express')
const cors = require('cors')
const compression = require('compression')

const path = require('path');
const fs = require('fs');
const cp = require('child_process');

const app = express()
app.use(cors())
app.use(compression())

const PORT = 9001

let result = cp.execSync("go env GOROOT");
const goRoot = result.toString().trim();
const wasmPath = path.join(goRoot, "misc", "wasm", "wasm_exec.js");

app.get('/', function (req, res) {
  res.contentType("text/html");
  res.sendFile(__dirname + '/index.html')
})

app.get('/wasm_exec.js', function (req, res) {
  res.contentType("text/javascript");
  res.sendFile(wasmPath)
})

app.get('/wasm', function (req, res) {
  res.contentType("application/wasm");
  console.log("Serving wasm...");
  res.sendFile(path.resolve(__dirname + '/../pkg/demo-ls.wasm'))
  console.log("Served WASM");
})

app.listen(PORT, function(error) {
  if (error) {
    console.error(error);
  } else {
    console.info("❤️  Listening on port %s. Visit http://localhost:%s/ in your browser.", PORT, PORT);
  }
});
