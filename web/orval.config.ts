import { defineConfig } from "orval";

export default defineConfig({
  "virsh-sandbox-api": {
    output: {
      client: "react-query",
      mode: "tags-split",
      target: "src/virsh-sandbox",
      schemas: "src/virsh-sandbox/model",
      mock: true,
    },
    input: {
      target: "../virsh-sandbox/docs/swagger.yaml",
    },
  },
  // "virsh-sandbox-api-zod": {
  //   output: {
  //     client: "zod",
  //     mode: "tags-split",
  //     target: "./src/virsh-sandbox-client",
  //     fileExtension: ".zod.ts",
  //   },
  //   input: {
  //     target: "../virsh-sandbox/docs/swagger.yaml",
  //   },
  // },
  "tmux-client": {
    output: {
      client: "react-query",
      mode: "tags-split",
      target: "./src/tmux-client",
      schemas: "./src/tmux-client/model",
      mock: true,
    },
    input: {
      target: "../tmux-client/docs/swagger.yaml",
    },
  },
  // "tmux-client-zod": {
  //   output: {
  //     client: "zod",
  //     mode: "tags-split",
  //     target: "./src/tmux-client",
  //     fileExtension: ".zod.ts",
  //   },
  //   input: {
  //     target: "../tmux-client/docs/swagger.yaml",
  //   },
  // },
});
