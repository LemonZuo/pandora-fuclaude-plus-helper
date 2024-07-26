import apiClient from '../apiClient';

import {ClaudeAccount} from '#/entity';

export enum ClaudeAccountApi {
  list = '/claude-account/list',
  add = '/claude-account/add',
  search = '/claude-account/search',
  delete = '/claude-account/delete',
  update = '/claude-account/update',
  statistic = '/claude-account/statistic',
  chatAuth = '/auth',
  disable = '/claude-account/disable',
  enable = '/claude-account/enable',
}

const getAccountList = () => apiClient.get<ClaudeAccount[]>({url: ClaudeAccountApi.list});
const addAccount = (data: ClaudeAccount) => apiClient.post({url: ClaudeAccountApi.add, data});
const updateAccount = (data: ClaudeAccount) => apiClient.post({url: ClaudeAccountApi.update, data});
const deleteAccount = (data: ClaudeAccount) => apiClient.post({url: ClaudeAccountApi.delete, data});
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

const disableAccount = (data: ClaudeAccount) => apiClient.post({url: ClaudeAccountApi.disable, data});

const enableAccount = (data: ClaudeAccount) => apiClient.post({url: ClaudeAccountApi.enable, data});

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
