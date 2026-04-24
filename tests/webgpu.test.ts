import path from "node:path";
import { chromium } from "patchright-core";
import { describe, expect, it } from "vitest";

describe("webgpu tests", () => {
  it("should set vendor and arch", async () => {
    const browser = await chromium.launch({
      executablePath: path.join(__dirname, "../chromium/src/out/Dev/chrome.exe"),
      headless: false,
      args: ["--gpu-vendor=TestVendor", "--gpu-arch=TestArch"],
    });

    const page = await browser.newPage();
    await page.goto("https://browserleaks.com/webgpu");

    const gpuInfo = await page.evaluate(async () => {
      return navigator.gpu.requestAdapter().then((a) => ({
        vendor: a?.info.vendor,
        architecture: a?.info.architecture,
      }));
    });

    console.log(gpuInfo);

    expect(gpuInfo?.vendor).toBe("TestVendor");
    expect(gpuInfo?.architecture).toBe("TestArch");
  });
});
