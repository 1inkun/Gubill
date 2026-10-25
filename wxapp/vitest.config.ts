// 独立的测试配置（不复用 vite.config.ts：那个文件会加载 uni 插件）。
// 注意：这里故意不用 `import { defineConfig } from 'vitest/config'`，
// 因为该 import 需要从本项目 node_modules 解析 vitest；普通对象即可。
// `@` 别名指向 src，与 tsconfig paths 保持一致。
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.dirname(fileURLToPath(import.meta.url));

export default {
  // 不让 esbuild 去解析项目 tsconfig（它 extends 的 @vue/tsconfig 在 vitest 的
  // tsconfck 下解析会失败）；测试转译不需要 tsconfig，别名走下面的 resolve.alias。
  esbuild: {
    tsconfigRaw: "{}",
  },
  resolve: {
    alias: {
      "@": path.resolve(rootDir, "src"),
    },
  },
  test: {
    environment: "node",
    include: ["test/**/*.spec.ts"],
    globals: false,
  },
};
