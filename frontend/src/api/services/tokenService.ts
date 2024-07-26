import apiClient from '../apiClient';

import {OpenaiToken} from '#/entity';

export enum OpenaiTokenApi {
  list = '/openai-token/list',
  add = '/openai-token/add',
  update = '/openai-token/update',
  delete = '/openai-token/delete',
  refresh = '/openai-token/refresh',
  search = '/openai-token/search',
}

const getTokenList = () => apiClient.get<OpenaiToken[]>({ url: OpenaiTokenApi.list }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});

const searchTokenList = (tokenName: string) => apiClient.post<OpenaiToken[]>({ url: OpenaiTokenApi.search, data: {tokenName} }).then((res) => {
  // 将shareList转为json对象
  // res.forEach((item) => {
  //   if (item.shareList) {
  //     item.shareList = JSON.parse(item.shareList);
  //   }
  // });
  return res;
});
export interface OpenaiTokenAddReq {
  id?: number;
  tokenName: string;
  refreshToken: string;
}
export interface taskStatus {
  status: boolean;
}
const addToken = (data: OpenaiTokenAddReq) => apiClient.post({ url: OpenaiTokenApi.add, data });
const updateToken = (data: OpenaiTokenAddReq) => apiClient.post({ url: OpenaiTokenApi.update, data });
const deleteToken = (id: number) => apiClient.post({ url: OpenaiTokenApi.delete, data: { id } });
const refreshToken = (id: number) => apiClient.post({ url: OpenaiTokenApi.refresh, data: { id } })

export default {
  getTokenList,
  searchTokenList,
  addToken,
  updateToken,
  deleteToken,
  refreshToken,
};
