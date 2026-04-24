import path from "node:path";
import { chromium } from "patchright-core";
import { describe, expect, it } from "vitest";

describe("javascript tests", async () => {
  const browser = await chromium.launch({
    executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
    headless: false,
    args: ["--user-agent-version=99.0.4844.51"],
  });

  const page = await browser.newPage();
  await page.goto("https://browserleaks.com/javascript");

  it("should set user agent version", async () => {
    const userAgent = await page.evaluate(() => navigator.userAgent);
    expect(userAgent).toContain("Chrome/99.0.0.0");

    const highEntropyValues = await page.evaluate(async () => {
      // @ts-ignore
      return await navigator.userAgentData.getHighEntropyValues(["fullVersionList"]);
    });

    const chromeVersion = highEntropyValues.fullVersionList?.find((item: any) => item.brand === "Google Chrome");
    expect(chromeVersion?.version).toBe("99.0.4844.51");
  });

  it("should prevent document idle", async () => {
    const result = await page.evaluate(() => {
      return document.visibilityState;
    });

    expect(result).toBe("visible");
  });
});
