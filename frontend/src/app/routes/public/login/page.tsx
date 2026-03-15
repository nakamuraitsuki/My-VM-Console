import { useAuth } from "@/context/AuthContext";
import { useEffect } from "react";
import { useNavigate } from "react-router";

export const LoginPage = () => {
  const { session, login } = useAuth();
  const navigate = useNavigate();
  // 現在はDummyなので、ログイン処理をしてすぐにリダイレクトする
  useEffect(() => {
    // すでに認証済み、または今まさにリクエスト中なら何もしない
    if (session.status === "authenticated" || session.status === "loading") {
      if (session.status === "authenticated") navigate("/");
      return;
    }

    // 完全に 'unauthenticated' の時だけ、一度だけ実行する
    login();
  }, [session.status]);

  return <div>Logging in...</div>;
};