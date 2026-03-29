import { useMutation } from '@tanstack/react-query';
import type { RawPublicKey } from '@/domain/compute/key.model';
import type { KeyError } from '@/domain/compute/key.repository';
import type { IAuthorizePublicKeyUseCase } from '@/application/compute/key.usecase';

/**
 * 公開鍵認可の状態管理
 */
export type UseAuthorizePublicKeyState = {
  isAuthorized: boolean;
  error: KeyError | null;
  isLoading: boolean;
};

export const useAuthorizePublicKey = (useCase: IAuthorizePublicKeyUseCase) => {
  const mutation = useMutation({
    mutationFn: async ({
      instanceId,
      publicKey,
    }: {
      instanceId: string;
      publicKey: RawPublicKey;
    }) => {
      const result = await useCase.execute(instanceId, publicKey);

      if (result.type === 'failed') {
        throw new Error(`Authorization failed: ${result.error}`);
      }

      return true;
    },
  });

  return {
    isAuthorized: mutation.isSuccess,
    error: (mutation.error?.message.includes('failed')
      ? mutation.error.message.split(': ')[1] as KeyError
      : null),
    isLoading: mutation.isPending,
    authorizePublicKey: mutation.mutate,
    reset: mutation.reset,
  };
};
