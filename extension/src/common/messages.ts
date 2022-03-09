import { RequestType0 } from 'vscode-languageclient';

export interface LSPVersionDetails {
  lspVersion: string;
}

export namespace LSPVersionRequest {
  export const type = new RequestType0<LSPVersionDetails, void>('demo/getVersion');
}
