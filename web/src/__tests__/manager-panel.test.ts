import { fireEvent, render, waitFor } from "@testing-library/vue";
import { describe, expect, test, vi } from "vitest";
import ManagerPanel from "../components/manager/ManagerPanel.vue";
import { getServiceDetail, getServiceLogs, listServices } from "../api/services";

vi.mock("../api/services", () => ({
  listServices: vi.fn(),
  getServiceDetail: vi.fn(),
  getServiceLogs: vi.fn(),
  runServiceBatch: vi.fn()
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
});
