import { button, link, Snippet } from "@heroui/react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

import { subtitle, title } from "@/components/primitives";
export const Route = createFileRoute("/_public/")({
  component: Index,
});

function Index() {
  const { t } = useTranslation();

  return (
    <section className="flex flex-col items-center justify-center gap-4 py-8 md:py-10">
      <div className="inline-block text-center justify-center">
        <span className={title({ color: "violet" })}>
          {t("landing_page.title_1")}
        </span>
        <br />
        <span className={title()}>{t("landing_page.title_2")}</span>
        <div className={subtitle({ class: "mt-4" })}>
          {t("landing_page.description")}
        </div>
      </div>

      <div className="flex gap-3">
        <Link
          className={button({
            color: "primary",
            radius: "full",
            variant: "shadow",
          })}
          to="/team"
        >
          {t("navbar.use_app")}
        </Link>
      </div>

      <div className="mt-8">
        <Snippet hideCopyButton hideSymbol variant="bordered">
          <span>
            {t("landing_page.snippet_1")}
            <Link className={link({ color: "primary" })} to="/register">
              {t("navbar.register")}
            </Link>
          </span>
        </Snippet>
      </div>
    </section>
  );
}
