import type { CommonModel, Image } from "@/shared/types";

export interface Role extends CommonModel {
  name: string;
}

export interface ProfileImage extends CommonModel {
  user_id: number;
  image_id: number;
  image: Image;
}

export interface UserResponse extends CommonModel {
  username: string;
  email: string;
  email_verifyed: boolean;
  full_name: string;
  role_id: number;
  profile_image?: ProfileImage[];
  role?: Role;
}
