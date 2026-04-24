import path from "node:path";
import { chromium } from "patchright-core";
import { describe, expect, it } from "vitest";

describe("keyboard tests", () => {
  it("should support altgr and numlock modifiers", async () => {
    const browser = await chromium.launch({
      executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
      headless: false,
    });

    const page = await browser.newPage();
    await page.goto("https://example.com");

    await page.evaluate(() => {
      (window as any).keyboardEvents = [];
      window.addEventListener("keydown", (event) => {
        (window as any).keyboardEvents.push({
          modifiers: ["AltGraph", "NumLock"].filter((key) => event.getModifierState(key)),
        });
      });
    });

    const client = await page.context().newCDPSession(page);
    await client.send("Input.dispatchKeyEvent", {
      type: "keyDown",
      modifiers: 48,
      windowsVirtualKeyCode: 0,
      code: "KeyA",
    });

    const events = await page.evaluate(() => (window as any).keyboardEvents);
    console.log(events);

    expect(events).toEqual([{ modifiers: ["AltGraph", "NumLock"] }]);

    await browser.close();
  });
});
