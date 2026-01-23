import { computed, ref } from "vue";

type ApiResponse<T> = {
  code: number;
  message: string;
  data?: T;
};

export class ApiError extends Error {
  status: number;
  code: number;

  constructor(message: string, status: number, code: number) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

const storageKey = "sm_token";
const apiBase = import.meta.env.VITE_API_BASE ?? "/api/v1";

const tokenRef = ref<string>(localStorage.getItem(storageKey) ?? "");
const isAuthenticated = computed(() => tokenRef.value.length > 0);

export const getToken = () => tokenRef.value;

export const setToken = (token: string) => {
  tokenRef.value = token;
  localStorage.setItem(storageKey, token);
};

export const clearToken = () => {
  tokenRef.value = "";
  localStorage.removeItem(storageKey);
};

export const useAuthState = () => ({
  token: tokenRef,
  isAuthenticated
});

type RequestOptions = RequestInit & { auth?: boolean };

export async function request<T>(path: string, options: RequestOptions = {}) {
  const { auth = true, ...fetchOptions } = options;
  const headers = new Headers(fetchOptions.headers ?? {});

  if (!headers.has("Content-Type") && fetchOptions.body) {
    headers.set("Content-Type", "application/json");
  }
  if (auth) {
    const token = getToken();
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
  }

  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  const response = await fetch(`${apiBase}${normalizedPath}`, {
    ...fetchOptions,
    headers
  });

  const rawText = await response.text();
  let payload: ApiResponse<T> | null = null;
  if (rawText) {
    try {
      payload = JSON.parse(rawText) as ApiResponse<T>;
    } catch {
      payload = null;
    }
  }

  if (!response.ok || payload?.code !== 200) {
    if (response.status === 401) {
      clearToken();
    }
    throw new ApiError(payload?.message ?? response.statusText, response.status, payload?.code ?? response.status);
  }

  return payload?.data as T;
}
