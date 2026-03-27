import type { RawPublicKey, SignedCertificate } from "@/domain/compute/key.model";
import type { IKeyRepository, KeyError } from "@/domain/compute/key.repository";
import { failure, success, type Result } from "@/domain/core/result";

/**
 * 鍵リポジトリの実装
 */
export class KeyGateway implements IKeyRepository {
  async signCertificate(publicKey: RawPublicKey): Promise<Result<SignedCertificate, KeyError>> {
    try {
      const response = await fetch('/api/ssh-ca/sign', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ publicKey }),
      });

      if (!response.ok) {
        // HTTPエラーをKeyErrorにマッピング
        const errorType: KeyError = response.status === 400 ? 'INVALID_PUBLIC_KEY' : 'SERVER_ERROR';
        return failure(errorType);
      }

      const data = await response.json();
      const signedCert = data.signedCertificate as string;
      return success(signedCert as SignedCertificate);
    } catch (error) {
      console.error('Error signing certificate:', error);
      return failure('NETWORK_ERROR');
    }
  }

  async authorizePublicKey(instanceId: string, publicKey: RawPublicKey): Promise<Result<void, KeyError>> {
    try {
      const response = await fetch(`/api/ssh-ca/instances/authorize-key`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          instanceId,
          publicKey
        }),
      });

      if (!response.ok) {
        const errorType: KeyError = response.status === 400 ? 'INVALID_PUBLIC_KEY' : 'SERVER_ERROR';
        return failure(errorType);
      }

      return success<void>(void 0);
    } catch (error) {
      console.error('Error authorizing public key:', error);
      return failure('NETWORK_ERROR');

    }
  }
}