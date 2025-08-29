import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Form,
  Input,
  link,
} from "@heroui/react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { z } from "zod";
import { useTranslation } from "react-i18next";

import { LoginDto, loginSchema, useLogin } from "@/hooks/use-auth";
import { toastMessage } from "@/utils/toastMessage";
import { useAuth } from "@/auth";

const loginRedirect = z.object({
  redirect: z.string().optional(),
});

export const Route = createFileRoute("/_public/login")({
  component: RouteComponent,
  validateSearch: loginRedirect,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { onChangeIsAuthenticated } = useAuth();
  const { redirect } = Route.useSearch();
  const navigate = useNavigate();
  const { mutateAsync, isPending } = useLogin(async () => {
    onChangeIsAuthenticated(true);
    await navigate({
      to: redirect ? `/${redirect}` : "/team",
      replace: true,
    });
  });

  const { control, handleSubmit } = useForm<LoginDto>({
    resolver: zodResolver(loginSchema),
  });

  return (
    <div>
      <Form
        className="flex justify-center items-center"
        onSubmit={handleSubmit((data) =>
          toast.promise(mutateAsync(data), toastMessage),
        )}
      >
        <Card fullWidth className="max-w-sm">
          <CardHeader>{t("login.title")}</CardHeader>
          <CardBody className="gap-2">
            <Controller
              control={control}
              name="username"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("login.username")}
                  name={name}
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "Username is required." }}
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
                  label={t("login.password")}
                  name={name}
                  type="password"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "Password is required." }}
            />
            <Link className={link()} to="/register">
              {t("login.register")}
            </Link>
          </CardBody>
          <CardFooter>
            <Button
              fullWidth
              color="primary"
              isLoading={isPending}
              type="submit"
            >
              {t("login.submit")}
            </Button>
          </CardFooter>
        </Card>
      </Form>
    </div>
  );
}
