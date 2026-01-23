import { afterEach, describe, expect, test, vi } from "vitest";
import { login, logout } from "../api/auth";
import { clearToken, getToken } from "../api/client";

describe("auth", () => {
  afterEach(() => {
    clearToken();
    localStorage.clear();
    vi.unstubAllGlobals();
  });

  test("login stores token", async () => {
    const payload = {
      code: 200,
      message: "success",
      data: { token: "token-abc", expire_at: "2099-01-01" }
    };
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      text: async () => JSON.stringify(payload)
    });
    vi.stubGlobal("fetch", fetchMock as unknown as typeof fetch);

    await login("admin", "pass");

    expect(getToken()).toBe("token-abc");
    logout();
    expect(getToken()).toBe("");
  });
});
