import { ILogger } from '../logging';

export class NullLogger implements ILogger {
  show(): void { };
  verbose(_: string): void { };
  debug(_: string): void { };
  normal(_: string): void { };
  warning(_: string): void { };
  error(_: string): void { };
}
