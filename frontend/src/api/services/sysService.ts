import apiClient from '../apiClient';

export enum SysApi {
  pandoraUsage= '/sys_info/usage',
  setting = '/sys_info/info',
  info = '/info',
}

export type PandoraUsage = {
  total: number;
  current: number;
  ttl: number
}

const getPandoraUsage = () => apiClient.get<PandoraUsage>({ url: SysApi.pandoraUsage });
const getSetting = () => apiClient.get({ url: SysApi.setting });
const getInfo = () => apiClient.post({ url: SysApi.info });

export default {
  getPandoraUsage,
  getSetting,
  getInfo,
}
