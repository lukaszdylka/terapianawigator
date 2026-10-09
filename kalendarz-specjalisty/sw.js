const CACHE = 'kalendarz-specjalisty-pwa-v1';
const SHELL = ['./', './index.html', './manifest.webmanifest', './icon.svg'];

self.addEventListener('install', event => {
  self.skipWaiting();
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(SHELL)).catch(() => {}));
});

self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(keys => Promise.all(keys.filter(k => k.startsWith('kalendarz-specjalisty-pwa-') && k !== CACHE).map(k => caches.delete(k))))
      .then(() => self.clients.claim())
  );
});

self.addEventListener('fetch', event => {
  const req = event.request;
  if (req.method !== 'GET') return;
  const url = new URL(req.url);
  if (url.origin !== self.location.origin) return;

  event.respondWith((async () => {
    try {
      const fresh = await fetch(req, { cache: 'no-store' });
      if (fresh && fresh.ok) {
        const cache = await caches.open(CACHE);
        let key = req;
        if (url.pathname.endsWith('/index.html') || req.mode === 'navigate') {
          key = new Request(new URL('./index.html', self.location.href).href);
        }
        cache.put(key, fresh.clone()).catch(() => {});
      }
      return fresh;
    } catch {
      const cache = await caches.open(CACHE);
      if (req.mode === 'navigate') {
        return (await cache.match('./index.html')) || (await cache.match('./')) || Response.error();
      }
      return (await cache.match(req)) || Response.error();
    }
  })());
});
