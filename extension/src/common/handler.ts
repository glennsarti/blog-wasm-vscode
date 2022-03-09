import * as vscode from 'vscode';
import { CommonLanguageClient, LanguageClientOptions, RevealOutputChannelOn } from 'vscode-languageclient';
import { IStatusBar } from './features/statusbar';
import { ConnectionStatus } from './interfaces';
import { ConnectionType, ProtocolType } from './settings';
import { DemoLangID } from './contants';
import { ILogger } from './logging';
import { LSPVersionDetails, LSPVersionRequest } from './messages';

export abstract class ConnectionHandler {
  private timeSpent: number;

  private _status: ConnectionStatus = ConnectionStatus.NotStarted;
  public get status(): ConnectionStatus {
    return this._status;
  }

  private _languageClient: CommonLanguageClient | undefined = undefined;
  public get languageClient(): CommonLanguageClient | undefined {
    return this._languageClient;
  }

  abstract get connectionType(): ConnectionType;

  public get protocolType(): ProtocolType {
    return ProtocolType.TCP;
  }

  protected constructor(
    protected context: vscode.ExtensionContext,
    protected statusBar: IStatusBar,
    protected logger: ILogger,
  ) {
    this.timeSpent = Date.now();
  }

  createClientOptions(langid: string): LanguageClientOptions {
    const documents = [
      { scheme: 'file', language: langid },
    ];

    return {
      documentSelector: documents,
      //outputChannel: this.logger.logChannel,
      revealOutputChannelOn: RevealOutputChannelOn.Info,
    };
  }

  abstract createLanguageClient(name: string, clientOptions: LanguageClientOptions): Promise<CommonLanguageClient>
  abstract cleanup(): void;

  async init(): Promise<void> {
    this.timeSpent = Date.now();
    this.setConnectionStatus('Initializing', ConnectionStatus.Initializing);

    this.logger.debug('Configuring language client options');
    const clientOptions: LanguageClientOptions = this.createClientOptions(DemoLangID);

    this.logger.debug('Creating language client');
    this._languageClient = await this.createLanguageClient('DemoVSCode', clientOptions);
    this._languageClient
      .onReady()
      .then(
        () => {
          this.setConnectionStatus('Loading Demo', ConnectionStatus.Starting);
          this.queryLanguageServerStatusWithProgress();
        },
        (reason) => {
          this.setConnectionStatus('Starting error', ConnectionStatus.Starting);
          if (this.languageClient !== undefined) { this.languageClient.error(reason) };
        },
      )
      .catch((reason) => {
        this.setConnectionStatus('Failure', ConnectionStatus.Failed);
      });
    this.setConnectionStatus('Initialization Complete', ConnectionStatus.InitializationComplete);

    this.start();
  }

  start(): void {
    this.logger.debug('Starting the language client...');
    if (this.languageClient === undefined) {
      this.logger.error("Language Client is undefined")
      return;
    }
    this.setConnectionStatus('Starting languageserver', ConnectionStatus.Starting, '');
    this.context.subscriptions.push(this.languageClient.start());
  }

  stop(): void {
    this.setConnectionStatus('Stopping languageserver', ConnectionStatus.Stopping, '');
    if (this.languageClient !== undefined) {
      this.timeSpent = Date.now() - this.timeSpent;
      this.languageClient.stop();
    }

    this.logger.debug('Running cleanup');
    this.cleanup();
    this.setConnectionStatus('Stopped languageserver', ConnectionStatus.Stopped, '');
  }

  public setConnectionStatus(message: string, status: ConnectionStatus, toolTip?: string) {
    var tipText = ""
    if (toolTip !== undefined) { tipText = toolTip; }
    this._status = status;
    this.statusBar.setConnectionStatus(message, status, tipText);
  }

  private queryLanguageServerStatusWithProgress() {
    return new Promise((resolve) => {
      let count = 0;
      let lastVersionResponse: LSPVersionDetails;
      const handle = setInterval(() => {
        count++;

        // After 30 seonds timeout the progress
        if (count >= 30 || this._languageClient === undefined) {
          clearInterval(handle);
          this.setConnectionStatus(lastVersionResponse.lspVersion, ConnectionStatus.RunningLoaded, '');
          resolve;
          return;
        }

        this._languageClient.sendRequest(LSPVersionRequest.type).then((details) => {
          lastVersionResponse = details;
          clearInterval(handle);
          this.setConnectionStatus(details.lspVersion, ConnectionStatus.RunningLoaded, '');
          resolve;
        });
      }, 1000);
    });
  }
}