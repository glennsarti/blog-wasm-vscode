// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import { IFeature } from '../common/feature';
import { NullLogger } from '../common/logging/null';
import { StatusBarFeature } from '../common/features/statusbar';
import { TcpConnectionHandler } from './handlers/tcp';
import { ConnectionHandler } from '../common/handler';
import { DemoLangID } from '../common/contants';

let extensionFeatures: IFeature[] = [];
let connectionHandler: ConnectionHandler;

export function activate(context: vscode.ExtensionContext) {
  let logger = new NullLogger();

  const statusBar = new StatusBarFeature([DemoLangID], logger, context);
  connectionHandler = new TcpConnectionHandler(context, statusBar, logger);
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
