import apiClient from '../apiClient';

import {OpenaiAccount} from '#/entity';

export enum ClaudeAccountApi {
  list = '/openai-account/list',
  add = '/openai-account/add',
  search = '/openai-account/search',
  delete = '/openai-account/delete',
  update = '/openai-account/update',
  statistic = '/openai-account/statistic',
  chatAuth = '/auth',
  disable = '/openai-account/disable',
  enable = '/openai-account/enable',
}

const getAccountList = () => apiClient.get<OpenaiAccount[]>({url: ClaudeAccountApi.list});
const addAccount = (data: OpenaiAccount) => apiClient.post({url: ClaudeAccountApi.add, data});
const updateAccount = (data: OpenaiAccount) => apiClient.post({url: ClaudeAccountApi.update, data});
const deleteAccount = (data: OpenaiAccount) => apiClient.post({url: ClaudeAccountApi.delete, data});
const searchAccount = (tokenId?: number) => apiClient.post({
  url: ClaudeAccountApi.search, data: {
    tokenId,
  }
});
const chatAuthAccount = (type: number, accountId:number) => apiClient.post({
    url: ClaudeAccountApi.chatAuth,
    data: {
      type: type,
      accountId: accountId,
      password: 'auth',
    }
});

const disableAccount = (data: OpenaiAccount) => apiClient.post({url: ClaudeAccountApi.disable, data});

const enableAccount = (data: OpenaiAccount) => apiClient.post({url: ClaudeAccountApi.enable, data});

type AccountStatistic = {
  series: ApexAxisChartSeries;
  categories: string[]
}


const getAccountStatistic = (tokenId: number) => apiClient.post<AccountStatistic>({
  url: ClaudeAccountApi.statistic,
  data: {tokenId},
});

export default {
  getAccountList,
  addAccount,
  updateAccount,
  searchAccount,
  deleteAccount,
  getAccountStatistic,
  chatAuthAccount,
  disableAccount,
  enableAccount,
};
