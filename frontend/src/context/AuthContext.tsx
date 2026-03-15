import React, { createContext, useEffect, useState } from "react";
import type { AuthSession } from "../domain/auth/auth.model";
import type { AuthError } from "../domain/auth/auth.repository";
import { failure, success, type Result } from "../domain/core/result";
import { useServices } from "./ServiceContext";

// NOTE: AuthContextType には、Sessionと更新メソッドをカプセル化
interface AuthContextType {
  session: AuthSession;
  login: () => Promise<Result<AuthSession, AuthError>>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const { authRepo } = useServices();

  // 初期状態 `loading`
  const [session, setSession] = useState<AuthSession>({
    status: 'loading',
    user: null,
  });

  useEffect(() => {
    const initAuth = async () => {
      try {
        const currentSession = await authRepo.fetchCurrentSession();
        setSession(currentSession);
      } catch (_error) {
        setSession({ status: 'unauthenticated', user: null });
      }
    }
    initAuth();
  }, [authRepo]);

  const login = async (): Promise<Result<AuthSession, AuthError>> => {
    setSession({ status: 'loading', user: null });

    const result = await authRepo.loginOIDC();

    if (result.success) {
      const sessionData: AuthSession = result.data
        ? { status: 'authenticated', user: result.data }
        : { status: 'unauthenticated', user: null };

      setSession(sessionData);
      return success(sessionData);
    }

    // 失敗時
    const errorSession: AuthSession = { status: 'unauthenticated', user: null };
    setSession(errorSession);
    return failure(result.error);
  }

  const logout = async () => {
    await authRepo.logout();
    setSession({ status: 'unauthenticated', user: null });
  }

  return (
    <AuthContext.Provider value={{ session, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = (): AuthContextType => {
  const context = React.useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}