import { useState } from 'react';
import styles from "./Dashboard.module.css";
import { CreateInstanceModal } from "@/features/compute/create";
import { ListMine } from "@/features/compute/dashboard";
import { IconButton } from "@/ui/IconButton/IconButton";
import { IoMdAdd } from 'react-icons/io';
import { useNavigate } from 'react-router';
import type { AppRoute } from '@/app/router';

/**
 * ダッシュボードページ
 */
export const DashboardPage = () => {
  const navigate = useNavigate();
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <h1 className={styles.title}>Dashboard</h1>

        <section className={styles.section}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>インスタンス管理</h2>
            <IconButton
              icon={<IoMdAdd size={24} />}
              label="新規作成"
              onClick={() => setIsModalOpen(true)}
            />
          </div>

          <ListMine
            onSelect={(instanceId) => {
              navigate(`/dashboard/instances/${instanceId}`);
            }}
          />

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

export const route: AppRoute = {
  path: "/dashboard",
  element: <DashboardPage />,
  meta: { name: "dashboard", requiresAuth: true },
};