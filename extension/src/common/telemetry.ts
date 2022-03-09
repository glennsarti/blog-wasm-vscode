/* eslint-disable @typescript-eslint/no-use-before-define */
import TelemetryReporter from '@vscode/extension-telemetry';

export const reporter: TelemetryReporter = getTelemetryReporter();

function getTelemetryReporter() {
  const pkg = getPackageInfo();
  const reporter: TelemetryReporter = new TelemetryReporter(pkg.name, pkg.version, pkg.aiKey);
  return reporter;
}

function getPackageInfo(): IPackageInfo {
  return {
    name: 'bogus',
    version:'bogus',
    aiKey: 'bogus'
  };
}

interface IPackageInfo {
  name: string;
  version: string;
  aiKey: string;
}