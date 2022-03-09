/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Microsoft Corporation. All rights reserved.
 *  Licensed under the MIT License. See License.txt in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

//@ts-check
"use strict";

//@ts-check
/** @typedef {import('webpack').Configuration} WebpackConfig **/

const path = require("path");
//const webpack = require('webpack');

/** @type WebpackConfig */
const nodeConfig = {
  target: "node",
  mode: "none",
  entry: {
    extension: "./src/node/extension.ts",
  },
  output: {
    filename: "[name].js",
    path: path.join(__dirname, "out", "node"),
    libraryTarget: "commonjs",
    devtoolModuleFilenameTemplate: "../../[resource-path]",
  },
  externals: {
    vscode: "commonjs vscode", // ignored because it doesn't exist
  },
  devtool: "nosources-source-map",
  resolve: {
    extensions: [".ts", ".js"],
    alias: {
      debug: path.join(__dirname, "polyfill", "debug.js"),
    },
  },
  module: {
    rules: [
      {
        test: /\.ts$/,
        exclude: /node_modules/,
        use: [
          {
            loader: "ts-loader",
          },
        ],
      },
    ],
  },
  infrastructureLogging: {
    level: "log", // enables logging required for problem matchers
  },
};

module.exports = [nodeConfig];
