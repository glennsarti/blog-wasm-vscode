import * as vscode from 'vscode';
import { LanguageClient } from 'vscode-languageclient/browser';
import { CommonLanguageClient, LanguageClientOptions } from 'vscode-languageclient';
import { IStatusBar } from '../common/features/statusbar';
import { ILogger } from '../common/logging';
import { ConnectionHandler } from '../common/handler';
import { DemoLangID } from '../common/contants';

export abstract class BrowserConnectionHandler extends ConnectionHandler {
  protected constructor(
    protected context: vscode.ExtensionContext,
    protected statusBar: IStatusBar,
    protected logger: ILogger
  ) {
    super(context, statusBar, logger);
  }

  async createLanguageClient(name: string, clientOptions: LanguageClientOptions): Promise<CommonLanguageClient>{
    this.logger.debug('Configuring web worker');
    const worker = await this.createWorker();
    return Promise.resolve(new LanguageClient(DemoLangID, name, clientOptions, worker));
  }

  abstract createWorker(): Promise<Worker>;
}
