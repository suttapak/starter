import { createFileRoute, notFound } from "@tanstack/react-router";
import {
  Button,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
  Image,
  Input,
  NumberInput,
  Select,
  SelectItem,
  Textarea,
  Tooltip,
} from "@heroui/react";
import { useTranslation } from "react-i18next";
import z from "zod";
import { useRef, useState } from "react";
import { Plus, Upload } from "lucide-react";
import toast from "react-hot-toast";

import { useGetProduct, useUploadProductImages } from "@/hooks/use-product";
import MSplashPage from "@/components/splash";
import ProductActionModal from "@/components/product-action-modal";
import { toastMessage } from "@/utils/toastMessage";
import ProductImage from "@/components/product-image";
const validateSearch = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
  transactionPage: z.number().prefault(1),
  transactionLimit: z.number().prefault(10),
});

export const Route = createFileRoute(
  "/_authed/team/_id/$id/_layout/product/$pid/",
)({
  validateSearch: validateSearch,
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { pid } = Route.useParams();
  const { data, isLoading, isError, error } = useGetProduct(pid);

  const { mutateAsync: uploadProductImage, isPending: isUploading } =
    useUploadProductImages();

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

  const handleFileUpload = async () => {
    if (!files) return;
    try {
      const productId = parseInt(pid);

      if (isNaN(productId)) throw new Error("Invalid product ID");
      await uploadProductImage({ product_id: productId, files: files });
      setFiles(null);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error(error);
    } finally {
      setFiles(null);
    }
  };

  if (isLoading) return <MSplashPage />;

  if (data?.data.status === 404) throw notFound();
  if (!data?.data.data) throw error;
  if (isError) throw error;

  return (
    <div className="flex flex-col gap-2">
      <div className="flex gap-2 flex-col lg:flex-row  w-full items-start">
        <Card fullWidth className="flex-1">
          <CardHeader>{t("product.product_info_title")}</CardHeader>
          <CardBody className="gap-2">
            <Input
              isReadOnly
              label={t("product.state.code")}
              value={data?.data.data.code}
            />
            <Input
              isReadOnly
              label={t("product.state.name")}
              value={data?.data.data.name}
            />
            <Textarea
              isReadOnly
              label={t("product.state.description")}
              value={data?.data.data.description}
            />
            <NumberInput
              isReadOnly
              formatOptions={{
                style: "currency",
                currency: "THB",
                currencyDisplay: "symbol",
              }}
              label={t("product.state.price")}
              value={data?.data.data.price * 0.01}
            />
            <Input
              isReadOnly
              label={t("product.state.oum")}
              value={data?.data.data.uom}
            />
            <Select
              isRequired
              isLoading={isLoading}
              items={data?.data.data.product_product_category || []}
              label={t("product.state.category")}
              selectedKeys={data?.data.data.product_product_category?.map(
                (id) => id.category.id.toString(),
              )}
              selectionMode="multiple"
              validationBehavior="aria"
            >
              {(item) => (
                <SelectItem key={item.category.id.toString()}>
                  {item.category.name}
                </SelectItem>
              )}
            </Select>
          </CardBody>
          <CardFooter>
            <ProductActionModal isText item={data?.data.data} />
          </CardFooter>
        </Card>
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
          <CardHeader>{t("product.state.images")}</CardHeader>
          <CardBody>
            <div className="relative -mb-8 columns-1 lg:columns-3 gap-2 *:mb-2">
              {data.data.data.product_image?.map((file, i) => (
                <ProductImage key={`product-image-${i}`} file={file} />
              ))}
            </div>

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
            <Button
              isIconOnly
              isLoading={isUploading}
              size="lg"
              onPress={() => imageUploadInputRef.current?.click()}
            >
              <Plus />
            </Button>
            {files && (
              <Button
                color="primary"
                isLoading={isUploading}
                size="lg"
                onPress={() =>
                  toast.promise(() => handleFileUpload(), toastMessage)
                }
              >
                <Upload />
              </Button>
            )}
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
