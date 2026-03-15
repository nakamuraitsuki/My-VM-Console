import axios from "axios";
import { apiClient } from "../../api/client";
import type { AuthSession } from "../../domain/auth/auth.model";
import type { AuthError, IAuthRepository } from "../../domain/auth/auth.repository";
import type { User } from "../../domain/user/user.model";
import { type Result } from "../../domain/core/result";

export class AuthRepositoryImpl implements IAuthRepository {
  async loginOIDC(): Promise<Result<User, AuthError>> {
    window.location.href = "/api/users/login";

    return new Promise(() => {() => {}}); // OIDCのリダイレクトでページ遷移するため、Promiseは完了しない
  }

  async fetchCurrentSession(): Promise<AuthSession> {
    try {
      const { data } = await apiClient.get<User>("/api/users/me");
      return {
        status: "authenticated",
        user: data,
      };
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 401) {
          return { status: "unauthenticated", user: null };
        }
      }

      console.error("Unexpected error during fetching session:", error);
      throw error;
    }
  }
}
