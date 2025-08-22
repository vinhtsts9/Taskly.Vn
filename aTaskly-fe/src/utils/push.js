export async function registerPush() {
  if (!("serviceWorker" in navigator)) {
    console.error("Trình duyệt không hỗ trợ Service Worker");
    return;
  }

  // 1. Đăng ký service worker
  const registration = await navigator.serviceWorker.register("/sw.js");

  // 2. Lấy VAPID public key từ environment variables
  const key = import.meta.env.VITE_PUBLIC_KEY_WEB_NOTI;
  if (!key) {
    console.error("VITE_PUBLIC_KEY_WEB_NOTI is not defined in .env file");
    return;
  }

  // 3. Đăng ký push
  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlBase64ToUint8Array(key),
  });

  // 4. Gửi subscription lên BE
  await fetch("/save-subscription", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(subscription),
  });

  console.log("Push subscription:", subscription);
}

function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
  const rawData = window.atob(base64);
  return Uint8Array.from([...rawData].map((c) => c.charCodeAt(0)));
}
