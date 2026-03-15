import axios from "axios";
import { apiClient } from "../../api/client";
import { failure, success, type Result } from "../../domain/core/result";
import type { IUserRepository, UpdateIconError } from "../../domain/user/user.repository";

export class UserRepositoryImpl implements IUserRepository {
  async updateProfile(name: string | null, bio: string | null): Promise<void> {
    await apiClient.patch("/api/users/me/profile", {
      name,
      bio,
    });
  }

  async updateIcon(icon: File | null): Promise<Result<void, UpdateIconError>> {
    if (!icon) return failure<UpdateIconError>("INVALID_FORMAT");

    const formData = new FormData();
    formData.append("file", icon);

    try {
      await apiClient.put("/api/users/me/icon", formData);
      return success(undefined);
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 413) return failure("FILE_TOO_LARGE");
        if (error.response?.status === 415) return failure("INVALID_FORMAT");
      }
      return failure("NETWORK_ERROR");
    }
  }
}