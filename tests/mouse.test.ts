import path from "node:path";
import { chromium } from "patchright-core";
import { describe, it } from "vitest";

describe("mouse tests", () => {
  it("should schedule mouse events", async () => {
    const browser = await chromium.launch({
      executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
      headless: false,
    });

    const page = await browser.newPage();
    await page.goto("https://example.com");

    const client = await page.context().newCDPSession(page);

    // Create list to record mouse events on page
    await page.evaluate(() => {
      (window as any).mouseEvents = [];
      document.addEventListener("mousemove", (e) => {
        (window as any).mouseEvents.push({
          x: e.clientX,
          y: e.clientY,
          time: Date.now(),
        });
      });
    });

    const events: any[] = [];
    const startTime = Date.now();
    for (let i = 0; i < 200; i++) {
      events.push({
        type: "mouseMoved",
        x: 100 + i,
        y: 100 + i,
        timestamp: startTime + i * 10,
      });
    }

    // @ts-ignore
    await client.send("Input.scheduleMouseEvents", {
      events,
    });

    // Wait for all events to complete (200 * 10ms = 2 seconds + buffer)
    await page.waitForTimeout(2500);

    // Kaydedilen eventleri kontrol et
    const recordedEvents = await page.evaluate(() => (window as any).mouseEvents);
    console.log("Recorded events count:", recordedEvents.length);
    console.log("First event:", recordedEvents[0]);
    console.log("Last event:", recordedEvents[recordedEvents.length - 1]);

    const timeDiff =
      recordedEvents.length > 1 ? recordedEvents[recordedEvents.length - 1].time - recordedEvents[0].time : 0;
    console.log("Total time span:", timeDiff, "ms");
  });
});
