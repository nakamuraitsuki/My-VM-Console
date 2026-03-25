import type { Instance } from "@/domain/compute/compute.model";
import styles from "./ListMine.module.css";

type ListMineItemProps = {
  instance: Instance;
  onClick: (instanceId: string) => void;
};

const statusText: Record<Instance["status"], string> = {
  PENDING: "起動準備中",
  RUNNING: "稼働中",
  STOPPED: "停止中",
  ERROR: "異常",
};

const statusClass: Record<Instance["status"], string> = {
  PENDING: styles.statusPending,
  RUNNING: styles.statusRunning,
  STOPPED: styles.statusStopped,
  ERROR: styles.statusError,
};

export const ListMineItem = ({ instance, onClick }: ListMineItemProps) => {
  return (
    <button
      type="button"
      className={styles.item}
      onClick={() => onClick(instance.id)}
      aria-label={`${instance.name} の詳細を開く`}
    >
      <div className={styles.itemMain}>
        <p className={styles.itemName}>{instance.name}</p>
        <p className={styles.itemMeta}>ID: {instance.id}</p>
      </div>

      <div className={styles.itemSide}>
        <span className={`${styles.status} ${statusClass[instance.status]}`}>
          {statusText[instance.status]}
        </span>
        <p className={styles.itemMeta}>Private IP: {instance.privateIp}</p>
      </div>
    </button>
  );
};
