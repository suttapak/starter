import { createFileRoute, Link } from "@tanstack/react-router";
import { Avatar, Button, Card, CardBody, CardFooter } from "@heroui/react";
import { Link as TanstackLink } from "@tanstack/react-router";
import { Box, FilePlus, Search } from "lucide-react";
import { useDateFormatter } from "@react-aria/i18n";
import { useTranslation } from "react-i18next";

import { useGetTeamMe } from "@/hooks/use-team";
import MSplashPage from "@/components/splash";

export const Route = createFileRoute("/_authed/team/_team/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();
  const { isLoading, data } = useGetTeamMe();
  const teams = data?.data.data || [];

  const formatDate = useDateFormatter({
    dateStyle: "medium",
  });

  if (isLoading) return <MSplashPage />;

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end gap-2">
        <Button
          as={Link}
          color="secondary"
          endContent={<Search />}
          to="/team/search"
        >
          {t("team.search_join")}
        </Button>
        <Button
          as={Link}
          color="primary"
          endContent={<FilePlus />}
          to="/team/new"
        >
          {t("team.new_team")}
        </Button>
      </div>
      <div className="grid grid-cols-1 gap-2 sm:grid-cols-3">
        {teams.map((team) => (
          <Card
            key={team.id}
            isPressable
            as={TanstackLink}
            to={`/team/${team.id}`}
          >
            <CardBody>
              <div className="flex gap-2">
                <Avatar className="h-14 w-14" icon={<Box />} radius="md" />
                <div className="flex flex-col gap-1">
                  <h3 className="text-lg font-semibold">{team.name}</h3>
                  <h3 className="text-default-400">@{team.username}</h3>
                  <p className="text-gray-600">{team.description}</p>
                </div>
              </div>
            </CardBody>
            <CardFooter>
              <p className="text-sm text-default-400">
                {formatDate.format(new Date(team.created_at))}
              </p>
            </CardFooter>
          </Card>
        ))}
        {teams.length === 0 && (
          <p className="text-center text-gray-500">
            {t("not_found.global.title")}
          </p>
        )}
      </div>
    </div>
  );
}
