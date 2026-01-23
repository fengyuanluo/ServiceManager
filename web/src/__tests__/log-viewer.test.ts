import { fireEvent, render } from "@testing-library/vue";
import { describe, expect, test } from "vitest";
import LogViewer from "../components/manager/LogViewer.vue";

describe("log viewer", () => {
  test("filters by keyword", async () => {
    const { getByTestId, queryAllByTestId } = render(LogViewer, {
      props: {
        lines: ["[info] boot", "[error] failed"],
        lang: "zh",
        active: true
      }
    });

    expect(queryAllByTestId("log-line").length).toBe(2);
    await fireEvent.update(getByTestId("log-search"), "error");
    expect(queryAllByTestId("log-line").length).toBe(1);
  });
});
