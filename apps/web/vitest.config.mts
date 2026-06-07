import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import { fileURLToPath } from "node:url";

const r = (p: string) => fileURLToPath(new URL(p, import.meta.url));

export default defineConfig({
  plugins: [react()],
  resolve: {
    // Monorepo has React both at root and under apps/web. Alias to a single
    // copy so the component and @testing-library share one React instance
    // (otherwise the hooks dispatcher is null during render).
    dedupe: ["react", "react-dom"],
    alias: {
      react: r("../../node_modules/react"),
      "react-dom": r("../../node_modules/react-dom"),
    },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./vitest.setup.ts"],
    coverage: {
      provider: "v8",
      include: ["app/**/*.{ts,tsx}", "components/**/*.{ts,tsx}"],
      exclude: [
        "app/layout.tsx",
        "**/*.d.ts",
        "**/*.test.{ts,tsx}",
        "vitest.setup.ts",
      ],
      reporter: ["text-summary", "text"],
      thresholds: { lines: 80 },
    },
  },
});
