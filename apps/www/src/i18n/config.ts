import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import en from "./en/translation.json";
import th from "./th/translation.json";

i18next
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    supportedLngs: ["en", "th"],
    fallbackLng: "th", // use en if detected lng is not available
    debug: true,
    resources: {
      en: {
        translation: en,
      },
      th: {
        translation: th,
      },
    },
    // if you see an error like: "Argument of type 'DefaultTFuncReturn' is not assignable to parameter of type xyz"
    // set returnNull to false (and also in the i18next.d.ts options)
    // returnNull: false,
  });
