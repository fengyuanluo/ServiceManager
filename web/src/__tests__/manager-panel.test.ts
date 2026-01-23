import { fireEvent, render, waitFor } from "@testing-library/vue";
import { describe, expect, test, vi } from "vitest";
import ManagerPanel from "../components/manager/ManagerPanel.vue";
import { ApiError } from "../api/client";
import {
  getServiceDetail,
  getServiceLogs,
  installManager,
  listServices
} from "../api/services";

vi.mock("../api/services", () => ({
  listServices: vi.fn(),
  getServiceDetail: vi.fn(),
  getServiceLogs: vi.fn(),
  runServiceBatch: vi.fn(),
  installManager: vi.fn()
}));

describe("manager panel", () => {
  test("shows detail after selection", async () => {
    const mockList = vi.mocked(listServices);
    const mockDetail = vi.mocked(getServiceDetail);
    const mockLogs = vi.mocked(getServiceLogs);
    mockList.mockResolvedValue([
      {
        name: "sshd",
        status: "running",
        detail: "Secure Shell",
        manager: "systemd"
      }
    ]);
    mockDetail.mockResolvedValue({
      name: "sshd",
      status: "running",
      detail: "Secure Shell",
      manager: "systemd",
      description: "OpenSSH"
    });
    mockLogs.mockResolvedValue(["[info] boot"]);

    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    expect(getByTestId("detail-empty")).toBeTruthy();

    await waitFor(() => {
      expect(mockList).toHaveBeenCalled();
    });

    const target = await waitFor(() => getByTestId("service-item-sshd"));
    await fireEvent.click(target);

    await waitFor(() => {
      expect(getByTestId("detail-name").textContent).toContain("sshd");
    });
    expect(mockDetail).toHaveBeenCalledWith("systemd", "sshd", "");
    expect(mockLogs).toHaveBeenCalledWith(
      "systemd",
      "sshd",
      expect.objectContaining({ lines: 500 })
    );
  });

  test("shows install prompt when manager missing", async () => {
    const mockList = vi.mocked(listServices);
    mockList.mockRejectedValue(new ApiError("manager unavailable", 400, 400));

    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    await waitFor(() => {
      expect(getByTestId("install-card")).toBeTruthy();
    });
  });

  test("installs manager after confirm", async () => {
    const mockList = vi.mocked(listServices);
    const mockInstall = vi.mocked(installManager);
    mockList
      .mockRejectedValueOnce(new ApiError("manager unavailable", 400, 400))
      .mockResolvedValueOnce([]);
    mockInstall.mockResolvedValue({ manager: "systemd", success: true });
    vi.stubGlobal("confirm", vi.fn().mockReturnValue(true));

    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    const button = await waitFor(() => getByTestId("install-action"));
    await fireEvent.click(button);

    await waitFor(() => {
      expect(mockInstall).toHaveBeenCalledWith("systemd");
    });
  });
});
