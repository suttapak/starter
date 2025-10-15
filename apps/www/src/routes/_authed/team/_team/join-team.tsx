import { Button, Card, CardBody, CardFooter, CardHeader } from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";
import { z } from "zod";
import toast from "react-hot-toast";
import { Box } from "lucide-react";
import { useTranslation } from "react-i18next";

import { toastMessage } from "@/core/utils";
import { useJoinTeam } from "@/features/team";

const validateSearch = z.object({
  token: z.string(),
});

export const Route = createFileRoute("/_authed/team/_team/join-team")({
  component: RouteComponent,
  validateSearch,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { mutateAsync, isPending } = useJoinTeam();

  return (
    <div>
      <Card className="max-w-sm mx-auto">
        <CardHeader className="flex justify-center">
          {t("team.join_with_link")}
        </CardHeader>

        <CardBody>
          <CardBody className="flex flex-col items-center justify-center gap-4">
            <Button isIconOnly size="lg">
              <Box />
            </Button>
            <h3 className="text-lg font-bold">{t("team.join_team")}</h3>
          </CardBody>
        </CardBody>
        <CardFooter>
          <Button
            fullWidth
            color="primary"
            isLoading={isPending}
            onPress={() => toast.promise(() => mutateAsync(), toastMessage)}
          >
            {t("team.join_team")}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
