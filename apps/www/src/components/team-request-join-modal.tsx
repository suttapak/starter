import {
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@heroui/react";
import { MoreVertical, Plus } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { TeamResponse } from "@/types/team_response";
import { useRequestJoinTeam } from "@/hooks/use-team";
import { toastMessage } from "@/utils/toastMessage";

type Props = {
  team: TeamResponse;
};

const TeamRequestJoinModal = ({ team }: Props) => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { mutateAsync, isPending } = useRequestJoinTeam(onClose);

  return (
    <div>
      <Button
        isIconOnly
        color="primary"
        size="sm"
        variant="bordered"
        onPress={onOpen}
      >
        <MoreVertical size={18} />
      </Button>
      <Modal {...rest}>
        <ModalContent>
          <ModalHeader>
            {t("team.request_join")} [{team.name}]
          </ModalHeader>
          <ModalBody className="flex justify-center items-center">
            <Button isIconOnly color="success" size="lg">
              <Plus />
            </Button>
            <h3 className="text-lg font-bold">{t("team.request_join")}</h3>
          </ModalBody>
          <ModalFooter>
            <Button
              color="primary"
              isLoading={isPending}
              onPress={() =>
                toast.promise(() => mutateAsync(team.id), toastMessage)
              }
            >
              {t("team.request_join")}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
};

export default TeamRequestJoinModal;
