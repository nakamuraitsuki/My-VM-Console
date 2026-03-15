import type { Result } from "../core/result";
import type { User } from "../user/user.model";
import type { AuthSession } from "./auth.model";

export interface IAuthRepository {
  loginOIDC(): Promise<Result<User, AuthError>>; // OIDC 用
  fetchCurrentSession(): Promise<AuthSession>;
}

export type AuthError =
  | "INVALID_CREDENTIALS"
  | "NETWORK_ERROR"
  | "SERVER_ERROR"
  | "UNKNOWN_ERROR";