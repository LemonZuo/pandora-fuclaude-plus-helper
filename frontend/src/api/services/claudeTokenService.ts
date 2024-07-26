import apiClient from '../apiClient';

import {ClaudeToken} from '#/entity';

export enum ClaudeTokenApi {
  list = '/claude-token/list',
  add = '/claude-token/add',
  update = '/claude-token/update',
  delete = '/claude-token/delete',
  refresh = '/claude-token/refresh',
  search = '/claude-token/search',
}

const getTokenList = () => apiClient.get<ClaudeToken[]>({ url: ClaudeTokenApi.list }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});

const searchTokenList = (tokenName: string) => apiClient.post<ClaudeToken[]>({ url: ClaudeTokenApi.search, data: {tokenName} }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});
export interface ClaudeTokenAddReq {
  id?: number;
  tokenName: string;
  sessionToken: string;
}
export interface taskStatus {
  status: boolean;
}
const addToken = (data: ClaudeTokenAddReq) => apiClient.post({ url: ClaudeTokenApi.add, data });
const updateToken = (data: ClaudeTokenAddReq) => apiClient.post({ url: ClaudeTokenApi.update, data });
const deleteToken = (id: number) => apiClient.post({ url: ClaudeTokenApi.delete, data: { id } });
const refreshToken = (id: number) => apiClient.post({ url: ClaudeTokenApi.refresh, data: { id } })

export default {
  getTokenList,
  searchTokenList,
  addToken,
  updateToken,
  deleteToken,
  refreshToken,
};
