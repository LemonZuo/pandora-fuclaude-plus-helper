import React from 'react';
import {useQuery} from "@tanstack/react-query";
import sysService from "@/api/services/sysService.ts";

interface SystemInfoResponse {
  startTime: string;
  status: boolean;
  systemName: string;
  version: string;
}

interface SystemInfoProps {
  version: string;
}

const VersionDisplay: React.FC<SystemInfoProps> = ({version}) => (
  <div className="flex cursor-pointer items-center justify-center rounded-md p-2 hover:bg-hover font-bold transition-colors duration-500 ease-out">
    {version}
  </div>
);

export const SystemInfo: React.FC = () => {
  const {data: sysInfo, isLoading} = useQuery<SystemInfoResponse, Error>({
    queryKey: ['sysInfo'],
    queryFn: sysService.getInfo,
    staleTime: Infinity, // 假设版本信息不经常变化
  });

  if (isLoading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  if (!sysInfo) {
    return null;
  }

  return <VersionDisplay version={sysInfo.version}/>;
};
