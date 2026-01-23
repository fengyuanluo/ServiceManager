import { afterEach, describe, expect, test, vi } from "vitest";
import { clearToken } from "../api/client";
import { createService } from "../api/services";

describe("services api", () => {
  afterEach(() => {
    clearToken();
    vi.unstubAllGlobals();
  });

  test("createService posts payload", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      text: async () => JSON.stringify({ code: 200, message: "ok", data: "ok" })
    });
    vi.stubGlobal("fetch", fetchMock as unknown as typeof fetch);

    const payload = {
      name: "demo",
      description: "Demo",
      exec: "/usr/bin/demo",
      workdir: "/srv/demo",
      user: "root",
      env: "A=B"
    };

    const result = await createService("systemd", payload);

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, options] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain("/api/v1/systemd/create");
    expect(options?.method).toBe("POST");
    expect(options?.body).toBe(JSON.stringify(payload));
    expect(result).toBe("ok");
  });
});
