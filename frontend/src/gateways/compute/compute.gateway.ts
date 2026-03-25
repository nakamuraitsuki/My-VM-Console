import { apiClient } from '../../api/client';
import type { Result } from '../../domain/core/result';
import { success, failure } from '../../domain/core/result';
import type { IComputeRepository, ComputeError, CreateInstanceRequest } from '../../domain/compute/compute.repository';
import type { Instance } from '../../domain/compute/compute.model';
import {
  createInstanceID,
  createSubnetID,
  type InstanceStatus,
} from '../../domain/compute/compute.model';

/**
 * DTO（キャメルケース）
 * Snake -> Camel 変換は、API Clientで済ませておく
 */
interface CreateInstanceApiResponse {
  instanceId: string;
  name: string;
  status: string;
  subnetId: string;
  privateIp: string;
}
// Frontend の Domain はIDなどに厳密なガードがかかっているので、中間体が必要
interface ListMineResponseItem {
  id: string;
  name: string;
  status: string;
  subnetId: string;
  privateIp: string;
}

/**
 * コンピュートリポジトリの実装
 */
export class ComputeGateway implements IComputeRepository {
  async createInstance(request: CreateInstanceRequest): Promise<Result<Instance, ComputeError>> {
    try {
      const response = await apiClient.post<CreateInstanceApiResponse>(
        '/api/computes/instances',
        {
          name: request.name,
          image_id: request.imageId,
          vpc_id: request.vpcId,
          subnet_id: request.subnetId,
          cpu: request.cpu,
          memory: request.memory,
        }
      );

      const data = response.data;
      const instance: Instance = {
        id: createInstanceID(data.instanceId),
        name: data.name,
        status: data.status as InstanceStatus,
        subnetId: createSubnetID(data.subnetId),
        privateIp: data.privateIp,
      };

      return success(instance);
    } catch (error) {
      console.error('Failed to create instance:', error);
      return failure(this.mapError(error));
    }
  }

  async listMine(): Promise<Result<Instance[] | null, ComputeError>> {
    try {
      const response = await apiClient.get<ListMineResponseItem[]>('/api/users/me/instances');
      if (!response.data || response.data.length === 0) {
        return success(null); // インスタンスがない場合はnullを返す
      }
      const instances = response.data.map((data) => ({
        id: createInstanceID(data.id),
        name: data.name,
        status: data.status as InstanceStatus,
        subnetId: createSubnetID(data.subnetId),
        privateIp: data.privateIp,
      }));
      return success(instances);
    } catch (error) {
      console.error('Failed to list instances:', error);
      return failure(this.mapError(error));
    }
  }

  private mapError(error: unknown): ComputeError {
    if (error instanceof Error) {
      if (error.message.includes('401') || error.message.includes('Unauthorized')) {
        return 'UNAUTHORIZED';
      }
      if (error.message.includes('400') || error.message.includes('Bad Request')) {
        return 'INVALID_REQUEST';
      }
      if (error.message.includes('409')) {
        return 'QUOTA_EXCEEDED';
      }
      if (error.message.includes('404')) {
        return 'RESOURCE_NOT_FOUND';
      }
      if (error.message.includes('5')) {
        return 'SERVER_ERROR';
      }
    }
    return 'NETWORK_ERROR';
  }
}
