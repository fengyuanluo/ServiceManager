import { fireEvent, render } from "@testing-library/vue";
import { describe, expect, test } from "vitest";
import ServiceEditor from "../components/manager/ServiceEditor.vue";
import type { ServiceDetail } from "../types/service";

const sampleService: ServiceDetail = {
  id: "sshd",
  name: "sshd",
  status: "running",
  summary: "Secure Shell",
  manager: "systemd",
  description: "OpenSSH",
  unitPath: "/lib/systemd/system/ssh.service"
};

describe("service editor", () => {
  test("toggles raw mode", async () => {
    const { getByTestId } = render(ServiceEditor, {
      props: {
        service: sampleService,
        managerKey: "systemd",
        lang: "zh"
      }
    });

    expect(getByTestId("editor-form")).toBeTruthy();
    await fireEvent.click(getByTestId("toggle-raw"));
    expect(getByTestId("editor-raw")).toBeTruthy();
  });
});
