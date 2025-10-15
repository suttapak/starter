import {
  Button,
  Modal,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/react";
import { useTranslation } from "react-i18next";

type Props = {
  text?: string;
  title?: string;
  isLoading?: boolean;
  onSubmit?: () => void;
};

const ConfirmDangerModal = ({ onSubmit, text, title, isLoading }: Props) => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();

  const handleSubmit = () => {
    onSubmit?.();
    onClose();
  };

  return (
    <>
      <Button color="danger" onPress={onOpen}>
        {text || t("danger_modal.button")}
      </Button>
      <Modal {...rest}>
        <ModalContent>
          <ModalHeader>{title || t("danger_modal.title")}</ModalHeader>
          <ModalFooter>
            <Button onPress={onClose}>{t("danger_modal.cancel")}</Button>
            <Button color="danger" isLoading={isLoading} onPress={handleSubmit}>
              {t("danger_modal.confirm")}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};

export default ConfirmDangerModal;
