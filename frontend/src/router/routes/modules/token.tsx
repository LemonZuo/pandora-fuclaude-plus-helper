import { Suspense, lazy } from 'react';
import { Navigate, Outlet } from 'react-router-dom';

import { SvgIcon } from '@/components/icon';
import { CircleLoading } from '@/components/loading';

import { AppRouteObject } from '#/router';

const OpenaiAccountPage = lazy(() => import(`@/pages/token/openai/account`));
const ClaudeTokenPage = lazy(() => import(`@/pages/token/claude/token`));
const ClaudeAccountPage = lazy(() => import(`@/pages/token/claude/account`));
const UserPage = lazy(() => import(`@/pages/token/user`));

const token: AppRouteObject = {
  order: 10,
  path: 'openai-token',
  element: (
    <Suspense fallback={<CircleLoading />}>
      <Outlet />
    </Suspense>
  ),
  meta: {
    label: 'sys_info.menu.dashboard',
    icon: <SvgIcon icon="ic-analysis" className="ant-menu-item-icon" size="24" />,
    key: '/token',
  },
  children: [
    {
      index: true,
      element: <Navigate to="openai-account" replace />,
    },
    {
      path: 'openai-account',
      element: <OpenaiAccountPage />,
      meta: { label: 'sys_info.menu.openai-account', key: '/token/openai-account' },
    },
    {
      path: 'openai-account',
      element: <ClaudeTokenPage />,
      meta: { label: 'sys_info.menu.claude-token', key: '/token/claude-token' },
    },
    {
      path: 'openai-account',
      element: <ClaudeAccountPage />,
      meta: { label: 'sys_info.menu.claude-account', key: '/token/claude-account' },
    },
    {
      path: 'user',
      element: <UserPage />,
      meta: { label: 'sys_info.menu.account', key: '/token/user' },
    },
  ],
};

export default token;
