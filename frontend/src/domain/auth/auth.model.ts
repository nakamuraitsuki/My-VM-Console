import type { User } from "../user/user.model";

export type AuthStatus = 
  |'authenticated'
  | 'unauthenticated'
  | 'loading';

export interface AuthSession {
  readonly status: AuthStatus;
  readonly user: User | null;
}