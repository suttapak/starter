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
import { MoreVertical } from "lucide-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { ProductCategoryResponse } from "@/types/product_category_response";
import {
  UpdateProductCategoryDto,
  updateProductCategorySchema,
  useUpdateProductCategory,
} from "@/hooks/use-product-category";
import { toastMessage } from "@/utils/toastMessage";

type CategoryEditModalProps = {
  item: ProductCategoryResponse;
};

const CategoryEditModal = ({ item }: CategoryEditModalProps) => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { mutateAsync, isPending } = useUpdateProductCategory(item.id, () => {
    onClose();
    reset();
  });
  const { control, handleSubmit, reset } = useForm<UpdateProductCategoryDto>({
    resolver: zodResolver(updateProductCategorySchema),
    defaultValues: {
      name: item.name,
    },
  });

  return (
    <>
      <Button
        isIconOnly
        color="primary"
        size="sm"
        variant="flat"
        onPress={onOpen}
      >
        <MoreVertical size={18} />
      </Button>
      <Modal {...rest}>
        <Form
          className="flex justify-center items-center"
          onSubmit={handleSubmit((data) =>
            toast.promise(mutateAsync(data), toastMessage),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("category.edit.title")}</ModalHeader>
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
                rules={{ required: "required." }}
              />
            </ModalBody>
            <ModalFooter>
              <Button
                color="primary"
                isLoading={isPending}
                type="submit"
                variant="solid"
              >
                {t("category.edit.action")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </>
  );
};

export default CategoryEditModal;
