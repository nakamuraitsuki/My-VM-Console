import type { RawPublicKey, SignedCertificate } from "@/domain/compute/key.model";
import type { IKeyRepository, KeyError } from "@/domain/compute/key.repository";

/**
 * 鍵署名の実行結果
 */
export type SignCertificateResult =
  | { type: 'success'; certificate: SignedCertificate }
  | { type: 'failed'; error: KeyError };

/**
 * 鍵認可の実行結果
 */
export type AuthorizePublicKeyResult =
  | { type: 'success' }
  | { type: 'failed'; error: KeyError };

/**
 * Key管理のビジネスロジック
 */
export type KeyDeps = {
  keyRepo: IKeyRepository;
};

export interface ISignCertificateUseCase {
  execute(publicKey: RawPublicKey): Promise<SignCertificateResult>;
}

export interface IAuthorizePublicKeyUseCase {
  execute(instanceId: string, publicKey: RawPublicKey): Promise<AuthorizePublicKeyResult>;
}

/**
 * 公開鍵の署名を実行
 */
export const signCertificate = ({
  keyRepo,
}: KeyDeps): ISignCertificateUseCase => ({
  execute: async (publicKey: RawPublicKey) => {
    const result = await keyRepo.signCertificate(publicKey);

    if (!result.success) {
      return { type: 'failed', error: result.error };
    }

    return { type: 'success', certificate: result.data };
  },
});

/**
 * インスタンスに公開鍵を認可
 */
export const authorizePublicKey = ({
  keyRepo,
}: KeyDeps): IAuthorizePublicKeyUseCase => ({
  execute: async (instanceId: string, publicKey: RawPublicKey) => {
    const result = await keyRepo.authorizePublicKey(instanceId, publicKey);

    if (!result.success) {
      return { type: 'failed', error: result.error };
    }

    return { type: 'success' };
  },
});
