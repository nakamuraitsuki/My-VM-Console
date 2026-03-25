import { useListMineQuery } from "../hooks/useListMineQuery";
import { ListMineItem } from "./ListMineItem";
import styles from "./ListMine.module.css";

type ListMineProps = {
  onSelect?: (instanceId: string) => void;
};

export const ListMine = ({ onSelect }: ListMineProps) => {
  const { data, isPending, isError, error } = useListMineQuery();

  if (isPending) {
    return (
      <div className={styles.stateBox}>
        <div className={styles.inlineSpinner} aria-hidden="true" />
        <p className={styles.stateText}>インスタンスを読み込み中です...</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className={styles.stateBox}>
        <p className={styles.errorText}>インスタンス一覧の取得に失敗しました。</p>
        <p className={styles.stateSubText}>{error.message}</p>
      </div>
    );
  }

  if (!data.insts || data.insts.length === 0) {
    return (
      <div className={styles.stateBox}>
        <p className={styles.stateText}>インスタンスがまだありません。</p>
        <p className={styles.stateSubText}>右上の「新規作成」から作成を開始できます。</p>
      </div>
    );
  }

  return (
    <div className={styles.list}>
      {data.insts.map((instance) => (
        <ListMineItem
          key={instance.id}
          instance={instance}
          onClick={(instanceId) => onSelect?.(instanceId)}
        />
      ))}
    </div>
  );
};
