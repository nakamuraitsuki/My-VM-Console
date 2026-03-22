import type { IComputeRepository, ComputeError, CreateInstanceRequest } from '../../domain/compute/compute.repository';
import type { Instance } from '../../domain/compute/compute.model';

export type RequestCreateInstanceResult =
  | { type: 'success'; instance: Instance }
  | { type: 'create_failed'; error: ComputeError };

export type RequestCreateInstanceDeps = {
  computeRepo: IComputeRepository;
};

export type RequestCreateInstanceParams = CreateInstanceRequest;

export interface IRequestCreateInstanceUseCase {
  execute(params: RequestCreateInstanceParams): Promise<RequestCreateInstanceResult>;
}

export const requestCreateInstance = ({
  computeRepo,
}: RequestCreateInstanceDeps): IRequestCreateInstanceUseCase => ({
  execute: async (params) => {
    // アプリケーション層でRepository結果をユースケース結果へ写像する。
    const res = await computeRepo.createInstance(params);

    if (!res.success) {
      return { type: 'create_failed', error: res.error };
    }

    return { type: 'success', instance: res.data };
  },
});
