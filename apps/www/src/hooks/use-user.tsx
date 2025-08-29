import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { getJson, postJson } from "@/utils/fetch";
import { Response } from "@/types/api_response";
import { UserResponse } from "@/types/user_response";

const keys = {
  userMe: ["users", "me"] as const,
};

const muKey = {
  uploadImageProfile: ["users", "profile-image"] as const,
};

export const useGetUserMe = () => {
  return useQuery({
    queryKey: keys.userMe,
    queryFn: () => getJson<Response<UserResponse>>("/users/me"),
  });
};

export const useGetImageProfile = () => {
  const { data, ...rest } = useGetUserMe();
  const image = data?.data.data.profile_image?.at(-1)?.image;
  const src = image ? `/api/v1/${image.path}` : undefined;

  return { image, src, ...rest };
};

export const useGetMutateUserMe = () => {
  return useMutation({
    mutationKey: keys.userMe,
    mutationFn: () => getJson<Response<UserResponse>>("/users/me"),
  });
};

// use upload image profile
export const useUploadProfileImage = (onSuccess?: () => void) => {
  const client = useQueryClient();

  return useMutation({
    mutationKey: muKey.uploadImageProfile,
    mutationFn: (file: File) => {
      const formData = new FormData();

      formData.append("file", file);

      return postJson<Response<UserResponse>>(
        "/users/profile-image",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        },
      );
    },
    onSuccess: () => {
      client.invalidateQueries({ queryKey: keys.userMe });
      onSuccess?.();
    },
  });
};
