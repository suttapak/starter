export interface Response<T> {
  message: string;
  status: number;
  data: T;
  meta: null;
}
