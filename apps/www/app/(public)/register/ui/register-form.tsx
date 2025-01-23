"use client";
import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Divider,
  Form,
  Input,
} from "@heroui/react";
import { Send } from "lucide-react";
import React, { useActionState } from "react";

import { register } from "@/actions/authentication";

export const RegisterForm = () => {
  const [state, action, pending] = useActionState(register, {
    data: undefined,
  });

  return (
    <React.Fragment>
      <div className="flex justify-center items-center">
        <Card fullWidth className="max-w-sm">
          <CardHeader>Register</CardHeader>
          <CardBody className="flex flex-col gap-2">
            <Form action={action}>
              <Input
                defaultValue={state?.data?.username}
                label="Username"
                name="username"
              />
              <Input
                errorMessage={state?.errors?.email?.join(",")}
                id="email"
                isInvalid={!!state?.errors?.email}
                label="Email"
                name="email"
              />
              <Input
                errorMessage={state?.errors?.password?.join(",")}
                id="password"
                isInvalid={!!state?.errors?.password}
                label="Password"
                name="password"
                type="password"
              />
              <Input
                errorMessage={state?.errors?.cf_password?.join(",")}
                id="cf_password"
                isInvalid={!!state?.errors?.cf_password}
                label="Confirm Password"
                name="cf_password"
                type="password"
              />
              <Divider />
              <Input
                errorMessage={state?.errors?.first_name?.join(",")}
                id="first_name"
                isInvalid={!!state?.errors?.first_name}
                label="First Name"
                name="first_name"
              />
              <Input
                errorMessage={state?.errors?.last_name?.join(",")}
                id="last_name"
                isInvalid={!!state?.errors?.last_name}
                label="Last Name"
                name="last_name"
              />
              <Button
                fullWidth
                color="primary"
                isLoading={pending}
                startContent={<Send />}
                type="submit"
              >
                Register
              </Button>
            </Form>
          </CardBody>
          <CardFooter />
        </Card>
      </div>
    </React.Fragment>
  );
};
