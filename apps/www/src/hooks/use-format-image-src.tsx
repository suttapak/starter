const BASE_URL = "/api/v1/";

export const useFormatImageSrc = () => {
  const format = (src: string = "") => {
    if (src === "") return undefined;
    if (src.startsWith("/")) {
      src = src.slice(1);
    }

    return `${BASE_URL}${src}`;
  };

  return { format };
};
