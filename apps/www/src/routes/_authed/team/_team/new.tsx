import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Form,
  Input,
  Textarea,
} from "@heroui/react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { toastMessage } from "@/utils/toastMessage";
import {
  CreateTeamDto,
  createTeamSchema,
  useCreateTeam,
} from "@/hooks/use-team";
export const Route = createFileRoute("/_authed/team/_team/new")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { mutateAsync, isPending } = useCreateTeam(
    async () =>
      await navigate({
        to: "/team",
      }),
  );

  const { control, handleSubmit } = useForm<CreateTeamDto>({
    resolver: zodResolver(createTeamSchema),
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
          <CardHeader>{t("team.new_team")}</CardHeader>
          <CardBody className="gap-2">
            <Controller
              control={control}
              name="name"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("team.state.name")}
                  name={name}
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
              rules={{ required: "ชื่อแผนก is required." }}
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
                  label={t("team.state.username")}
                  name={name}
                  type="text"
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
              name="description"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Textarea
                  ref={ref}
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  label={t("team.state.description")}
                  name={name}
                  type="text"
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onChange={onChange}
                />
              )}
            />
          </CardBody>
          <CardFooter>
            <Button
              fullWidth
              color="primary"
              isLoading={isPending}
              type="submit"
            >
              {t("team.new_team")}
            </Button>
          </CardFooter>
        </Card>
      </Form>
    </div>
  );
}
