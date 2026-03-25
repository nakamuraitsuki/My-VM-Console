import { useParams } from "react-router";
import styles from "./InstanceDetail.module.css";

export const InstanceDetailPage = () => {
  const { instanceId } = useParams();

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>インスタンス詳細</h1>
      <p className={styles.lead}>対象インスタンス: {instanceId}</p>

      <section className={styles.placeholderSection}>
        <h2 className={styles.sectionTitle}>今後表示予定</h2>
        <ul className={styles.list}>
          <li>SSH鍵 CA 情報</li>
          <li>SSH公開鍵の追加状況</li>
          <li>インスタンスのドメイン情報</li>
        </ul>
      </section>
    </div>
  );
};
