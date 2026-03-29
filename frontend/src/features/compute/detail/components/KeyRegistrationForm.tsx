'use client';

import { useState } from 'react';
import type { RawPublicKey, SignedCertificate } from '@/domain/compute/key.model';
import { createRawPublicKey } from '@/domain/compute/key.model';
import type { KeyError } from '@/domain/compute/key.repository';
import type { UseSignCertificateState } from '../hooks/useSignCertificate';
import type { UseAuthorizePublicKeyState } from '../hooks/useAuthorizePublicKey';
import styles from './KeyRegistrationForm.module.css';

interface KeyRegistrationFormProps {
  instanceId: string;
  signCertificateState: UseSignCertificateState & {
    signCertificate: (publicKey: RawPublicKey) => void;
    reset: () => void;
  };
  authorizePublicKeyState: UseAuthorizePublicKeyState & {
    authorizePublicKey: (params: {
      instanceId: string;
      publicKey: RawPublicKey;
    }) => void;
    reset: () => void;
  };
}

export const KeyRegistrationForm = ({
  instanceId,
  signCertificateState,
  authorizePublicKeyState,
}: KeyRegistrationFormProps) => {
  const [publicKeyInput, setPublicKeyInput] = useState('');
  const [authorizeError, setAuthorizeError] = useState<KeyError | null>(null);
  const [authorizationSuccess, setAuthorizationSuccess] = useState(false);

  const handleSignCertificate = () => {
    if (!publicKeyInput.trim()) {
      return;
    }

    try {
      const publicKey = createRawPublicKey(publicKeyInput.trim());
      signCertificateState.signCertificate(publicKey);
    } catch (error) {
      console.error('Failed to parse public key:', error);
    }
  };

  const handleAuthorizePublicKey = () => {
    if (!publicKeyInput.trim()) {
      setAuthorizeError('INVALID_PUBLIC_KEY');
      return;
    }

    setAuthorizationSuccess(false);
    setAuthorizeError(null);

    try {
      const publicKey = createRawPublicKey(publicKeyInput.trim());
      authorizePublicKeyState.authorizePublicKey({
        instanceId,
        publicKey,
      });
    } catch (error) {
      console.error('Failed to authorize public key:', error);
      setAuthorizeError('UNKNOWN_ERROR');
    }
  };

  const handleCopyCertificate = (certificate: SignedCertificate) => {
    navigator.clipboard.writeText(certificate).catch((error) => {
      console.error('Failed to copy certificate:', error);
    });
  };

  // 認可成功時の処理
  if (authorizePublicKeyState.isAuthorized && !authorizationSuccess) {
    setAuthorizationSuccess(true);
    setTimeout(() => {
      setPublicKeyInput('');
      authorizePublicKeyState.reset();
      setAuthorizationSuccess(false);
    }, 2000);
  }

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>SSH鍵の登録</h3>

      {/* Step 1: 公開鍵入力 */}
      <section className={styles.section}>
        <h4 className={styles.sectionTitle}>ステップ1: 公開鍵の入力</h4>
        <textarea
          className={styles.textarea}
          placeholder="ssh-rsa AAAA...（OpenSSH形式の公開鍵を貼り付けてください）"
          value={publicKeyInput}
          onChange={(e) => {
            setPublicKeyInput(e.target.value);
            setAuthorizeError(null);
            setAuthorizationSuccess(false);
          }}
          disabled={signCertificateState.isLoading || authorizePublicKeyState.isLoading}
        />
      </section>

      {/* Step 2: 証明書署名 */}
      <section className={styles.section}>
        <h4 className={styles.sectionTitle}>ステップ2: SSH証明書の生成</h4>
        <button
          className={styles.button}
          onClick={handleSignCertificate}
          disabled={
            !publicKeyInput.trim() ||
            signCertificateState.isLoading ||
            authorizePublicKeyState.isLoading
          }
        >
          {signCertificateState.isLoading ? '生成中...' : '証明書を生成'}
        </button>

        {signCertificateState.error && (
          <div className={styles.error}>
            ❌ 証明書生成失敗: {formatKeyError(signCertificateState.error)}
          </div>
        )}

        {signCertificateState.certificate && (
          <div className={styles.certificateResult}>
            <div className={styles.certificateBox}>
              <textarea
                className={styles.textarea}
                value={signCertificateState.certificate}
                readOnly
              />
            </div>
            <button
              className={styles.buttonSecondary}
              onClick={() => handleCopyCertificate(signCertificateState.certificate!)}
            >
              📋 コピー
            </button>
          </div>
        )}
      </section>

      {/* Step 3: インスタンスへの認可 */}
      <section className={styles.section}>
        <h4 className={styles.sectionTitle}>ステップ3: インスタンスへの簿可</h4>
        <p className={styles.description}>
          公開鍵をこのインスタンス上に登録します。
        </p>
        <button
          className={`${styles.button} ${styles.authorizeButton}`}
          onClick={handleAuthorizePublicKey}
          disabled={
            !publicKeyInput.trim() ||
            signCertificateState.isLoading ||
            authorizePublicKeyState.isLoading
          }
        >
          {authorizePublicKeyState.isLoading ? '登録中...' : 'インスタンスに登録'}
        </button>

        {authorizeError && (
          <div className={styles.error}>
            ❌ 認可失敗: {formatKeyError(authorizeError)}
          </div>
        )}

        {authorizationSuccess && (
          <div className={styles.success}>
            ✅ インスタンスへの登録に成功しました！
          </div>
        )}
      </section>
    </div>
  );
};

/**
 * KeyError型をユーザーフレンドリーなメッセージに変換
 */
function formatKeyError(error: KeyError): string {
  const errorMessages: Record<KeyError, string> = {
    INVALID_PUBLIC_KEY: '無効なSSH公開鍵形式です',
    UNAUTHORIZED: '認可されていません',
    RESOURCE_NOT_FOUND: 'リソースが見つかりません',
    NETWORK_ERROR: 'ネットワークエラーが発生しました',
    SERVER_ERROR: 'サーバーエラーが発生しました',
    UNKNOWN_ERROR: '不明なエラーが発生しました',
  };
  return errorMessages[error] || error;
}
