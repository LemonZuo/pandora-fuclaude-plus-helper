import apiClient from '../apiClient';
import {UserInfo} from "#/entity.ts";

// import {UserInfo, UserToken} from '#/entity';

export interface SignInReq {
  type: number;
  accountId?: number;
  password?: string;
  token?: string
}

// export type SignInRes = UserToken & {user: UserInfo};
export type SignInRes = {
  type: number
  accessToken: string
  user: UserInfo
  loginUrl: string
};

export enum AuthApi {
  SignIn = 'auth',
  Logout = '/auth/logout',
}

const signin = (data: SignInReq) => apiClient.post<SignInRes>({ url: AuthApi.SignIn, data });
const logout = () => apiClient.get({ url: AuthApi.Logout });

export default {
  signin,
  logout,
};
