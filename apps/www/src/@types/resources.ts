import th from "../i18n/th/translation.json";
import en from "../i18n/en/translation.json";

const resources = {
  translation: { ...th, en },
} as const;

export default resources;
