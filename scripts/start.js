import { execSync } from "child_process";
import os from "os";

// Get first non-internal IPv4 (works on Wi-Fi, hotspot, ethernet)
function getLocalIp() {
  const interfaces = os.networkInterfaces();
  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name]) {
      if (iface.family === "IPv4" && !iface.internal) {
        return iface.address;
      }
    }
  }
  return "127.0.0.1"; // fallback
}

const ip = getLocalIp();
console.log(`🚀 Starting Expo on LAN (${ip})`);

try {
  execSync(
    `REACT_NATIVE_PACKAGER_HOSTNAME=${ip} npx expo start --host lan -c`,
    { stdio: "inherit" }
  );
} catch (err) {
  console.error("❌ Failed to start Expo:", err.message);
}
