import {
  Button,
  Form,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  NumberInput,
  Select,
  SelectItem,
  Textarea,
  useDisclosure,
} from "@heroui/react";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import toast from "react-hot-toast";
import { MoreVertical } from "lucide-react";
import React from "react";
import { useTranslation } from "react-i18next";

import { toastMessage } from "@/utils/toastMessage";
import {
  createUpdateSchema,
  UpdateProductDto,
  useUpdateProduct,
} from "@/hooks/use-product";
import { useGetProductCategoryInProductPage } from "@/hooks/use-product-category";
import { ProductResponse } from "@/types/product_response";
type ProductActionModalProps = {
  item: ProductResponse;
  isText?: boolean;
};
const ProductActionModal = ({ item, isText }: ProductActionModalProps) => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { data, isLoading } = useGetProductCategoryInProductPage();
  const { mutateAsync, isPending } = useUpdateProduct(item.id, () => {
    onClose();
    reset();
  });
  const { control, handleSubmit, reset } = useForm<UpdateProductDto>({
    resolver: zodResolver(createUpdateSchema),
    defaultValues: {
      code: item.code,
      name: item.name,
      description: item.description || "",
      uom: item.uom,
      price: item.price * 0.01 || 0,
      category_id: item.product_product_category.map((c) => c.category.id),
    },
  });

  return (
    <React.Fragment>
      {isText ? (
        <Button color="secondary" size="sm" variant="bordered" onPress={onOpen}>
          {t("product.edit_product.action")}
        </Button>
      ) : (
        <Button
          isIconOnly
          color="secondary"
          size="sm"
          variant="bordered"
          onPress={onOpen}
        >
          <MoreVertical size={18} />
        </Button>
      )}
      <Modal {...rest}>
        <Form
          className="flex justify-center items-center"
          onSubmit={handleSubmit((data) =>
            toast.promise(mutateAsync(data), toastMessage),
          )}
        >
          <ModalContent>
            <ModalHeader>{t("product.edit_product.title")}</ModalHeader>
            <ModalBody>
              <Controller
                control={control}
                name="code"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <Input
                    ref={ref}
                    isRequired
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("product.state.code")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
                rules={{ required: t("product.state.validation.code") }}
              />
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
                    label={t("product.state.name")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
                rules={{ required: t("product.state.validation.name") }}
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
                    isRequired
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("product.state.description")}
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
                name="price"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <NumberInput
                    ref={ref}
                    isRequired
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("product.state.price")}
                    name={name}
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onValueChange={onChange}
                  />
                )}
                rules={{ required: t("product.state.validation.price") }}
              />
              <Controller
                control={control}
                name="uom"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <Input
                    ref={ref}
                    isRequired
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    label={t("product.state.oum")}
                    name={name}
                    type="text"
                    validationBehavior="aria"
                    value={value}
                    onBlur={onBlur}
                    onChange={onChange}
                  />
                )}
                rules={{ required: t("product.state.validation.oum") }}
              />
              <Controller
                control={control}
                name="category_id"
                render={({
                  field: { name, value, onChange, onBlur, ref },
                  fieldState: { invalid, error },
                }) => (
                  <Select
                    ref={ref}
                    isRequired
                    errorMessage={error?.message}
                    isInvalid={invalid}
                    isLoading={isLoading}
                    items={data?.data.data || []}
                    label={t("product.state.category")}
                    name={name}
                    selectedKeys={value?.map((v) => v.toString())}
                    selectionMode="multiple"
                    validationBehavior="aria"
                    onBlur={onBlur}
                    onSelectionChange={(keys) => {
                      onChange(
                        Array.from(keys)?.map((key) =>
                          parseInt(key.toString()),
                        ),
                      );
                    }}
                  >
                    {(item) => (
                      <SelectItem key={item.id.toString()}>
                        {item.name}
                      </SelectItem>
                    )}
                  </Select>
                )}
              />
            </ModalBody>
            <ModalFooter>
              <Button
                color="primary"
                isLoading={isPending}
                type="submit"
                variant="solid"
              >
                {t("product.edit_product.action")}
              </Button>
            </ModalFooter>
          </ModalContent>
        </Form>
      </Modal>
    </React.Fragment>
  );
};

export default ProductActionModal;
