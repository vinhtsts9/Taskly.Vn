import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: ["e12d412ecfe6.ngrok-free.app"],
    proxy: {
      // Proxy tất cả các request HTTP bắt đầu bằng /v1
      "/v1": {
        target: "https://taskly-vn-2.onrender.com", // Địa chỉ backend của bạn
        changeOrigin: true,
        secure: false,
      },
      // Proxy yêu cầu WebSocket
      "/ws": {
        target: "https://taskly-vn-2.onrender.com", // Địa chỉ WebSocket backend
        ws: true,
        // Viết lại đường dẫn: Bỏ /ws và thay bằng /v1/2024/ws
        rewrite: (path) => path.replace(/^\/ws/, "/v1/2024/ws"),
      },
    },
  },
});
