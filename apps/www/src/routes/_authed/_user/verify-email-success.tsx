import { Avatar, Button, Card, CardBody, CardFooter } from "@heroui/react";
import { Link } from "@tanstack/react-router";
import { createFileRoute } from "@tanstack/react-router";
import { ShieldCheck, User } from "lucide-react";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_authed/_user/verify-email-success")({
  component: RouteComponent,
});

function RouteComponent() {
  const { t } = useTranslation();

  // create a static page that shows a success message for email verification
  return (
    <Card className="max-w-sm mx-auto mt-10">
      <CardBody className="flex flex-col items-center justify-center gap-4 text-center my-10">
        <Avatar
          className="h-48 w-48"
          color="success"
          icon={<ShieldCheck size={140} />}
        />
        <p>{t("verify_email_success.message")}</p>
      </CardBody>
      <CardFooter>
        <Button
          fullWidth
          as={Link}
          color="primary"
          endContent={<User />}
          to="/profile"
        >
          {t("verify_email_success.go_to_profile")}
        </Button>
      </CardFooter>
    </Card>
  );
}
