/* eslint-disable @typescript-eslint/no-use-before-define */
'use strict';

import * as vscode from 'vscode';
import { IFeature } from '../feature';
import { ConnectionStatus } from '../interfaces';
import { ILogger } from '../logging';

class StatusBarProvider {
  private statusBarItem: vscode.StatusBarItem;
  private logger: ILogger;

  constructor(langIDs: string[], logger: ILogger) {
    this.logger = logger;
    this.statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 1);
    this.statusBarItem.show();

    vscode.window.onDidChangeActiveTextEditor((textEditor) => {
      if (textEditor === undefined || langIDs.indexOf(textEditor.document.languageId) === -1) {
        this.statusBarItem.hide();
      } else {
        this.statusBarItem.show();
      }
    });
  }

  public setConnectionStatus(statusText: string, status: ConnectionStatus, toolTip: string): void {
    this.logger.debug(`Setting status bar to ${statusText}`);
    // Icons are from https://octicons.github.com/
    let statusIconText: string;
    let statusColor: string;

    switch (status) {
      case ConnectionStatus.RunningLoaded:
        statusIconText = '$(terminal) ';
        statusColor = '#affc74';
        break;
      case ConnectionStatus.RunningLoading:
        statusIconText = '$(sync~spin) ';
        statusColor = '#affc74';
        break;
      case ConnectionStatus.Failed:
        statusIconText = '$(alert) ';
        statusColor = '#fcc174';
        break;
      default:
        // ConnectionStatus.NotStarted
        // ConnectionStatus.Starting
        // ConnectionStatus.Stopping
        statusIconText = '$(gear) ';
        statusColor = '#f3fc74';
        break;
    }

    statusIconText = (statusIconText + statusText).trim();
    this.statusBarItem.color = statusColor;
    // Using a conditional here because resetting a $(sync~spin) will cause the animation to restart. Instead
    // Only change the status bar text if it has actually changed.
    if (this.statusBarItem.text !== statusIconText) {
      this.statusBarItem.text = statusIconText;
    }
    this.statusBarItem.tooltip = toolTip;
  }

  public showConnectionMenu() {
    const menuItems: StatusBarMenuItem[] = [];

    vscode.window.showQuickPick<StatusBarMenuItem>(menuItems).then((selectedItem) => {
      if (selectedItem) {
        selectedItem.callback();
      }
    });
  }
}

class StatusBarMenuItem implements vscode.QuickPickItem {
  public description = '';

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  constructor(public readonly label: string, public readonly callback: () => void = () => { }) { }

  // TODO Add a demo item
}

export interface IStatusBar {
  setConnectionStatus(statusText: string, status: ConnectionStatus, toolTip: string): void;
}

export class StatusBarFeature implements IFeature, IStatusBar {
  private provider: StatusBarProvider;

  constructor(langIDs: string[], logger: ILogger, context: vscode.ExtensionContext) {
    this.provider = new StatusBarProvider(langIDs, logger);
  }

  public setConnectionStatus(statusText: string, status: ConnectionStatus, toolTip: string): void {
    this.provider.setConnectionStatus(statusText, status, toolTip);
  }

  public dispose(): any {
    return undefined;
  }
}
