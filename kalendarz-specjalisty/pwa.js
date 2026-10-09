(() => {
  let installPrompt = null;
  let reloading = false;

  window.addEventListener('beforeinstallprompt', e => {
    e.preventDefault();
    installPrompt = e;
  });

  async function register() {
    if (!('serviceWorker' in navigator) || location.protocol !== 'https:') return null;
    try {
      const reg = await navigator.serviceWorker.register('./sw.js', {updateViaCache:'none'});
      navigator.serviceWorker.addEventListener('controllerchange', () => {
        if (reloading) return;
        reloading = true;
        location.reload();
      });
      return reg;
    } catch (e) {
      console.warn('PWA registration', e);
      return null;
    }
  }

  let registrationPromise = register();

  window.KSPWA = {
    async install() {
      if (window.matchMedia('(display-mode: standalone)').matches) return true;
      if (!installPrompt) return false;
      installPrompt.prompt();
      const result = await installPrompt.userChoice;
      installPrompt = null;
      return result && result.outcome === 'accepted';
    },
    async checkUpdate() {
      const reg = await registrationPromise;
      if (!reg) return 'unavailable';
      const old = reg.installing || reg.waiting;
      await reg.update();
      await new Promise(r => setTimeout(r, 1200));
      const nw = reg.installing || reg.waiting;
      if (nw && nw !== old) return 'updated';
      return 'current';
    }
  };

  if (navigator.storage && navigator.storage.persist) {
    navigator.storage.persist().catch(()=>{});
  }
})();
