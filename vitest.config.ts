import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    globals: true,
    environment: "node",
    isolate: false,
    include: ["tests/**/*.test.ts"],
    watch: false,
    testTimeout: 1000 * 60 * 60,
  },
});
