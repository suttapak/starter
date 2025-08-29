export const errors = (error: any): string => {
  if (typeof error === "string") {
    return error;
  }
  if (error.response) {
    // The request was made and the server responded with a status code
    // that falls out of the range of 2xx
    return error.response.data.message;
  }

  return error.message || "something went wrong";
};

export const toastMessage = {
  loading: "กำลังดำเนินการ",
  success: "ดำเนินการสำเร็จ",
  error: errors,
};
