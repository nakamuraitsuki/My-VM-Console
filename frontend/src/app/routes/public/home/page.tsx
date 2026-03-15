import { useNavigate } from "react-router";

export const HomePage = () => {
  const navigate = useNavigate();

  const handleVideoSelect = (videoId: string) => {
    // 選択した動画の詳細ページに遷移
    navigate(`/video/${videoId}`);
  };

  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", padding: "1rem" }}>
      <h2 style={{ textAlign: "center", marginBottom: "1rem" }}>Home</h2>

      <section>
        <h3 style={{ marginBottom: "1rem" }}>Public Videos</h3>
      </section>
    </div>
  );
};
