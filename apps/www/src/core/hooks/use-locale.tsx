import axios from "axios";
import { useCallback, useEffect, useState } from "react";

export const useLocale = () => {
  // get localstorage locale th-TH en-US
  const [locale, setLocale] = useState(() => {
    if (typeof window !== "undefined") {
      return localStorage.getItem("locale") || "th-TH";
    }

    return "th-TH"; // default locale
  });

  useEffect(() => {
    const init = () => {
      if (typeof window === "undefined") {
        return "th-TH";
      }

      return (localStorage.getItem("locale") as "th-TH" | "en-US") || "th-TH";
    };
    const fn = () => {
      onChangeLocale(init());
    };

    fn();
  }, []);

  const getLocaleShortCode = (locale: "th-TH" | "en-US"): "th" | "en" => {
    return locale === "en-US" ? "en" : "th";
  };

  const onChangeLocale = useCallback((newLocale: "th-TH" | "en-US") => {
    axios.defaults.headers.common["lng"] = getLocaleShortCode(newLocale);
    setLocale(newLocale);
    if (typeof window !== "undefined") {
      localStorage.setItem("locale", newLocale);
    }
  }, []);

  return { locale, onChangeLocale };
};
