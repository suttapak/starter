import { Button, Form, Input } from "@nextui-org/react";
import { createLazyFileRoute } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";

export const Route = createLazyFileRoute("/login")({
  component: RouteComponent,
});
type FormData = {
  username_email: string;
  password: string;
};

function RouteComponent() {
  const { handleSubmit, control } = useForm<FormData>();

  const onSubmit = (data: FormData) => {
    alert(JSON.stringify(data));
    // Call your API here.
  };

  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <Controller
        control={control}
        name="username_email"
        render={({
          field: { name, value, onChange, onBlur, ref },
          fieldState: { invalid, error },
        }) => (
          <Input
            ref={ref}
            isRequired
            errorMessage={error?.message}
            isInvalid={invalid}
            label="Username"
            name={name}
            validationBehavior="aria"
            value={value}
            onBlur={onBlur}
            onChange={onChange}
          />
        )}
        rules={{ required: "Name is required." }}
      />
      <Controller
        control={control}
        name="password"
        render={({
          field: { name, value, onChange, onBlur, ref },
          fieldState: { invalid, error },
        }) => (
          <Input
            ref={ref}
            isRequired
            errorMessage={error?.message}
            isInvalid={invalid}
            label="Password"
            name={name}
            type="password"
            validationBehavior="aria"
            value={value}
            onBlur={onBlur}
            onChange={onChange}
          />
        )}
        rules={{ required: "Name is required." }}
      />
      <Button type="submit">Login</Button>
    </Form>
  );
}
