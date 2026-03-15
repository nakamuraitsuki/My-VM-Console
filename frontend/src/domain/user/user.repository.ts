import type { Result } from "../core/result";
import type { User } from "./user.model";

export interface IUserRepository {
  getCurrentUser(): Promise<Result<User>>;
}