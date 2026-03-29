import type { Result } from '../core/result';
import type { ImageID, Instance, InstanceDetail, SubnetID, VPCID } from './compute.model';

export type ComputeError =
  | 'INVALID_REQUEST'
  | 'UNAUTHORIZED'
  | 'QUOTA_EXCEEDED'
  | 'RESOURCE_NOT_FOUND'
  | 'NETWORK_ERROR'
  | 'SERVER_ERROR'
  | 'UNKNOWN_ERROR';

export interface CreateInstanceRequest {
  name: string;
  imageId: ImageID;
  vpcId?: VPCID;
  subnetId?: SubnetID;
  cpu: number;
  memory: number;
}

export interface IComputeRepository {
  /**
   * インスタンスを作成します
   */
  createInstance(request: CreateInstanceRequest): Promise<Result<Instance, ComputeError>>;
  /**
   * インスタンス列挙
   */
  listMine(): Promise<Result<Instance[] | null, ComputeError>>;
  /**
   * インスタンス詳細取得
   */
  getByID(instanceId: string): Promise<Result<InstanceDetail, ComputeError>>;
}
