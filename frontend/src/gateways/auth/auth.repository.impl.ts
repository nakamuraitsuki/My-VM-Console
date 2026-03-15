import axios from "axios";
import { apiClient } from "../../api/client";
import type { AuthSession } from "../../domain/auth/auth.model";
import type { AuthError, IAuthRepository } from "../../domain/auth/auth.repository";
import type { User } from "../../domain/user/user.model";
import { failure, success, type Result } from "../../domain/core/result";

export class AuthRepositoryImpl implements IAuthRepository {
  async login(_email: string, _password: string): Promise<Result<User, AuthError>> {
    // NOTE: 現在はBackendにDummyLoginがあるので引数は使わない
    try {
      const { data } = await apiClient.post<User>("/api/users/login");
      return success(data);
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        return failure("INVALID_CREDENTIALS");
      }
      throw error;
    }
  }

  async loginOIDC(): Promise<Result<User, AuthError>> {
    window.location.href = "/api/users/login";

    return new Promise(() => {() => {}}); // OIDCのリダイレクトでページ遷移するため、Promiseは完了しない
  }

  async logout(): Promise<void> {
    await apiClient.post("/api/users/logout");
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
