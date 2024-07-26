import apiClient from '../apiClient';

import {User} from '#/entity';

export enum UserApi {
  list = '/user/list',
  add = '/user/add',
  update = '/user/update',
  delete = '/user/delete',
  refresh = '/user/refresh',
  search = '/user/search',
}

const getUserList = () => apiClient.get<User[]>({ url: UserApi.list }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});

const searchUserList = (uniqueName: string) => apiClient.post<User[]>({ url: UserApi.search, data: {uniqueName} }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});
export interface UserAddReq {
  id?: number;
  uniqueName: string;
  password: string;
  enable: 0 | 1;
  openai: 0 | 1;
  openaiToken?: number;
  claude: 0 | 1;
  claudeToken?: number;
  expirationTime?: string;
  createTime?: string;
  updateTime?: string;
}
export interface taskStatus {
  status: boolean;
}
const addUser = (data: UserAddReq) => apiClient.post({ url: UserApi.add, data });
const updateUser = (data: UserAddReq) => apiClient.post({ url: UserApi.update, data });
const deleteUser = (id: number) => apiClient.post({ url: UserApi.delete, data: { id } });
const refreshUser = (id: number) => apiClient.post({ url: UserApi.refresh, data: { id } })

export default {
  getUserList,
  searchUserList,
  addUser,
  updateUser,
  deleteUser,
  refreshUser,
};
