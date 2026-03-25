import { useState } from 'react';
import { IconButton } from '../../../../ui/IconButton/IconButton';
import { useCreateInstance } from '../hooks/useCreateInstance';
import { createImageID } from '../../../../domain/compute/compute.model';
import styles from './CreateInstanceButton.module.css';
import type { CreateInstanceRequest } from '@/domain/compute/compute.repository';

interface CreateInstanceButtonProps {
  onSuccess?: (instanceId: string) => void;
  onError?: (error: string) => void;
}

/**
 * インスタンス作成ボタンコンポーネント
 * Pure なUIコンポーネント
 */
export const CreateInstanceButton = ({
  onSuccess,
  onError,
}: CreateInstanceButtonProps) => {
  const { isLoading, error, execute } = useCreateInstance();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleClick = async () => {
    setSuccessMessage(null);

    // TODO: 後々詳細設定ができるようにする。
    // 現在はサンプル値で実装
    const request: CreateInstanceRequest = {
      name: `instance-${Date.now()}`,
      imageId: createImageID('img-ubuntu-2404'),
      vpcId: undefined,
      subnetId: undefined,
      cpu: 2,
      memory: 2048,
    };

    const result = await execute(request);
    if (result.type === 'success') {
      const message = `インスタンス作成に成功しました。ID: ${result.instance.id}`;
      setSuccessMessage(message);
      onSuccess?.(result.instance.id);
      return;
    }

    {
      const errorMessage = result.error;
      onError?.(errorMessage);
    }
  };

  return (
    <div className={styles.container}>
      <IconButton
        icon="+"
        label="インスタンスを作成"
        onClick={handleClick}
        disabled={isLoading}
        ariaLabel="新しいインスタンスを作成"
      />

      {error && (
        <div className={styles.errorMessage}>
          <p>エラー: {error}</p>
        </div>
      )}

      {successMessage && (
        <div className={styles.successMessage}>
          <p>{successMessage}</p>
        </div>
      )}
    </div>
  );
};
