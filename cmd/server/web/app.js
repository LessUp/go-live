(() => {
  const storageKeys = {
    room: 'go-live.room',
    token: 'go-live.token',
  };

  class HTTPError extends Error {
    constructor(message, status) {
      super(message);
      this.name = 'HTTPError';
      this.status = status;
    }
  }

  let bootstrapCache = null;

  function qs(id) {
    return document.getElementById(id);
  }

  function normalizeRoom(value) {
    const room = String(value || '').trim();
    return room || 'demo';
  }

  function readQueryParam(name) {
    return new URLSearchParams(window.location.search).get(name) || '';
  }

  function initRoomAndToken(roomInput, tokenInput) {
    if (roomInput) {
      const roomFromQuery = readQueryParam('room');
      const roomFromStorage = window.localStorage.getItem(storageKeys.room) || '';
      roomInput.value = normalizeRoom(roomFromQuery || roomFromStorage || roomInput.value);
      roomInput.addEventListener('input', () => {
        window.localStorage.setItem(storageKeys.room, normalizeRoom(roomInput.value));
      });
    }

    if (tokenInput) {
      const tokenFromStorage = window.sessionStorage.getItem(storageKeys.token) || '';
      tokenInput.value = tokenFromStorage || tokenInput.value;
      tokenInput.addEventListener('input', () => {
        const token = tokenInput.value.trim();
        if (token) {
          window.sessionStorage.setItem(storageKeys.token, token);
        } else {
          window.sessionStorage.removeItem(storageKeys.token);
        }
      });
    }
  }

  function createHeaders(token, baseHeaders = {}) {
    const headers = new Headers(baseHeaders);
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
    return headers;
  }

  async function readErrorMessage(response) {
    const text = (await response.text()).trim();
    if (text) {
      return text;
    }
    return `请求失败（${response.status}）`;
  }

  async function fetchJSON(url, options = {}) {
    const response = await fetch(url, options);
    if (!response.ok) {
      throw new HTTPError(await readErrorMessage(response), response.status);
    }
    return response.json();
  }

  async function fetchText(url, options = {}) {
    const response = await fetch(url, options);
    if (!response.ok) {
      throw new HTTPError(await readErrorMessage(response), response.status);
    }
    return response.text();
  }

  async function getBootstrap(force = false) {
    if (!force && bootstrapCache) {
      return bootstrapCache;
    }
    bootstrapCache = await fetchJSON('/api/bootstrap');
    return bootstrapCache;
  }

  function setStatus(target, message, tone = 'info') {
    if (!target) {
      return;
    }
    target.hidden = !message;
    target.textContent = message || '';
    target.className = `status status--${tone}`;
  }

  function clearLog(target) {
    if (target) {
      target.textContent = '';
    }
  }

  function appendLog(target, message) {
    if (!target || !message) {
      return;
    }
    const timestamp = new Date().toLocaleTimeString('zh-CN', { hour12: false });
    target.textContent += `[${timestamp}] ${message}\n`;
    target.scrollTop = target.scrollHeight;
  }

  function formatBytes(bytes) {
    const units = ['B', 'KB', 'MB', 'GB'];
    let value = Number(bytes) || 0;
    let index = 0;
    while (value >= 1024 && index < units.length - 1) {
      value /= 1024;
      index += 1;
    }
    return `${value.toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
  }

  function formatDate(value) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value || '-';
    }
    return date.toLocaleString('zh-CN', { hour12: false });
  }

  function buildPageLink(path, room) {
    const url = new URL(path, window.location.origin);
    if (room) {
      url.searchParams.set('room', normalizeRoom(room));
    }
    return `${url.pathname}${url.search}`;
  }

  function describeBootstrap(bootstrap) {
    const items = [];
    if (bootstrap.authEnabled) {
      items.push('已开启鉴权');
    } else {
      items.push('未开启鉴权');
    }
    if (bootstrap.recordEnabled) {
      items.push('已开启录制');
    } else {
      items.push('未开启录制');
    }
    items.push(`ICE 服务器 ${Array.isArray(bootstrap.iceServers) ? bootstrap.iceServers.length : 0} 个`);
    return items.join(' · ');
  }

  window.GoLiveApp = {
    HTTPError,
    qs,
    normalizeRoom,
    readQueryParam,
    initRoomAndToken,
    createHeaders,
    fetchJSON,
    fetchText,
    getBootstrap,
    setStatus,
    clearLog,
    appendLog,
    formatBytes,
    formatDate,
    buildPageLink,
    describeBootstrap,
  };
})();
