import {
  useDisclosure,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  Input,
  ModalFooter,
  Button,
  Form,
} from "@heroui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Upload } from "lucide-react";
import { useForm, Controller } from "react-hook-form";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { FileTrigger } from "./file-trigger";

import { toastMessage } from "@/utils/toastMessage";
import {
  useUploadReportTemplate,
  UploadReportTemplateDto,
  uploadReportTemplateSchema,
} from "@/hooks/use-report";

const UploadReportTemplateModal = () => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { mutateAsync, isPending } = useUploadReportTemplate(onClose);
  const { control, handleSubmit } = useForm<UploadReportTemplateDto>({
    resolver: zodResolver(uploadReportTemplateSchema),
  });

  return (
    <>
      <Button color="secondary" endContent={<Upload />} onPress={onOpen}>
        {t("report.new.title")}
      </Button>
      <Modal {...rest}>
        <Form
          onSubmit={handleSubmit((data) =>
            toast.promise(() => mutateAsync(data), toastMessage),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("report.new.title")}</ModalHeader>
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
                rules={{ required: "Report file is required." }}
              />
            </ModalBody>
            <ModalFooter>
              <Button color="primary" isLoading={isPending} type="submit">
                {t("report.new.action")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </>
  );
};

export default UploadReportTemplateModal;
