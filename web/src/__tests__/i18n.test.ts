import { describe, expect, test } from "vitest";
import { t } from "../i18n";

describe("i18n", () => {
  test("lookup", () => {
    expect(t("nav.systemd", "zh")).toBe("Systemd");
  });
});
