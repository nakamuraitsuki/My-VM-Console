import { useMutation } from '@tanstack/react-query';
import type { RawPublicKey, SignedCertificate } from '@/domain/compute/key.model';
import type { KeyError } from '@/domain/compute/key.repository';
import type { ISignCertificateUseCase } from '@/application/compute/key.usecase';

/**
 * 公開鍵署名の状態管理
 */
export type UseSignCertificateState = {
  certificate: SignedCertificate | null;
  error: KeyError | null;
  isLoading: boolean;
};

export const useSignCertificate = (useCase: ISignCertificateUseCase) => {
  const mutation = useMutation({
    mutationFn: async (publicKey: RawPublicKey) => {
      const result = await useCase.execute(publicKey);
      
      if (result.type === 'failed') {
        throw new Error(`Certificate signing failed: ${result.error}`);
      }

      return result.certificate;
    },
  });

  return {
    certificate: mutation.data ?? null,
    error: (mutation.error?.message.includes('signing failed') 
      ? mutation.error.message.split(': ')[1] as KeyError 
      : null),
    isLoading: mutation.isPending,
    signCertificate: mutation.mutate,
    reset: mutation.reset,
  };
};
