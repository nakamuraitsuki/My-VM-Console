import type { AuthSession } from "../../domain/auth/auth.model";
import type { AuthError, IAuthRepository } from "../../domain/auth/auth.repository";
import { failure, success, type Result } from "../../domain/core/result";
import type { User } from "../../domain/user/user.model";

export class AuthRepositoryMock implements IAuthRepository {
  async login(_email: string, _password: string): Promise<Result<User, AuthError>> {
    try {
      const user: User = {
        id: "1",
        name: "Locak Mock User",
        bio: "This is a mock user for testing purposes.",
        role: "user",
      };
      // NOTE: 簡易的にLocalStorageにデータ保存
      localStorage.setItem("user_id", user.id);

      return success(user);
    } catch (error) {
      return failure("UNKNOWN_ERROR");
    }
  }

  async logout(): Promise<void> {
    localStorage.removeItem("user_id");
  }

  async loginOIDC(): Promise<Result<User, AuthError>> {
    // OIDCのリダイレクトは実際には行わないが、成功したと仮定してユーザーデータを返す
    try {
      const user: User = {
        id: "1",
        name: "Locak Mock User",
        bio: "This is a mock user for testing purposes.",
        role: "user",
      };
      localStorage.setItem("user_id", user.id);
      return success(user);
    } catch (error) {
      return failure("UNKNOWN_ERROR");
    }
  }

  async fetchCurrentSession(): Promise<AuthSession> {
    const userId = localStorage.getItem("user_id");
    if (userId) {
      return {
        status: "authenticated",
        user: {
          id: userId,
          name: "Locak Mock User",
          bio: "This is a mock user for testing purposes.",
          role: "user",
        },
      };
    } else {
      return { status: "unauthenticated", user: null };
    }
  }
}