import {useMutation, useQueryClient} from "@tanstack/react-query";
import userService from "@/api/services/userService.ts";
import {message} from "antd";

export const useAddUserMutation = () => {
  const client = useQueryClient();
  return useMutation(userService.addUser, {
    onSuccess: () => {
      /* onSuccess */
      message.success('Add Token Success')
      client.invalidateQueries(['users']);
    },
  });
}

export const useUpdateUserMutation = () => {
  const client = useQueryClient();
  return useMutation(userService.updateUser, {
    onSuccess: () => {
      /* onSuccess */
      message.success('Update Token Success')
      client.invalidateQueries(['users']);
    },
  });
}

export const useDeleteUserMutation = () => {
  const client = useQueryClient();
  return useMutation(userService.deleteUser, {
    onSuccess: () => {
      /* onSuccess */
      message.success('Delete Token Success')
      client.invalidateQueries(['users']);
    }
  });
}
