import { render, fireEvent } from "@testing-library/vue";
import { beforeEach, describe, expect, test } from "vitest";
import App from "../App.vue";

describe("app shell", () => {
  beforeEach(() => {
    localStorage.clear();
    document.documentElement.className = "";
  });

  test("toggles language", async () => {
    const { getByTestId, getByText } = render(App);
    expect(getByText("服务管理器")).toBeTruthy();

    const toggle = getByTestId("lang-toggle");
    await fireEvent.click(toggle);

    expect(getByText("Service Manager")).toBeTruthy();
  });

  test("switches theme class", async () => {
    const { getByTestId } = render(App);
    const shell = getByTestId("app-shell");
    expect(shell.className).toContain("theme-white");

    const darkButton = getByTestId("theme-dark");
    await fireEvent.click(darkButton);
    expect(shell.className).toContain("theme-dark");
  });

  test("opens mobile drawer", async () => {
    const { getByTestId } = render(App);
    const toggle = getByTestId("menu-toggle");
    const drawer = getByTestId("mobile-drawer");

    expect(drawer.getAttribute("data-state")).toBe("closed");
    await fireEvent.click(toggle);
    expect(drawer.getAttribute("data-state")).toBe("open");
  });
});
