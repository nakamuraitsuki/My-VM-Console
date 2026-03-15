import { useNavigate, useParams } from "react-router";
import { VideoPlayerContainer } from "@/features/video/play";

export const VideoPlayPage = () => {
  const { videoId } = useParams<{ videoId: string }>();
  const navigate = useNavigate();

  // IDがない場合は不正なアクセスとして処理
  if (!videoId) {
    navigate("/my", { replace: true });
    return null;
  }

  return (
    <div style={{ container: "var(--container-width)", margin: "0 auto" }}>
      <VideoPlayerContainer
        videoId={videoId}
        onErrorBack={() => navigate("/my")}
      />
    </div>
  );
};