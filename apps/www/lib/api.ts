"use client";
import axios from "axios";

export const API_URL = "/api/v1";

export const getJson = async <T = any>(url: string, opt?: any) => {
  return axios.get<T>(API_URL + url, {
    params: {
      ...opt,
    },
  });
};

export const postJson = async <T, K = any>(
  url: string,
  data?: K,
  header = {},
) => {
  return axios.post<T>(API_URL + url, data, header);
};

export const putJson = async <T, R = any>(url: string, data: T) => {
  return axios.put<R>(API_URL + url, data, {});
};

export const deleteJson = async <R = any>(url: string) => {
  return axios.delete<R>(API_URL + url);
};
