import type { InstanceDetail } from "@/domain/compute/compute.model";
import type { ComputeError, IComputeRepository } from "@/domain/compute/compute.repository";

export type GetDetailResult =
  | { type: "success"; detail: InstanceDetail }
  | { type: "fetch_failed"; error: ComputeError };

export type GetDetailDeps = {
  computeRepo: IComputeRepository;
};

export interface IGetDetailUseCase {
  execute(instanceId: string): Promise<GetDetailResult>;
}

export const getDetail = ({
  computeRepo,
}: GetDetailDeps): IGetDetailUseCase => ({
  execute: async (instanceId: string) => {
    const res = await computeRepo.getByID(instanceId);

    if (!res.success) {
      return { type: "fetch_failed", error: res.error };
    }

    return { type: "success", detail: res.data };
  },
});
