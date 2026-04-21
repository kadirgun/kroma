import path from "node:path";
import { chromium } from "patchright-core";
import { describe, expect, it } from "vitest";

const getCanvasHash = async (noised: boolean) => {
  const args = [];
  if (noised) {
    args.push("--canvas-noise=123456789");
  }

  const browser = await chromium.launch({
    executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
    headless: false,
    args,
  });

  const page = await browser.newPage();
  await page.goto("https://browserleaks.com/canvas");

  const canvasHash = await page.locator("#canvas-hash").textContent();

  return canvasHash;
};

describe("canvas fingerprint", () => {
  it("should be add noise to canvas fingerprint", async () => {
    const [originalHash, noisedHash] = await Promise.all([getCanvasHash(false), getCanvasHash(true)]);

    console.log("Original Hash:", originalHash);
    console.log("Noised Hash:", noisedHash);

    expect(originalHash).not.toBe(noisedHash);
  });
});
