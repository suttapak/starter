import { Card, CardHeader, CardBody, Input, Textarea } from "@heroui/react";
import { useTranslation } from "react-i18next";

import { TeamResponse } from "@/types/team_response";

type Props = {
  item: TeamResponse;
};

const TeamInfoCard = (props: Props) => {
  const { t } = useTranslation();
  const { item } = props;

  return (
    <>
      <Card>
        <CardHeader>{t("team.info_title")}</CardHeader>
        <CardBody className="gap-2">
          <Input isReadOnly label={t("team.state.name")} value={item.name} />
          <Input
            isReadOnly
            label={t("team.state.username")}
            value={item.username}
          />
          <Textarea
            isReadOnly
            label={t("team.state.description")}
            value={item.description}
          />
        </CardBody>
      </Card>
      <Card>
        <CardHeader>{t("team.contact_title")}</CardHeader>
        <CardBody className="gap-2">
          <Input isReadOnly label={t("team.state.email")} value={item.email} />
          <Input isReadOnly label={t("team.state.phone")} value={item.phone} />
          <Textarea
            isReadOnly
            label={t("team.state.address")}
            value={item.address}
          />
        </CardBody>
      </Card>
    </>
  );
};

export default TeamInfoCard;
