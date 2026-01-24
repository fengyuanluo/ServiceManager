import { afterEach, describe, expect, test, vi } from "vitest";
import { clearToken, request, setToken } from "../api/client";

describe("api client", () => {
  afterEach(() => {
    clearToken();
    vi.unstubAllGlobals();
  });

  test("adds bearer token", async () => {
    setToken("token-123");
    const payload = { code: 200, message: "success", data: [] };
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      text: async () => JSON.stringify(payload)
    });
    vi.stubGlobal("fetch", fetchMock as unknown as typeof fetch);

    await request("/systemd/list");

    const [, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    const headers = options.headers instanceof Headers
      ? options.headers
      : new Headers(options.headers);
    expect(headers.get("Authorization")).toBe("Bearer token-123");
  });

  test("handles empty response", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 204,
      text: async () => ""
    });
    vi.stubGlobal("fetch", fetchMock as unknown as typeof fetch);

    const result = await request("/systemd/list");

    expect(result).toBeUndefined();
  });
});
