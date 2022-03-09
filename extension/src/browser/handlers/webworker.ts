import * as vscode from 'vscode';
import { IStatusBar } from '../../common/features/statusbar';
import { BrowserConnectionHandler } from '../handler';
import { ILogger } from '../../common/logging';
import { ConnectionType } from '../../common/settings';

export class WebWorkerConnectionHandler extends BrowserConnectionHandler {
  workers: Worker[] = [];

  constructor(
    context: vscode.ExtensionContext,
    statusBar: IStatusBar,
    logger: ILogger,
  ) {
    super(context, statusBar, logger);
    this.logger.debug(`Configuring ${ConnectionType[this.connectionType]}::Web Worker connection handler`);
    this.start();
  }

  get connectionType(): ConnectionType {
    return ConnectionType.Local;
  }

  timeout(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async workerReady(worker: Worker): Promise<Boolean> {
    var isReady: boolean = false;

    const listener = (v: Event) => {
      // Do we really care about the message?  Nope.
      isReady = true;
    };
    worker.addEventListener("message",listener);
    const waitTill = new Date(new Date().getTime() + 4 * 1000);
    while (!isReady || waitTill > new Date()) {
      this.logger.debug("Waiting for the Worker to be ready...")
      await this.timeout(500);
    }
    worker.removeEventListener("message",listener);

    if (isReady) {
      this.logger.debug("Worker is ready")
    } else {
      this.logger.debug("Worker is not ready")
    }
    return Promise.resolve(isReady);
  }

  async createWorker(): Promise<Worker> {
    // Create a worker. The worker main file implements the language server.
    const serverMain = vscode.Uri.joinPath(this.context.extensionUri, 'out/browser/wasm-lsp.js');
    const worker = new Worker(serverMain.toString());
    if (this.workers === undefined) {
      this.workers = [];
    }
    this.workers.push(worker);

    // TODO: What do you do if it's not ready?
    const ready = await this.workerReady(worker);

    return Promise.resolve(worker);
  }

  cleanup(): void {
    this.logger.debug(`Stopping Web Workers`);
    this.workers.forEach(worker => {
      worker.terminate();
    });
  }
}