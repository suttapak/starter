import {
  Card,
  CardHeader,
  CardBody,
  useDisclosure,
  Modal,
  ModalContent,
  ModalHeader,
  ModalFooter,
  Button,
  ModalBody,
  Form,
  Input,
} from "@heroui/react";
import { File, Upload } from "lucide-react";
import { useDateFormatter } from "@react-aria/i18n";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { FileTrigger } from "./file-trigger";
import ConfirmDangerModal from "./confirm-danger-modal";

import { ReportResponse } from "@/types/report_response";
import {
  UpdateReportTemplateDto,
  updateReportTemplateSchema,
  useDeleteReportTemplate,
  useUpdateReportTemplate,
} from "@/hooks/use-report";
import { toastMessage } from "@/utils/toastMessage";

type Props = {
  item: ReportResponse;
};

const ReportTemplateCard = ({ item }: Props) => {
  const { t } = useTranslation();
  const formatDate = useDateFormatter({ dateStyle: "full" });
  const { onOpen, onClose, ...rest } = useDisclosure();
  const {
    mutateAsync: mutateAsyncUpdateReportTemplate,
    isPending: isPendingReportTemplate,
  } = useUpdateReportTemplate(onClose);
  const {
    mutateAsync: mutateAsyncDeleteReportTemplate,
    isPending: isPendingDeleteReportTemplate,
  } = useDeleteReportTemplate(onClose);
  const { control, handleSubmit } = useForm<UpdateReportTemplateDto>({
    defaultValues: {
      name: item.name,
    },
    resolver: zodResolver(updateReportTemplateSchema),
  });

  return (
    <>
      <Card isPressable onPress={onOpen}>
        <CardHeader className="gap-2">
          <File />
          {item.display_name}
        </CardHeader>
        <CardBody>
          <p className="text-xs text-default-500">
            {formatDate.format(new Date(item.created_at))}
          </p>
        </CardBody>
      </Card>
      <Modal {...rest}>
        <Form
          onSubmit={handleSubmit((data) =>
            toast.promise(
              () =>
                mutateAsyncUpdateReportTemplate({
                  ...data,
                  templateId: item.id,
                }),
              toastMessage,
            ),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("report.edit.title")}</ModalHeader>
            <ModalBody className="gap-2">
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
                    label={t("report.state.name")}
                    name={name}
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onChange={onChange}
                  />
                )}
                rules={{ required: "Template name is required." }}
              />
              <Controller
                control={control}
                name="file"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <>
                    <FileTrigger
                      ref={ref}
                      buttonProps={{
                        color: "secondary",
                        startContent: <Upload />,
                      }}
                      errorMessage={error?.message}
                      invalid={invalid}
                      name={name}
                      value={value?.name}
                      onBlur={onBlur}
                      onSelect={(file) => {
                        if (!file) return;
                        if (file.length <= 0) return;
                        onChange(file[0]);
                      }}
                    >
                      {t("report.state.file")}
                    </FileTrigger>
                  </>
                )}
              />
            </ModalBody>
            <ModalFooter>
              <ConfirmDangerModal
                isLoading={isPendingDeleteReportTemplate}
                onSubmit={() =>
                  toast.promise(
                    () => mutateAsyncDeleteReportTemplate({ id: item.id }),
                    toastMessage,
                  )
                }
              />
              <Button
                color="primary"
                isLoading={isPendingReportTemplate}
                type="submit"
              >
                {t("report.edit.action")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </>
  );
};

export default ReportTemplateCard;
