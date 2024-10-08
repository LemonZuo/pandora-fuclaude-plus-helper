import {useMutation, useQueryClient} from "@tanstack/react-query";
import accountService from "@/api/services/claudeAccountService.ts";
import {message} from "antd";


export const useAddAccountMutation = () => {
  const client = useQueryClient();
  return useMutation(accountService.addAccount, {
    onSuccess: () => {
      /* onSuccess */
      client.invalidateQueries(['accounts']);
      message.success('Success')
    },
  });
}

export const useUpdateAccountMutation = () => {
  const client = useQueryClient();
  return useMutation(accountService.updateAccount, {
    onSuccess: () => {
      /* onSuccess */
      client.invalidateQueries(['shareList']);
      message.success('Success')
    },
  });
}

export const useDeleteAccountMutation = () => {
  const client = useQueryClient();
  return useMutation(accountService.deleteAccount, {
    onSuccess: () => {
      /* onSuccess */
      client.invalidateQueries(['shareList']);
      message.success('Success')
    },
  })
}

export const useDisableAccountMutation = () => {
  const client = useQueryClient();
  return useMutation(accountService.disableAccount, {
    onSuccess: () => {
      /* onSuccess */
      client.invalidateQueries(['shareList']);
      // message.success('Success')
    },
  })
}

export const useEnableAccountMutation = () => {
  const client = useQueryClient();
  return useMutation(accountService.enableAccount, {
    onSuccess: () => {
      /* onSuccess */
      client.invalidateQueries(['shareList']);
      // message.success('Success')
    },
  })
}

export default {
  useAddAccountMutation,
  useDeleteAccountMutation,
}

