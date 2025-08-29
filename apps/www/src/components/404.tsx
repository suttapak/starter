import { Card, CardHeader, Button, CardBody, CardFooter } from "@heroui/react";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

export default function NotFound() {
  const { t } = useTranslation();

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <Card className="max-w-md w-full">
        <CardHeader>{t("not_found.global.title")}</CardHeader>
        <CardBody>
          <p className="mb-4 text-gray-600">
            {t("not_found.global.description")}
          </p>
        </CardBody>
        <CardFooter>
          <Button>
            <Link to="/team">{t("not_found.global.action")}</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
