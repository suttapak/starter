import {
  createFileRoute,
  useNavigate,
  useParams,
} from "@tanstack/react-router";
import {
  Autocomplete,
  AutocompleteItem,
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Form,
  Image,
  Input,
  NumberInput,
  Select,
  SelectItem,
  Textarea,
  Tooltip,
} from "@heroui/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, Controller } from "react-hook-form";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { useRef, useState } from "react";
import { Plus } from "lucide-react";

import {
  CreateProductDto,
  createProductSchema,
  useCreateProduct,
  useUploadProductImages,
} from "@/hooks/use-product";
import { useGetProductCategoryInProductPage } from "@/hooks/use-product-category";
import { toastMessage } from "@/utils/toastMessage";
import { uom } from "@/global/constant/uom";

export const Route = createFileRoute(
  "/_authed/team/_id/$id/_layout/product/new",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { id } = useParams({
    from: "/_authed/team/_id/$id/_layout/product/new",
  });
  const { data, isLoading } = useGetProductCategoryInProductPage();
  const navigate = useNavigate();
  const { mutateAsync: uploadProductImage, isPending: isUploading } =
    useUploadProductImages();
  const { mutateAsync, isPending } = useCreateProduct(async (data) => {
    try {
      if (files && files.length > 0) {
        await toast.promise(
          () => uploadProductImage({ files: files, product_id: data.id }),
          toastMessage,
        );
      }
    } catch {}
    reset();
    navigate({ to: `/team/${id}/product` });
  });
  const { control, handleSubmit, reset } = useForm<CreateProductDto>({
    resolver: zodResolver(createProductSchema),
  });

  const [files, setFiles] = useState<File[] | null>(null);
  const imageUploadInputRef = useRef<HTMLInputElement | null>(null);
  const handleImageUploadChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const files = event.target.files;

    if (!files) return;
    setFiles((p) => (p ? [...p, ...Array.from(files)] : Array.from(files)));
  };
  const handleRemoveImageByIndex = (index: number) => {
    setFiles((prev) => {
      if (!prev) return null;
      const newFiles = [...prev];

      newFiles.splice(index, 1);

      return newFiles;
    });
  };

  return (
    <div className="flex gap-2 flex-col lg:flex-row flex-1 w-full items-start">
      <Form
        className="flex justify-center sticky top-14 flex-1"
        onSubmit={handleSubmit((data) =>
          toast.promise(mutateAsync(data), toastMessage),
        )}
      >
        <Card className="w-full">
          <CardHeader>{t("product.new_product.title")}</CardHeader>
          <CardBody className="gap-2">
            <Controller
              control={control}
              name="code"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Input
                  ref={ref}
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
              rules={{ required: "required." }}
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
                  formatOptions={{
                    style: "currency",
                    currency: "THB",
                  }}
                  isInvalid={invalid}
                  label={t("product.state.price")}
                  name={name}
                  validationBehavior="aria"
                  value={value}
                  onBlur={onBlur}
                  onValueChange={onChange}
                />
              )}
              rules={{ required: "required." }}
            />
            <Controller
              control={control}
              name="uom"
              render={({
                field: { name, value, onChange, onBlur, ref },
                fieldState: { invalid, error },
              }) => (
                <Autocomplete
                  ref={ref}
                  isRequired
                  errorMessage={error?.message}
                  isInvalid={invalid}
                  items={uom}
                  label={t("product.state.oum")}
                  name={name}
                  selectedKey={value ? value.toString() : null}
                  type="text"
                  validationBehavior="aria"
                  onBlur={onBlur}
                  onSelectionChange={(lot) => {
                    if (!lot) return onChange(null);
                    onChange(lot?.toString());
                  }}
                >
                  {(item) => (
                    <AutocompleteItem key={item.label}>
                      {item.label}
                    </AutocompleteItem>
                  )}
                </Autocomplete>
              )}
              rules={{ required: "หน่วยสินค้า is required." }}
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
                      Array.from(keys)?.map((key) => parseInt(key.toString())),
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
          </CardBody>
          <CardFooter>
            <Button
              color="primary"
              isLoading={isPending || isUploading}
              type="submit"
              variant="solid"
            >
              {t("product.new_product.action")}
            </Button>
          </CardFooter>
        </Card>
      </Form>
      <input
        ref={imageUploadInputRef}
        multiple
        accept="image/*"
        id="files"
        name="files"
        style={{ display: "none" }}
        type="file"
        onChange={handleImageUploadChange}
      />
      <Card className="flex-1">
        <CardHeader>{t("product.upload_image.title")}</CardHeader>
        <CardBody>
          <div className="relative -mb-8 columns-1 lg:columns-3 gap-2 *:mb-2">
            {files?.map((file, i) => (
              <Tooltip
                key={`product-image-upload-id-${i}`}
                color="danger"
                content={t("product.upload_image.remove_image")}
              >
                <Image
                  isZoomed
                  src={URL.createObjectURL(file)}
                  onClick={() => handleRemoveImageByIndex(i)}
                />
              </Tooltip>
            ))}
          </div>
        </CardBody>
        <CardFooter>
          <Tooltip content={t("product.upload_image.title")}>
            <Button
              isIconOnly
              size="lg"
              onPress={() => imageUploadInputRef.current?.click()}
            >
              <Plus />
            </Button>
          </Tooltip>
        </CardFooter>
      </Card>
    </div>
  );
}
