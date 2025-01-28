import {
  LoginSchema,
  LoginState,
  RegisterSchema,
  RegisterState,
  RegisterType,
} from "@/lib/types/auth";

export const login = async (state: LoginState, data: FormData) => {
  const obj = Object.fromEntries(data);

  const validatedFields = LoginSchema.safeParse(obj);

  // If any form fields are invalid, return early
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }

  // eslint-disable-next-line no-console
  return {
    errors: {
      name: ["error"],
    },
  };
};

export const register = async (
  preventState: RegisterState,
  formData: FormData,
): Promise<RegisterState> => {
  const obj = Object.fromEntries(formData);

  const validatedFields = RegisterSchema.safeParse(obj);
  const state = validatedFields.data;

  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      data: state || (obj as RegisterType),
    };
  }

  return {
    message: "register success",
  };
};
