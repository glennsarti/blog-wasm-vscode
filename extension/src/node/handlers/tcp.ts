import * as net from 'net';
import * as vscode from 'vscode';
import { ServerOptions, StreamInfo } from 'vscode-languageclient/node';
import { IStatusBar } from '../../common/features/statusbar';
import { NodeConnectionHandler } from '../handler';
import { ILogger } from '../../common/logging';
import { ConnectionType } from '../../common/settings';

export class TcpConnectionHandler extends NodeConnectionHandler {
  constructor(
    context: vscode.ExtensionContext,
    statusBar: IStatusBar,
    protected logger: ILogger,
  ) {
    super(context, statusBar, logger);
    this.logger.debug(`Configuring ${ConnectionType[this.connectionType]}::TCP connection handler`);
    this.start();
  }

  get connectionType(): ConnectionType {
    return ConnectionType.Local;
  }

  createServerOptions(): ServerOptions {
    var host = "127.0.0.1";
    var port = 30337;

    this.logger.debug(
      `Starting language server client (host ${host} port ${port})`,
    );

    const serverOptions = () => {
      const socket = new net.Socket();

      socket.connect({
        port: port,
        host: host,
      });

      const result: StreamInfo = {
        writer: socket,
        reader: socket,
      };

      return Promise.resolve(result);
    };
    return serverOptions;
  }

  cleanup(): void {
    this.logger.debug(`No cleanup needed for ${this.protocolType}`);
  }
}