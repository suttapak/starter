import {
  Button,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Skeleton,
  Snippet,
  useDisclosure,
} from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";
import toast from "react-hot-toast";
import React from "react";
import { Share } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useGetTeamById, useShareTeam } from "@/hooks/use-team";
import { toastMessage } from "@/utils/toastMessage";
import TeamInfoCard from "@/components/team-info-card";
import TeamInfoEditModal from "@/components/team-info-edit-modal";

export const Route = createFileRoute("/_authed/team/_id/$id/_layout/team/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const [codeString, setCodeString] = React.useState("");
  const { mutateAsync, isPending } = useShareTeam(setCodeString);
  const { onOpen, ...rest } = useDisclosure();
  const { data, isLoading } = useGetTeamById();

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end items-center gap-2">
        {data?.data.data && <TeamInfoEditModal item={data?.data.data} />}
        <Button color="secondary" endContent={<Share />} onPress={onOpen}>
          {t("team.share")}
        </Button>
      </div>
      {isLoading && (
        <>
          <Skeleton className="h-16 rounded" />
          <Skeleton className="h-16 rounded" />
        </>
      )}
      {data?.data.data && <TeamInfoCard item={data?.data.data} />}

      <Modal {...rest}>
        <ModalContent>
          <ModalHeader>{t("team.share_title")}</ModalHeader>
          <ModalBody>
            <Snippet codeString={codeString} disableCopy={!codeString}>
              {codeString
                ? t("team.share_description_generated")
                : t("team.share_description_not_generated")}
            </Snippet>
          </ModalBody>
          <ModalFooter>
            <Button
              color="primary"
              isLoading={isPending}
              onPress={() => toast.promise(() => mutateAsync(), toastMessage)}
            >
              {t("team.share")}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
}
