import { useState } from 'react';
import { useCreateInstance } from '../hooks/useCreateInstance';
import { createImageID } from '../../../../domain/compute/compute.model';
import type { CreateInstanceRequest } from '@/domain/compute/compute.repository';
import styles from './CreateInstanceModal.module.css';

interface CreateInstanceModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: (instanceId: string) => void;
}

export const CreateInstanceModal = ({ isOpen, onClose, onSuccess }: CreateInstanceModalProps) => {
  const { isLoading, execute } = useCreateInstance();
  const [name, setName] = useState(`instance-${Date.now().toString().slice(-4)}`);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);

    const request: CreateInstanceRequest = {
      name: name,
      imageId: createImageID('img-ubuntu-2404'), // 今はイメージ固定
      vpcId: undefined,
      subnetId: undefined,
      cpu: 2,
      memory: 2048,
    };

    const result = await execute(request);
    if (result.type === 'success') {
      onSuccess?.(result.instance.id);
      onClose();
      return;
    }
    setError(result.error);
  };

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>インスタンスの新規作成</h2>
        
        <form onSubmit={handleSubmit}>
          <div className={styles.field}>
            <label htmlFor="inst-name">インスタンス名</label>
            <input
              id="inst-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              autoFocus
              className={styles.input}
            />
          </div>

          <div className={styles.configSummary}>
            <h3>構成の概要 (デフォルト)</h3>
            <ul>
              <li><strong>OS:</strong> Ubuntu 24.04 LTS</li>
              <li><strong>Spec:</strong> 2 vCPU / 2048 MiB</li>
              <li><strong>Network:</strong> Default VPC / Subnet</li>
            </ul>
            <p className={styles.note}>※ 詳細な構成変更は後ほど可能になる予定です。</p>
          </div>

          {error && <div className={styles.errorAlert}>{error}</div>}

          <div className={styles.actions}>
            <button type="button" onClick={onClose} className={styles.cancelBtn} disabled={isLoading}>
              キャンセル
            </button>
            <button type="submit" className={styles.submitBtn} disabled={isLoading}>
              {isLoading ? '準備中...' : '作成を開始する'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};