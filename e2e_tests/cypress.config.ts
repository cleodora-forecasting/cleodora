import { defineConfig } from "cypress";

export default defineConfig({
  env: {
    cleocPath: '../dist/cleoc_linux_amd64_v1/cleoc',
  },
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    baseUrl: 'http://localhost:8080',
  },
});
