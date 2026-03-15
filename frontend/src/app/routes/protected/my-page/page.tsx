import { useNavigate } from "react-router";
import { useAuth } from "@/context/AuthContext";
import { useEffect } from "react";
import { MyVideoList } from "@/features/video/list-mine";
import { IconButton } from "@/ui/IconButton/IconButton";
import { RiAddLine } from "react-icons/ri";

const MY_PAGE_LIMIT = 10;

export const MyPage = () => {
  const navigate = useNavigate();
  const { session } = useAuth();

  useEffect(() => {
    if (session.status === "unauthenticated") {
      navigate("/");
    }
  }, [session.status, navigate]);

  if (session.status !== "authenticated") return null;

  return (
    <div style={{ padding: "20px", maxWidth: "800px", margin: "0 auto" }}>
      <h2>My Page</h2>
      <p>Welcome, {session.user?.name}!</p>

      <IconButton
        icon={<RiAddLine size={20} />}
        onClick={() => navigate("/upload")}
        label="Create New Video"
        ariaLabel="Upload New Video"
      />
      <hr />

      <h3>My Videos</h3>
      <MyVideoList
        limit={MY_PAGE_LIMIT}
        onCardSelect={(video) => navigate(`/video/${video.id}`, { state: { video } })}
        onSettingsClick={(video) => navigate(`/manage/${video.id}`, { state: { video } })}
      />
    </div>
  );
};
