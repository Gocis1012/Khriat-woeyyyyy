import { describe, it, expect } from "vitest";
import { TARGETS, LEVELS, targetLabel, levelLabel } from "./translateOptions";

describe("translateOptions", () => {
  it("has 3 targets and 5 levels", () => {
    expect(TARGETS).toHaveLength(3);
    expect(LEVELS).toHaveLength(5);
  });

  it("targetLabel resolves known values and falls back to the raw value", () => {
    expect(targetLabel("boss")).toBe("หัวหน้า");
    expect(targetLabel("client")).toBe("ลูกค้า");
    expect(targetLabel("friend")).toBe("เพื่อน");
    expect(targetLabel("nope")).toBe("nope");
  });

  it("levelLabel resolves known values and falls back to default", () => {
    expect(levelLabel(1)).toBe("สวมวิญญาณผู้ดี");
    expect(levelLabel(3)).toBe("มนุษย์ปกติ");
    expect(levelLabel(5)).toBe("ตัวแม่จะแคร์เพื่อ");
    expect(levelLabel(999)).toBe("มนุษย์ปกติ");
  });

  it("every level has a feeling and an example", () => {
    for (const lvl of LEVELS) {
      expect(lvl.feeling.length).toBeGreaterThan(0);
      expect(lvl.example.length).toBeGreaterThan(0);
      expect(lvl.value).toBeGreaterThanOrEqual(1);
      expect(lvl.value).toBeLessThanOrEqual(5);
    }
  });
});
