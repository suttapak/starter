import {
  Button,
  Form,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { PlusCircle } from "lucide-react";
import { useTranslation } from "react-i18next";

import {
  CreateProductCategoryDto,
  createProductCategorySchema,
  useCreateProductCategory,
} from "@/hooks/use-product-category";
import { toastMessage } from "@/utils/toastMessage";

const CategoryCreateModal = () => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { mutateAsync, isPending } = useCreateProductCategory(() => {
    onClose();
    reset();
  });
  const { control, handleSubmit, reset } = useForm<CreateProductCategoryDto>({
    resolver: zodResolver(createProductCategorySchema),
  });

  return (
    <>
      <Button color="primary" endContent={<PlusCircle />} onPress={onOpen}>
        {t("category.new.action")}
      </Button>
      <Modal {...rest}>
        <Form
          className="flex justify-center items-center"
          onSubmit={handleSubmit((data) =>
            toast.promise(mutateAsync(data), toastMessage),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("category.new.title")}</ModalHeader>
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
                    label={t("category.state.name")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onChange={onChange}
                  />
                )}
                rules={{ required: "ชื่อหมวดหมู่ is required." }}
              />
            </ModalBody>
            <ModalFooter>
              <Button
                color="primary"
                isLoading={isPending}
                type="submit"
                variant="solid"
              >
                {t("category.new.action")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </>
  );
};

export default CategoryCreateModal;
