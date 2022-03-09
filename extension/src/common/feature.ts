import * as vscode from 'vscode';

export interface IFeature extends vscode.Disposable {
  dispose(): any;
}
