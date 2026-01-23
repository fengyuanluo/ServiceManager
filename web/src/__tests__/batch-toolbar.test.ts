import { fireEvent, render } from "@testing-library/vue";
import { describe, expect, test } from "vitest";
import BatchToolbar from "../components/manager/BatchToolbar.vue";

describe("batch toolbar", () => {
  test("emits action", async () => {
    const { getByTestId, emitted } = render(BatchToolbar, {
      props: {
        selectedCount: 2,
        totalCount: 5,
        allSelected: false,
        lang: "zh"
      }
    });

    await fireEvent.click(getByTestId("batch-start"));
    expect(emitted().action?.[0]).toEqual(["start"]);
  });
});
