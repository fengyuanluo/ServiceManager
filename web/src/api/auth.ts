import { clearToken, request, setToken } from "./client";

type LoginResponse = {
  token: string;
  expire_at: string;
};

export async function login(username: string, password: string, expire = 0) {
  const data = await request<LoginResponse>("/user/login", {
    method: "POST",
    auth: false,
    body: JSON.stringify({ username, password, expire })
  });
  setToken(data.token);
  return data;
}

export const logout = () => {
  clearToken();
};
