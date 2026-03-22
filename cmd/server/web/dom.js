(() => {
  function clear(node) {
    while (node && node.firstChild) {
      node.removeChild(node.firstChild);
    }
  }

  function el(tagName, options = {}, children = []) {
    const node = document.createElement(tagName);
    if (options.className) {
      node.className = options.className;
    }
    if (options.text !== undefined) {
      node.textContent = options.text;
    }
    if (options.attrs) {
      Object.entries(options.attrs).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          node.setAttribute(key, value);
        }
      });
    }
    children.forEach((child) => {
      if (child) {
        node.appendChild(child);
      }
    });
    return node;
  }

  function badge(text, tone = 'info') {
    return el('span', { className: `badge badge--${tone}`, text });
  }

  window.GoLiveDOM = {
    clear,
    el,
    badge,
  };
})();
