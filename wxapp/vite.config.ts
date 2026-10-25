import { defineConfig } from "vite";
import uni from "@dcloudio/vite-plugin-uni";
import path from "node:path";

// 后端跑在宿主机 Mac 上（localhost:8080），容器内要用 host.docker.internal 访问。
// 若后端改在本容器里直接 `go run`，把 VITE_PROXY_TARGET 设为 http://localhost:8080 即可。
const proxyTarget =
  (globalThis as any).process?.env?.VITE_PROXY_TARGET ?? "http://host.docker.internal:8080";

export default defineConfig({
  plugins: [uni()],
  resolve: {
    alias: { "@": path.resolve(__dirname, "src") },
  },
  server: {
    host: "0.0.0.0", // 让宿主 Mac 能打开 http://localhost:5173
    port: 5173,
    proxy: {
      "/api": { target: proxyTarget, changeOrigin: true },
    },
  },
});
