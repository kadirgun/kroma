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

const getWebGLHash = async (args?: string[]) => {
  const browser = await chromium.launch({
    executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
    headless: false,
    args,
  });

  const page = await browser.newPage();
  await page.goto("https://browserleaks.com/webgl");

  const webglHash = await page.locator("#gl-image-hash").textContent();
  const webglVendor = await page.locator("#UNMASKED_VENDOR_WEBGL").textContent();
  const webglRenderer = await page.locator("#UNMASKED_RENDERER_WEBGL").textContent();

  return {
    webglHash,
    webglVendor,
    webglRenderer,
  };
};

describe("canvas fingerprint", () => {
  it("should be add noise to canvas fingerprint", async () => {
    const [originalHash, noisedHash] = await Promise.all([getCanvasHash(false), getCanvasHash(true)]);

    console.log("Original Hash:", originalHash);
    console.log("Noised Hash:", noisedHash);

    expect(originalHash).not.toBe(noisedHash);
  });

  it("should be add noise to webgl fingerprint", async () => {
    const [original, fake] = await Promise.all([
      getWebGLHash(),
      getWebGLHash(["--canvas-noise=123456789", "--webgl-renderer=Noised Renderer", "--webgl-vendor=Noised Vendor"]),
    ]);

    console.log("Original:", original);
    console.log("Noised:", fake);

    expect(original.webglHash).not.toBe(fake.webglHash);
    expect(fake.webglVendor).contain("Noised Vendor");
    expect(fake.webglRenderer).contain("Noised Renderer");
  });
});
