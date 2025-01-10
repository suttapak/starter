import axios from "axios";
import FileDownload from "js-file-download";

import { Response } from "@/types/response";

export const API_URL = "/api/v1";

export const getJson = async <T = any, K = object>(
  url: string,
  opt?: any,
  body?: K,
) => {
  return axios.get<Response<T>>(API_URL + url, {
    params: {
      ...opt,
    },
    data: body,
  });
};

export const postJson = async <T, K = any>(
  url: string,
  data: K,
  header = {},
) => {
  return axios.post<Response<T>>(API_URL + url, data, header);
};

export const putJson = async <T, R = any>(url: string, data: T) => {
  return axios.put<R>(API_URL + url, data, {});
};

export const deleteJson = async <R = any,>(url: string) => {
  return axios.delete<R>(API_URL + url);
};

export const download = async (url: string, params = {}) => {
  const res = await axios.get(API_URL + url, {
    responseType: "blob",
    params: {
      ...params,
    },
  });
  const disposition = res.headers["content-disposition"];
  let filename = "default.xlsx"; // Fallback filename if none is provided by the server

  if (disposition && disposition.includes("filename")) {
    const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
    const matches = filenameRegex.exec(disposition);

    if (matches != null && matches[1]) {
      filename = matches[1].replace(/['"]/g, ""); // Remove quotes if present
    }
  }

  FileDownload(res.data, filename);
};
