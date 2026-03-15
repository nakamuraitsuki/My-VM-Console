import { useNavigate } from "react-router";
import { VideoUploadForm } from "@/features/video/create";
import styles from "./page.module.css";

export const VideoUploadPage = () => {
  const navigate = useNavigate();

  return (
    <div className={styles.container}>
      <VideoUploadForm
        onBack={() => navigate(-1)}
        onSuccess={() => navigate("/my-page")}
      />
    </div>
  );
};
