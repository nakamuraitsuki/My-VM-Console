'use client';

import { useMemo } from "react";
import type { AppRoute } from "@/app/router";
import { useParams } from "react-router";
import { useServices } from "@/context/ServiceContext";
import { signCertificate, authorizePublicKey } from "@/application/compute/key.usecase";
import { useSignCertificate } from "@/features/compute/detail/hooks/useSignCertificate";
import { useAuthorizePublicKey } from "@/features/compute/detail/hooks/useAuthorizePublicKey";
import { KeyRegistrationForm } from "@/features/compute/detail/components/KeyRegistrationForm";
import { useInstanceDetailQuery } from "@/features/compute/detail/hooks/useInstanceDetailQuery";
import styles from "./InstanceDetail.module.css";

export const InstanceDetailPage = () => {
  const { instanceId } = useParams();
  const { keyRepository } = useServices();

  // usecaseのセットアップ
  const signCertificateUseCase = useMemo(
    () => signCertificate({ keyRepo: keyRepository }),
    [keyRepository]
  );

  const authorizePublicKeyUseCase = useMemo(
    () => authorizePublicKey({ keyRepo: keyRepository }),
    [keyRepository]
  );

  // hooksの初期化
  const signCertificateState = useSignCertificate(signCertificateUseCase);
  const authorizePublicKeyState = useAuthorizePublicKey(authorizePublicKeyUseCase);
  const detailQuery = useInstanceDetailQuery(instanceId);

  if (!instanceId) {
    return (
      <div className={styles.container}>
        <h1 className={styles.title}>インスタンス詳細</h1>
        <p className={styles.error}>インスタンスIDが指定されていません</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>インスタンス詳細</h1>
      <p className={styles.lead}>対象インスタンス: {instanceId}</p>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>紐づくサブドメイン</h2>
        {detailQuery.isPending && (
          <p className={styles.muted}>サブドメインを読み込み中です...</p>
        )}
        {detailQuery.isError && (
          <p className={styles.error}>サブドメインの取得に失敗しました。</p>
        )}
        {detailQuery.data && detailQuery.data.detail.subdomains.length === 0 && (
          <p className={styles.muted}>現在、このインスタンスに紐づくサブドメインはありません。</p>
        )}
        {detailQuery.data && detailQuery.data.detail.subdomains.length > 0 && (
          <ul className={styles.list}>
            {detailQuery.data.detail.subdomains.map((subdomain) => (
              <li key={subdomain}>{subdomain}</li>
            ))}
          </ul>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>SSH鍵の登録</h2>
        <KeyRegistrationForm
          instanceId={instanceId}
          signCertificateState={signCertificateState}
          authorizePublicKeyState={authorizePublicKeyState}
        />
      </section>
    </div>
  );
};

export const route: AppRoute = {
  path: "/dashboard/instances/:instanceId",
  element: <InstanceDetailPage />,
  meta: { name: "instance-detail", requiresAuth: true },
};
