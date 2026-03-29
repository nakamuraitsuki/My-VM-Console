import { getDetail } from "@/application/compute/get-detail.usecase";
import { useAuth } from "@/context/AuthContext";
import { useServices } from "@/context/ServiceContext";
import type { InstanceDetail } from "@/domain/compute/compute.model";
import { useQuery } from "@tanstack/react-query";

export const useInstanceDetailQuery = (instanceId?: string) => {
  const { session } = useAuth();
  const { computeRepository } = useServices();

  const execute = getDetail({ computeRepo: computeRepository });

  return useQuery<{ detail: InstanceDetail }, Error>({
    queryKey: ["instanceDetail", instanceId],
    queryFn: async () => {
      if (!instanceId) {
        throw new Error("instance id is missing");
      }

      const res = await execute.execute(instanceId);
      if (res.type !== "success") {
        throw new Error("Failed to fetch instance detail");
      }

      return { detail: res.detail };
    },
    enabled: Boolean(instanceId) && session.status === "authenticated",
  });
};
