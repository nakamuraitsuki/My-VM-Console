export type userID = string;

export type UserStatus =
  | 'pending'
  | 'initializing'
  | 'active'
  | 'failed';

export type FailedPhase =
  | 'failed in pending'
  | 'failed in initializing';

export const Permissions = {
  All: '*',
  InstanceRead: 'instance:read',
  InstanceCreate: 'instance:create',
  InstanceStop: 'instance:stop',
  InstanceStopAll: 'instance:stop:all',
  InstanceUpdate: 'instance:update',
  InstanceDelete: 'instance:delete',
  InstanceDeleteAll: 'instance:delete:all',
  NetworkManage: 'network:manage',
} as const;

export type Permission = typeof Permissions[keyof typeof Permissions];

export interface UsageQuota {
  readonly maxInstance: number;
  readonly maxCPU: number;
  readonly maxMemory: number;
}

export interface UserProps {
  readonly id: string;
  readonly displayName: string;
  readonly permissions: Permission[];
  readonly quota: UsageQuota;
  readonly status: UserStatus;
  readonly errorPhase?: FailedPhase;
}

export interface User {
  id: userID
  displayName: string;
  profileImageURL: string;
  permissions: Permission[];
  quota: UsageQuota;
  status: UserStatus;
  errorPhase?: FailedPhase;
}
