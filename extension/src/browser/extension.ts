// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import { IFeature } from '../common/feature';
import { OutputChannelLogger } from '../common/logging/outputchannel';
import { StatusBarFeature } from '../common/features/statusbar';
import { WebWorkerConnectionHandler } from './handlers/webworker';
import { ConnectionHandler } from '../common/handler';

import { DemoLangID } from '../common/contants';

let logger: OutputChannelLogger;
let connectionHandler: ConnectionHandler;
let extensionFeatures: IFeature[] = [];

export async function activate(context: vscode.ExtensionContext) {
  // Use the console to output diagnostic information (console.log) and errors (console.error)
  // This line of code will only be executed once when your extension is activated

  let logger = new OutputChannelLogger('debug');

  const statusBar = new StatusBarFeature([DemoLangID], logger, context);
  connectionHandler = new WebWorkerConnectionHandler(context, statusBar, logger);
  connectionHandler.init();

  console.log('Congratulations, your extension "vscode-demo" is now active!');
}

// this method is called when your extension is deactivated
export function deactivate() {
  // Dispose all extension features
  extensionFeatures.forEach((feature) => {
    feature.dispose();
  });

  if (connectionHandler !== undefined) {
    connectionHandler.stop();
  }
}
