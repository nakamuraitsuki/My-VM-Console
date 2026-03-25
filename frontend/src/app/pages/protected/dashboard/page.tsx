import { useState } from 'react';
import styles from "./Dashboard.module.css";
import { CreateInstanceModal } from "@/features/compute/create";
import { IconButton } from "@/ui/IconButton/IconButton";

/**
 * ダッシュボードページ
 */
export const DashboardPage = () => {
  // 2. モーダルの開閉状態を定義
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <h1 className={styles.title}>Dashboard</h1>

        <section className={styles.section}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>インスタンス管理</h2>
            <IconButton
              icon="+"
              label="新規作成"
              onClick={() => setIsModalOpen(true)}
            />
          </div>
          <CreateInstanceModal
            isOpen={isModalOpen}
            onClose={() => setIsModalOpen(false)}
            onSuccess={(instanceId: string) => {
              console.log('作成開始:', instanceId);
            }}
          />
        </section>
      </div>
    </div>
  );
};