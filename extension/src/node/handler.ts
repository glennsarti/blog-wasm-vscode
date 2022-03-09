import * as vscode from 'vscode';
import { ServerOptions, LanguageClient } from 'vscode-languageclient/node';
import { CommonLanguageClient, LanguageClientOptions } from 'vscode-languageclient';
import { IStatusBar } from '../common/features/statusbar';
import { ILogger } from '../common/logging';
import { ConnectionHandler } from '../common/handler';

export abstract class NodeConnectionHandler extends ConnectionHandler {
  protected constructor(
    protected context: vscode.ExtensionContext,
    protected statusBar: IStatusBar,
    protected logger: ILogger,
  ) {
    super(context, statusBar, logger);

    this.logger.debug('Creating server options');
  }

  async createLanguageClient(name: string, clientOptions: LanguageClientOptions): Promise<CommonLanguageClient> {
    this.logger.debug('Configuring language server options');
    const serverOptions = this.createServerOptions();

    const langClient: CommonLanguageClient = new LanguageClient(name, serverOptions, clientOptions);

    return Promise.resolve(langClient);
  }

  abstract createServerOptions(): ServerOptions;
}