import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  cn,
  Form,
  Input,
  link,
} from "@heroui/react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { User } from "lucide-react";
import { useTranslation } from "react-i18next";

import { RegisterDto, registerSchema, useRegister } from "@/hooks/use-auth";
import { toastMessage } from "@/utils/toastMessage";
export const Route = createFileRoute("/_public/register")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const { mutateAsync, isPending } = useRegister(
    async () => await navigate({ to: "/login", replace: true }),
  );

  const { control, handleSubmit } = useForm<RegisterDto>({
    resolver: zodResolver(registerSchema),
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
          <CardHeader>{t("register.title")}</CardHeader>
          <CardBody className="gap-2">
            <Controller
              control={control}
              name="full_name"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("register.state.full_name")}
                  name={name}
                  type="text"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "ชื่อ-นามสกุล is required." }}
            />
            <Controller
              control={control}
              name="email"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("register.state.email")}
                  name={name}
                  type="email"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "อีเมล is required." }}
            />
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
                  label={t("register.state.username")}
                  name={name}
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "ชื่อผู้ใช้ (Username) is required." }}
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
                  label={t("register.state.password")}
                  name={name}
                  type="password"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "รหัสผ่าน is required." }}
            />
            <Controller
              control={control}
              name="confirm_password"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("register.state.confirm_password")}
                  name={name}
                  type="password"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "ยืนยันรหัสผ่าน is required." }}
            />

            <Link className={cn(["pt-2", link()])} to="/login">
              {t("register.login")}
            </Link>
          </CardBody>
          <CardFooter>
            <Button
              fullWidth
              color="primary"
              endContent={<User />}
              isLoading={isPending}
              type="submit"
            >
              {t("register.submit")}
            </Button>
          </CardFooter>
        </Card>
      </Form>
    </div>
  );
}
