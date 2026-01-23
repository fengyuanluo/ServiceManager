import { fireEvent, render, waitFor } from "@testing-library/vue";
import { describe, expect, test, vi } from "vitest";
import ManagerPanel from "../components/manager/ManagerPanel.vue";
import { createService, listServices } from "../api/services";

vi.mock("../api/services", () => ({
  listServices: vi.fn(),
  getServiceDetail: vi.fn(),
  getServiceLogs: vi.fn(),
  runServiceBatch: vi.fn(),
  installManager: vi.fn(),
  createService: vi.fn()
}));

describe("service create wizard", () => {
  test("opens wizard from batch toolbar", async () => {
    vi.mocked(listServices).mockResolvedValue([]);

    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    await fireEvent.click(getByTestId("create-service"));

    expect(getByTestId("create-wizard")).toBeTruthy();
  });

  test("submits create request", async () => {
    vi.mocked(listServices).mockResolvedValue([]);
    const mockCreate = vi.mocked(createService);
    mockCreate.mockResolvedValue("ok");

    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    await fireEvent.click(getByTestId("create-service"));
    await fireEvent.update(getByTestId("create-name"), "demo");
    await fireEvent.update(getByTestId("create-exec"), "/usr/bin/demo");
    await fireEvent.click(getByTestId("create-submit"));

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalledWith(
        "systemd",
        expect.objectContaining({ name: "demo" }),
        ""
      );
    });
  });
});
