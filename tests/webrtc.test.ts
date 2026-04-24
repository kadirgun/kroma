import path from "node:path";
import { chromium } from "patchright-core";
import { describe, expect, it } from "vitest";

describe("webrtc tests", () => {
  it("should set visitor-ip", async () => {
    const browser = await chromium.launch({
      executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
      headless: false,
      args: ["--visitor-ip=54.54.54.54"],
    });

    const page = await browser.newPage();
    await page.goto("https://browserleaks.com/webrtc");

    const rtcPublicLocator = page.locator("#rtc-public");
    await rtcPublicLocator.waitFor({ state: "visible" });

    const ipLocator = page.locator(`a`, { hasText: "54.54.54.54" });
    await ipLocator.waitFor({ state: "visible", timeout: 30000 });
    const ipText = await ipLocator.textContent();
    console.log("IP Text:", ipText);

    expect(ipText?.trim()).toBe("54.54.54.54");
  });
});
