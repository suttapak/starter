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
import { t } from "i18next";
import { Fragment } from "react";
import toast from "react-hot-toast";

import { ProductImageResponse } from "@/types/product_response";
import { useFormatImageSrc } from "@/hooks/use-format-image-src";
import { useDeleteProductImage } from "@/hooks/use-product";
import { toastMessage } from "@/utils/toastMessage";

type Props = { file: ProductImageResponse };

function ProductImage({ file }: Props) {
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
          <ModalHeader>ลบรูปภาพ</ModalHeader>
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
              ลบ
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Fragment>
  );
}

export default ProductImage;
