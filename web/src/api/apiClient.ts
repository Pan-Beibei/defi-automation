import axios, {
  type AxiosRequestConfig,
  type AxiosResponse,
  type AxiosError,
} from "axios";

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

interface Result {
  code: number;
}

class APIClient {
  get<T = unknown>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: "GET" });
  }

  post<T = unknown>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: "POST" });
  }

  request<T = unknown>(config: AxiosRequestConfig): Promise<T> {
    return new Promise((resolve, reject) => {
      axiosInstance
        .request<unknown, AxiosResponse<Result>>(config)
        .then((res: AxiosResponse<Result>) => {
          resolve(res.data as T);
        })
        .catch((e: Error | AxiosError) => {
          reject(e);
        });
    });
  }
}

export default APIClient;
