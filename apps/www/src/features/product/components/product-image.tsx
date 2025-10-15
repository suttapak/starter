import type { ProductImageResponse } from "../types/product";

import {
  Button,
  Image,
  Modal,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Tooltip,
  useDisclosure,
} from "@heroui/react";
import { Fragment } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { useDeleteProductImage } from "../api/use-product";

import { useFormatImageSrc } from "@/core/hooks/use-format-image-src";
import { toastMessage } from "@/core/utils/toastMessage";

type Props = { file: ProductImageResponse };

function ProductImage({ file }: Props) {
  const { t } = useTranslation();
  const { format } = useFormatImageSrc();
  const { onOpen, onClose, ...rest } = useDisclosure();

  const { mutateAsync, isPending } = useDeleteProductImage(onClose);

  return (
    <Fragment>
      <Tooltip color="danger" content={t("product.upload_image.remove_image")}>
        <Image isZoomed src={format(file.image.path)} onClick={onOpen} />
      </Tooltip>
      <Modal {...rest}>
        <ModalContent>
          <ModalHeader>{t("product.upload_image.remove_image")}</ModalHeader>
          <ModalFooter>
            <Button
              color="danger"
              isLoading={isPending}
              onPress={() =>
                toast.promise(
                  () =>
                    mutateAsync({
                      productId: file.product_id,
                      productImageId: file.id,
                    }),
                  toastMessage,
                )
              }
            >
              {t("product.upload_image.remove_image")}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Fragment>
  );
}

export default ProductImage;
