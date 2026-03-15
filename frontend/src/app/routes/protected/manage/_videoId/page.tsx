import { useNavigate, useLocation, useParams } from "react-router";
import { useAuth } from "@/context/AuthContext";
import { useEffect } from "react";
import { VideoEditPanel } from "@/features/video/edit";
import type { Video } from "@/domain/video/video.model";
import { IconButton } from "@/ui/IconButton/IconButton";
import { BiArrowBack } from "react-icons/bi";

export const VideoManagePage = () => {
  const navigate = useNavigate();
  const { session } = useAuth();
  const { videoId } = useParams<{ videoId: string }>();
  const location = useLocation();
  const video = location.state?.video as Video | undefined;

  useEffect(() => {
    if (session.status === "unauthenticated") {
      navigate("/");
    }
  }, [session.status, navigate]);

  if (session.status !== "authenticated") return null;

  if (!video || video.id !== videoId) {
    return (
      <div style={{ padding: "20px" }}>
        <p>動画情報が見つかりませんでした。</p>
        <IconButton
          onClick={() => navigate("/my-page")}
          icon={<BiArrowBack size={20} />}
          label="マイページに戻る"
        />
      </div>
    );
  }

  return (
    <div style={{ padding: "20px", maxWidth: "800px", margin: "0 auto" }}>
      <IconButton
        onClick={() => navigate("/my-page")}
        icon={<BiArrowBack size={20} />}
        label="マイページに戻る"
      />
      <h2 style={{ marginBottom: "24px" }}>{video.title}</h2>
      <VideoEditPanel video={video} />
    </div>
  );
};
