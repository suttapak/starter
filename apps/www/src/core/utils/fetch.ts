import axios, { AxiosRequestConfig } from "axios";

export const API_URL = "/api/v1";

export const getJson = async <T = any, K = object>(
  url: string,
  opt?: any,
  body?: K,
) => {
  return axios.get<T>(API_URL + url, {
    params: {
      ...opt,
    },
    data: body,
  });
};

export const postJson = async <T, K = any>(
  url: string,
  data: K,
  config: AxiosRequestConfig = {},
) => {
  return axios.post<T>(API_URL + url, data, config);
};

export const putJson = async <T, R = any>(url: string, data: T) => {
  return axios.put<R>(API_URL + url, data, {});
};

export const deleteJson = async <R = any>(url: string) => {
  return axios.delete<R>(API_URL + url);
};

export const downLoadFile = async (
  uri: string,
  config: AxiosRequestConfig,
  ext = ".xlsx",
) => {
  const response = await axios.get(API_URL + uri, {
    ...config,
    responseType: "blob",
  });
  // Extract filename from Content-Disposition header (optional)
  const disposition = response.headers["content-disposition"];
  let fileName = "example.xlsx";

  if (disposition && disposition.includes("filename=")) {
    fileName = disposition.split("filename=")[1].replace(/"/g, "");
  }
  if (!fileName.endsWith(ext)) {
    fileName += ext; // Ensure the file has the correct extension
  }

  // Create blob URL
  const blob = new Blob([response.data], {
    type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  });
  const url = window.URL.createObjectURL(blob);

  // Create a temporary anchor tag and trigger download
  const a = document.createElement("a");

  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
};
