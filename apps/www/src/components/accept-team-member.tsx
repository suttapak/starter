import React from "react";
import {
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Tooltip,
  useDisclosure,
  User,
} from "@heroui/react";
import { MoreVertical } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

import { TeamMemberResponse } from "@/types/team_member_response";
import { useAcceptTeamMember } from "@/hooks/use-team";
import { toastMessage } from "@/utils/toastMessage";

type Props = {
  team: TeamMemberResponse;
};

const AcceptTeamMember = ({ team }: Props) => {
  const { t } = useTranslation();
  const { onOpen, onClose, ...rest } = useDisclosure();
  const { mutateAsync, isPending } = useAcceptTeamMember(onClose);

  return (
    <React.Fragment>
      <Tooltip content={t("team.approve_title")}>
        <Button isIconOnly color="secondary" size="sm" onPress={onOpen}>
          <MoreVertical size={18} />
        </Button>
      </Tooltip>
      <Modal {...rest}>
        <ModalContent>
          <ModalHeader>{t("team.approve_title")}</ModalHeader>
          <ModalBody className="flex justify-start">
            <User
              avatarProps={{ name: team.user.full_name }}
              description={team.user.email}
              name={team.user.full_name}
            />
          </ModalBody>
          <ModalFooter>
            <Button
              color="success"
              isLoading={isPending}
              onPress={() =>
                toast.promise(
                  () =>
                    mutateAsync({
                      user_id: team.user_id,
                      role_id: team.team_role_id,
                    }),
                  toastMessage,
                )
              }
            >
              {t("team.approve")}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </React.Fragment>
  );
};

export default AcceptTeamMember;
