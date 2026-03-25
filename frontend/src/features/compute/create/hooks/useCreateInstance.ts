import { useState } from 'react';
import { requestCreateInstance } from '@/application/compute/create.usecase';
import type { Instance } from '@/domain/compute/compute.model';
import type { ComputeError, CreateInstanceRequest } from '@/domain/compute/compute.repository';
import { useServices } from '@/context/ServiceContext';

interface UseCreateInstanceState {
  isLoading: boolean;
  error: ComputeError | null;
  data: Instance | null;
}

/**
 * インスタンス作成ロジックのフック
 */
export const useCreateInstance = () => {
  const { computeRepository } = useServices();
  const [state, setState] = useState<UseCreateInstanceState>({
    isLoading: false,
    error: null,
    data: null,
  });

  const execute = async (request: CreateInstanceRequest) => {
    setState({ isLoading: true, error: null, data: null });

    const usecase = requestCreateInstance({ computeRepo: computeRepository });
    const result = await usecase.execute(request);

    if (result.type === 'success') {
      setState({ isLoading: false, error: null, data: result.instance });
      return result;
    }

    setState({ isLoading: false, error: result.error, data: null });
    return result;
  };

  const reset = () => {
    setState({ isLoading: false, error: null, data: null });
  };

  return {
    ...state,
    execute,
    reset,
  };
};
