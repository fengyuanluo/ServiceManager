import { render } from "@testing-library/vue";
import { describe, expect, test } from "vitest";
import App from "../App.vue";
import { clearToken } from "../api/client";

describe("app auth", () => {
  test("shows login when no token", () => {
    localStorage.clear();
    clearToken();
    const { getByTestId } = render(App);
    expect(getByTestId("login-panel")).toBeTruthy();
  });
});
