import {
  Input,
  Textarea,
  Modal,
  ModalContent,
  ModalBody,
  Button,
  useDisclosure,
  ModalHeader,
  ModalFooter,
  Form,
} from "@heroui/react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { TeamResponse } from "@/types/team_response";
import {
  UpdateTeamDto,
  updateTeamSchema,
  useUpdateTeamInfo,
} from "@/hooks/use-team";
import { toastMessage } from "@/utils/toastMessage";

type Props = {
  item: TeamResponse;
};

const TeamInfoEditModal = (props: Props) => {
  const { t } = useTranslation();
  const { item } = props;
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { control, handleSubmit } = useForm<UpdateTeamDto>({
    defaultValues: {
      ...item,
    },
    resolver: zodResolver(updateTeamSchema),
  });

  const { mutateAsync, isPending } = useUpdateTeamInfo(() => {
    onClose();
  });

  return (
    <>
      <Button onPress={onOpen}>{t("team.edit_title")}</Button>
      <Modal {...rest}>
        <Form
          onSubmit={handleSubmit((data) =>
            toast.promise(() => mutateAsync(data), toastMessage),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("team.edit_title")}</ModalHeader>
            <ModalBody>
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
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
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
                    onValueChange={onChange}
                  />
                )}
                rules={{ required: "(Username) ของแผนก is required." }}
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
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("team.state.email")}
                    name={name}
                    type="email"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
              />
              <Controller
                control={control}
                name="phone"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <Input
                    ref={ref}
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("team.state.phone")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
              />
              <Controller
                control={control}
                name="address"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <Textarea
                    ref={ref}
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("team.state.address")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
              />
            </ModalBody>
            <ModalFooter>
              <Button color="primary" isLoading={isPending} type="submit">
                {t("team.edit")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </>
  );
};

export default TeamInfoEditModal;
