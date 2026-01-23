import { fireEvent, render } from "@testing-library/vue";
import { describe, expect, test } from "vitest";
import ManagerPanel from "../components/manager/ManagerPanel.vue";

describe("manager panel", () => {
  test("shows detail after selection", async () => {
    const { getByTestId } = render(ManagerPanel, {
      props: { managerKey: "systemd", lang: "zh" }
    });

    expect(getByTestId("detail-empty")).toBeTruthy();

    const target = getByTestId("service-item-sshd");
    await fireEvent.click(target);
    expect(getByTestId("detail-name").textContent).toContain("sshd");
  });
});
