import type { Result } from "../core/result";
import type { RawPublicKey, SignedCertificate } from "./key.model";

type KeyError =
  | 'INVALID_PUBLIC_KEY'
  | 'UNAUTHORIZED'
  | 'RESOURCE_NOT_FOUND'
  | 'NETWORK_ERROR'
  | 'SERVER_ERROR'
  | 'UNKNOWN_ERROR';

export interface IKeyRepository {
  /**
     * 公開鍵をバックエンドに送り、署名済み証明書を取得する
     */
  signCertificate(publicKey: RawPublicKey): Promise<Result<SignedCertificate, KeyError>>;

  /**
   * 指定したインスタンスに公開鍵を配置（Authorize）する
   */
  authorizePublicKey(instanceId: string, publicKey: RawPublicKey): Promise<Result<void, KeyError>>;
}
