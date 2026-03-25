import type { Instance } from "@/domain/compute/compute.model";
import type { ComputeError, IComputeRepository } from "@/domain/compute/compute.repository";

export type ListResult = 
  | { type: 'success'; instances: Instance[] | null }
  | { type: 'fetch_failed'; error: ComputeError };

export type DashboardDeps = {
  computeRepo: IComputeRepository;
}

export interface IListMineUseCase {
  execute(): Promise<ListResult>;
}

export const listMine = ({
  computeRepo,
}: DashboardDeps): IListMineUseCase => ({
  execute: async () => {
    const res = await computeRepo.listMine();
    
    if (!res.success) {
      return { type: 'fetch_failed', error: res.error };
    }

    return { type: 'success', instances: res.data };
  }
})