import { listMine } from "@/application/compute/list-mine.usecase";
import { useAuth } from "@/context/AuthContext";
import { useServices } from "@/context/ServiceContext";
import type { Instance } from "@/domain/compute/compute.model";
import { useQuery } from "@tanstack/react-query";

export const useListMineQuery = () => {
  const { session } = useAuth();
  const { computeRepository } = useServices();

  const execute = listMine({ computeRepo: computeRepository });

  return useQuery<{ insts: Instance[] | null }, Error>({
    queryKey: ["myInstances"],
    queryFn: async () => {
      const res = await execute.execute();

      if (res.type !== "success") {
        throw new Error("Failed to fetch my instances");
      }

      return { insts: res.instances };
    },
    enabled: session.status === "authenticated", // ログイン時のみクエリを有効化
    // キャッシュなどが必要な場合は追記
  });
};
