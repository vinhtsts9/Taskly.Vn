self.addEventListener("push", function (event) {
  const data = event.data ? event.data.json() : {};
  const title = data.title || "Thông báo mới";
  const options = {
    body: data.body || "Bạn có đơn hàng mới!",
    icon: "/icon.png",
    badge: "/badge.png",
  };
  event.waitUntil(self.registration.showNotification(title, options));
});
