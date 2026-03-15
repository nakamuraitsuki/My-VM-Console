import type { Result } from "../core/result";

export interface IUserRepository {
  updateProfile: (name: string | null, bio: string | null) => Promise<void>;
  updateIcon: (icon: File | null) => Promise<Result<void, UpdateIconError>>;
}

export type UpdateIconError =
  | "INVALID_FORMAT"
  | "FILE_TOO_LARGE"
  | "NETWORK_ERROR"
  | "UNKNOWN_ERROR";