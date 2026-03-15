export type UserId = string;

export interface User {
  readonly id: UserId;
  readonly name: string;
  readonly bio: string;
  readonly iconKey?: string;
}