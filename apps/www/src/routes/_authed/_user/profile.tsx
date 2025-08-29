import {
  Avatar,
  Button,
  Card,
  CardBody,
  CardHeader,
  Chip,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Select,
  Selection,
  SelectItem,
} from "@heroui/react";
import { createFileRoute } from "@tanstack/react-router";
import { Send } from "lucide-react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { useCallback, useRef, useState } from "react";

import MSplashPage from "@/components/splash";
import {
  useGetImageProfile,
  useGetUserMe,
  useUploadProfileImage,
} from "@/hooks/use-user";
import { useSendVerifyEmail } from "@/hooks/use-auth";
import { toastMessage } from "@/utils/toastMessage";
import { useLocale } from "@/hooks/use-locale";

export const Route = createFileRoute("/_authed/_user/profile")({
  component: RouteComponent,
});
export const langs = [
  { key: "th", label: "ไทย" },
  { key: "en", label: "English" },
];

function RouteComponent() {
  const { t, i18n } = useTranslation();
  const { data, isLoading, error } = useGetUserMe();
  const user = data?.data.data;

  const { onChangeLocale } = useLocale();

  const { mutateAsync, isPending } = useSendVerifyEmail();
  const handleChangeLang = useCallback((lng: "th" | "en") => {
    i18n.changeLanguage(lng);
    if (lng === "th") {
      onChangeLocale("th-TH");
    }
    if (lng === "en") {
      onChangeLocale("en-US");
    }
  }, []);
  const { mutateAsync: uploadProfileImage, isPending: isUploading } =
    useUploadProfileImage(() => setFile(null));
  const imageProfileRef = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | null>(null);

  const handleChangeFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;

    if (!files) return;
    // check if files length is 0
    if (files.length <= 0) return;
    setFile(files[0]);
  };

  const handleClearFile = () => {
    setFile(null);
  };

  const handleClickUploadProfileImage = () => {
    if (!imageProfileRef.current) return;
    imageProfileRef.current.click();
  };

  const { src } = useGetImageProfile();

  if (isLoading) return <MSplashPage />;

  if (error) {
    throw error;
  }

  return (
    <div className="flex flex-col gap-2">
      <Card>
        <CardBody className="gap-2">
          <div className="flex justify-end">
            <Chip
              color={user?.email_verifyed ? "success" : "default"}
              radius="sm"
            >
              {user?.email_verifyed
                ? t("user.verify_email_status.verified")
                : t("user.verify_email_status.not_verified")}
            </Chip>
          </div>
          <div className="flex justify-center items-center">
            <Avatar
              className="w-28 h-28 text-large"
              name={user?.full_name}
              src={src}
            />
          </div>
          <div className="flex justify-center">
            {!user?.email_verifyed && (
              <Button
                color="secondary"
                endContent={<Send />}
                isLoading={isPending}
                onPress={() => toast.promise(() => mutateAsync(), toastMessage)}
              >
                {t("user.verify_email.action")}
              </Button>
            )}
          </div>
          <Modal
            isOpen={!!file}
            onClose={handleClearFile}
            onOpenChange={handleClearFile}
          >
            <ModalContent>
              <ModalHeader>{t("user.upload_image_profile.title")}</ModalHeader>
              <ModalBody className="items-center justify-center">
                {file && (
                  <Avatar
                    className="w-40 h-40 text-large"
                    src={URL.createObjectURL(file)}
                  />
                )}
              </ModalBody>
              <ModalFooter>
                <Button onPress={handleClearFile}>
                  {t("user.upload_image_profile.cancel")}
                </Button>
                <Button
                  color="primary"
                  endContent={<Send />}
                  isLoading={isUploading}
                  onPress={() =>
                    file &&
                    toast.promise(
                      () => uploadProfileImage(file as File),
                      toastMessage,
                    )
                  }
                >
                  {t("user.upload_image_profile.action")}
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
          <input
            ref={imageProfileRef}
            accept="image/*"
            multiple={false}
            style={{ display: "none" }}
            type="file"
            onChange={handleChangeFile}
          />
          <Button
            color="secondary"
            variant="ghost"
            onPress={handleClickUploadProfileImage}
          >
            {t("user.upload_image_profile.title")}
          </Button>
        </CardBody>
      </Card>
      <Card>
        <CardBody className="gap-2">
          <Input
            isReadOnly
            label={t("user.state.username")}
            labelPlacement="outside"
            value={user?.username}
          />
          <Input
            isReadOnly
            label={t("user.state.full_name")}
            labelPlacement="outside"
            value={user?.full_name}
          />
          <Input
            isReadOnly
            label={t("user.state.email")}
            labelPlacement="outside"
            value={user?.email}
          />
          <Input
            isReadOnly
            label={t("user.state.role")}
            labelPlacement="outside"
            value={user?.role?.name}
          />
        </CardBody>
      </Card>
      <Card>
        <CardHeader>{t("user.group_title.change_lang")}</CardHeader>
        <CardBody className="gap-2">
          <Select
            items={langs}
            label={t("user.change_lang.state.lang.label")}
            placeholder={t("user.change_lang.state.lang.placeholder")}
            selectedKeys={[i18n.language]}
            onSelectionChange={(key: Selection) => {
              if (key === "all") return;
              const lngs = Array.from(key);

              if (lngs.length === 0) return;
              handleChangeLang(lngs[0] as "th" | "en");
            }}
          >
            {(lng) => <SelectItem>{lng.label}</SelectItem>}
          </Select>
        </CardBody>
      </Card>
    </div>
  );
}
