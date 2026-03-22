import styles from "./Dashboard.module.css";
import { CreateInstanceButton } from "@/features/compute";

/**
 * ダッシュボードページ
 * 複数 feature の集約点
 */
export const DashboardPage = () => {
  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <h1 className={styles.title}>Dashboard</h1>

        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>インスタンス管理</h2>
          <CreateInstanceButton
            onSuccess={(instanceId: string) => {
              console.log('インスタンス作成成功:', instanceId);
            }}
            onError={(error: string) => {
              console.error('インスタンス作成失敗:', error);
            }}
          />
        </section>
      </div>
    </div>
  );
};
