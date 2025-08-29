export interface Response<T> {
  message: string;
  status: number;
  data: T;
  meta: any;
}

export interface PaginatedResponse<T> extends Response<T[]> {
  meta: {
    page: number;
    limit: number;
    next_page: number;
    previous_page: number;
    count: number;
    total_page: number;
  };
  data: T[];
}
