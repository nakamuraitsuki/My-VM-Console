// cf. https://zenn.dev/farstep/articles/typescript-branded-types
// branded type によって、別IDの混同などを防止する。 
export type InstanceID = string & { readonly __brand: 'InstanceID' };
export type ImageID = string & { readonly __brand: 'ImageID' };
export type VPCID = string & { readonly __brand: 'VPCID' };
export type SubnetID = string & { readonly __brand: 'SubnetID' };

export type InstanceStatus = 'PENDING' | 'RUNNING' | 'STOPPED' | 'ERROR';

/**
 * コンピュートインスタンスのドメインモデル
 */
export interface Instance {
  id: InstanceID;
  name: string;
  status: InstanceStatus;
  subnetId: SubnetID;
  privateIp: string;
}

export interface InstanceDetail extends Instance {
  subdomains: string[];
}



// 型ガード関数
export const createInstanceID = (id: string): InstanceID => id as InstanceID;
export const createImageID = (id: string): ImageID => id as ImageID;
export const createVPCID = (id: string): VPCID => id as VPCID;
export const createSubnetID = (id: string): SubnetID => id as SubnetID;
