import { BasicStatus, PermissionType } from './enum';

export interface UserToken {
  accessToken?: string
}

export interface UserInfo {
  id: string;
  email: string;
  username: string;
  password?: string;
  avatar?: string;
  role?: Role;
  status?: BasicStatus;
  permissions?: Permission[];
}

export interface Organization {
  id: string;
  name: string;
  status: 'enable' | 'disable';
  desc?: string;
  order?: number;
  children?: Organization[];
}

export interface Permission {
  id: string;
  parentId: string;
  name: string;
  label: string;
  type: PermissionType;
  route: string;
  status?: BasicStatus;
  order?: number;
  icon?: string;
  component?: string;
  hide?: boolean;
  frameSrc?: string;
  newFeature?: boolean;
  children?: Permission[];
}

export interface Role {
  id: string;
  name: string;
  label: string;
  status: BasicStatus;
  order?: number;
  desc?: string;
  permission?: Permission[];
}

export interface User {
  id: number;
  uniqueName: string;
  password: string;
  enable: 0 | 1;
  openai: 0 | 1;
  openaiToken?: number;
  claude: 0 | 1;
  expirationTime?: string;
  createTime?: string;
  updateTime?: string;
}

export interface OpenaiToken {
  id: number;
  tokenName: string;
  plusSubscription?: number;
  refreshToken: string;
  accessToken?: string;
  expireAt?: string;
  createTime?: string;
  updateTime?: string;
}

export interface OpenaiAccount {
  id?: number;
  userId: number;
  account: string;
  status: 1 | 0;
  expirationTime?: string;
  tokenId: number;
  gpt35Limit: number;
  gpt4Limit: number;
  gpt4oLimit: number;
  gpt4oMiniLimit: number;
  o1Limit: number;
  o1MiniLimit: number;
  showConversations: 1 | 0;
  temporaryChat: 0 | 1;
  shareToken?: string;
  expireAt?: string;
  createTime?: string;
  updateTime?: string;
}

export interface ClaudeToken {
  id: number;
  tokenName: string;
  sessionToken: string;
  expireAt?: string;
  createTime?: string;
  updateTime?: string;
}

export interface ClaudeAccount {
  id?: number;
  userId: number;
  tokenId: number;
  account: string;
  status: 1 | 0;
  createTime?: string;
  updateTime?: string;
}
