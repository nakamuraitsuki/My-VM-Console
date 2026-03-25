import type { Result } from '../core/result';
import type { ImageID, Instance, SubnetID, VPCID } from './compute.model';

export type ComputeError =
  | 'INVALID_REQUEST'
  | 'UNAUTHORIZED'
  | 'QUOTA_EXCEEDED'
  | 'RESOURCE_NOT_FOUND'
  | 'NETWORK_ERROR'
  | 'SERVER_ERROR'
  | 'UNKNOWN_ERROR';

/**
 * インスタンス作成のリクエストモデル
 */
export interface CreateInstanceRequest {
  name: string;
  imageId: ImageID;
  vpcId?: VPCID;
  subnetId?: SubnetID;
  cpu: number;
  memory: number;
}

/**
 * コンピュートリソースのリポジトリインターフェース
 */
export interface IComputeRepository {
  /**
   * インスタンスを作成します
   */
  createInstance(request: CreateInstanceRequest): Promise<Result<Instance, ComputeError>>;
}
