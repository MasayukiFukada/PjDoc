document.addEventListener('DOMContentLoaded', () => {
  let currentPath = '';
  let currentQuery = '';

  const fileTree = document.getElementById('fileTree');
  const searchInput = document.getElementById('searchInput');
  const searchResults = document.getElementById('searchResults');
  const markdownBody = document.getElementById('markdownBody');
  const breadcrumb = document.getElementById('breadcrumb');
  const tocList = document.getElementById('tocList');
  const themeToggle = document.getElementById('themeToggle');
  const themeMenu = document.getElementById('themeMenu');
  const hljsTheme = document.getElementById('hljsTheme');

  const THEME_NAMES = {
    dark: 'Dark',
    light: 'Light',
    sepia: 'Sepia',
    cyber: 'Cyberpunk',
    system: 'Auto'
  };

  const THEME_ICONS = {
    dark: `<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>`,
    light: `<circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>`,
    sepia: `<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>`,
    cyber: `<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>`,
    system: `<rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line>`
  };

  const HLJS_THEMES = {
    dark: 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css',
    light: 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css',
    sepia: 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/solarized-light.min.css',
    cyber: 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/tokyo-night-dark.min.css'
  };

  const MERMAID_THEMES = {
    dark: 'dark',
    light: 'default',
    sepia: 'neutral',
    cyber: 'dark'
  };

  const mediaQueryDark = window.matchMedia('(prefers-color-scheme: dark)');

  function getEffectiveTheme(setting) {
    if (setting === 'system') {
      return mediaQueryDark.matches ? 'dark' : 'light';
    }
    return setting || 'dark';
  }

  function applyTheme(selectedSetting) {
    const effectiveTheme = getEffectiveTheme(selectedSetting);

    // document attribute
    document.body.setAttribute('data-theme', effectiveTheme);

    // LocalStorage
    localStorage.setItem('pjdoc-theme', selectedSetting);

    // Active status in menu
    document.querySelectorAll('.theme-option').forEach(opt => {
      opt.classList.toggle('active', opt.dataset.themeValue === selectedSetting);
    });

    // Update label text in button
    const currentThemeNameEl = document.getElementById('currentThemeName');
    if (currentThemeNameEl) {
      currentThemeNameEl.textContent = THEME_NAMES[selectedSetting] || 'Dark';
    }

    // Update Highlight.js theme
    if (hljsTheme) {
      const cssUrl = HLJS_THEMES[effectiveTheme] || HLJS_THEMES.dark;
      if (hljsTheme.getAttribute('href') !== cssUrl) {
        hljsTheme.setAttribute('href', cssUrl);
      }
    }

    // Update Mermaid default theme setting
    if (window.mermaid) {
      mermaid.initialize({
        startOnLoad: false,
        theme: MERMAID_THEMES[effectiveTheme] || 'dark'
      });
      // Re-render mermaid diagrams if present
      const mermaidNodes = markdownBody.querySelectorAll('.mermaid');
      if (mermaidNodes.length > 0) {
        try {
          mermaid.run({ nodes: mermaidNodes });
        } catch (e) {
          console.error('Mermaid theme update error:', e);
        }
      }
    }
  }

  // Load saved theme or default to system
  const savedTheme = localStorage.getItem('pjdoc-theme') || 'system';
  applyTheme(savedTheme);

  // OS theme change listener for 'system' theme
  mediaQueryDark.addEventListener('change', () => {
    const currentSetting = localStorage.getItem('pjdoc-theme') || 'system';
    if (currentSetting === 'system') {
      applyTheme('system');
    }
  });

  // Toggle Menu
  themeToggle.addEventListener('click', (e) => {
    e.stopPropagation();
    const isHidden = themeMenu.classList.contains('hidden');
    if (isHidden) {
      themeMenu.classList.remove('hidden');
      themeToggle.setAttribute('aria-expanded', 'true');
    } else {
      themeMenu.classList.add('hidden');
      themeToggle.setAttribute('aria-expanded', 'false');
    }
  });

  // Theme option clicks
  document.querySelectorAll('.theme-option').forEach(btn => {
    btn.addEventListener('click', (e) => {
      e.stopPropagation();
      const themeValue = btn.dataset.themeValue;
      applyTheme(themeValue);
      themeMenu.classList.add('hidden');
      themeToggle.setAttribute('aria-expanded', 'false');
    });
  });

  // Close menu on outside click
  document.addEventListener('click', (e) => {
    if (!themeMenu.contains(e.target) && !themeToggle.contains(e.target)) {
      themeMenu.classList.add('hidden');
      themeToggle.setAttribute('aria-expanded', 'false');
    }
  });

  loadTree();

  let searchTimeout = null;
  searchInput.addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    currentQuery = e.target.value.trim();
    if (!currentQuery) {
      searchResults.classList.add('hidden');
      fileTree.classList.remove('hidden');
      if (markdownBody) {
        removeHighlights(markdownBody);
      }
      return;
    }

    searchTimeout = setTimeout(() => {
      performSearch(currentQuery);
      if (markdownBody) {
        highlightInElement(markdownBody, currentQuery);
      }
    }, 250);
  });

  async function loadTree() {
    try {
      const res = await fetch('/api/tree');
      const data = await res.json();
      fileTree.innerHTML = '';
      if (data && data.children) {
        renderTreeNodes(data.children, fileTree);
      }
    } catch (err) {
      console.error('Failed to load file tree:', err);
    }
  }

  function renderTreeNodes(nodes, container) {
    nodes.forEach(node => {
      const itemEl = document.createElement('div');
      itemEl.className = 'tree-item';
      
      const icon = node.isDir ? '📁' : '📄';
      itemEl.innerHTML = `<span>${icon}</span><span>${escapeHtml(node.name)}</span>`;
      
      if (!node.isDir) {
        itemEl.addEventListener('click', () => {
          document.querySelectorAll('.tree-item.active').forEach(el => el.classList.remove('active'));
          itemEl.classList.add('active');
          loadDocument(node.path);
        });
        container.appendChild(itemEl);
      } else {
        const dirWrapper = document.createElement('div');
        dirWrapper.appendChild(itemEl);
        
        const childrenContainer = document.createElement('div');
        childrenContainer.className = 'tree-children';
        renderTreeNodes(node.children, childrenContainer);
        
        itemEl.addEventListener('click', () => {
          childrenContainer.classList.toggle('hidden');
        });

        dirWrapper.appendChild(childrenContainer);
        container.appendChild(dirWrapper);
      }
    });
  }

  async function loadDocument(path) {
    currentPath = path;
    breadcrumb.textContent = path;
    try {
      const res = await fetch(`/api/document?path=${encodeURIComponent(path)}`);
      if (!res.ok) {
        markdownBody.innerHTML = `<div class="empty-state"><h2>Error</h2><p>ドキュメントの読み込みに失敗しました (${res.status})</p></div>`;
        return;
      }
      const data = await res.json();
      renderMarkdown(data.content);
    } catch (err) {
      console.error('Failed to load document:', err);
    }
  }

  let plantumlServer = 'https://kroki.io';

  async function loadConfig() {
    try {
      const res = await fetch('/api/config');
      if (res.ok) {
        const config = await res.json();
        if (config.plantumlServer) {
          plantumlServer = config.plantumlServer;
        }
      }
    } catch (e) {
      console.warn('Failed to load server config, using default PlantUML server:', e);
    }
  }
  loadConfig();

  function encode64(data) {
    let r = "";
    for (let i = 0; i < data.length; i += 3) {
      if (i + 2 < data.length) {
        r += append3bytes(data[i], data[i + 1], data[i + 2]);
      } else if (i + 1 < data.length) {
        r += append3bytes(data[i], data[i + 1], 0);
      } else {
        r += append3bytes(data[i], 0, 0);
      }
    }
    return r;
  }

  function encode6bit(b) {
    if (b < 10) return String.fromCharCode(48 + b);
    b -= 10;
    if (b < 26) return String.fromCharCode(65 + b);
    b -= 26;
    if (b < 26) return String.fromCharCode(97 + b);
    b -= 26;
    if (b === 0) return '-';
    if (b === 1) return '_';
    return '?';
  }

  function append3bytes(b1, b2, b3) {
    const c1 = b1 >> 2;
    const c2 = ((b1 & 0x3) << 4) | (b2 >> 4);
    const c3 = ((b2 & 0xF) << 2) | (b3 >> 6);
    const c4 = b3 & 0x3F;
    return encode6bit(c1 & 0x3F) + encode6bit(c2 & 0x3F) + encode6bit(c3 & 0x3F) + encode6bit(c4 & 0x3F);
  }

  function encodePlantUML(text) {
    if (!window.pako) return null;
    const textEncoder = new TextEncoder();
    const data = textEncoder.encode(text);
    const compressed = pako.deflateRaw(data, { level: 9 });
    return encode64(compressed);
  }

  function renderMarkdown(content) {
    if (!window.marked) return;

    const rawHtml = marked.parse(content);
    markdownBody.innerHTML = rawHtml;

    markdownBody.querySelectorAll('pre code').forEach((block) => {
      hljs.highlightElement(block);
    });

    markdownBody.querySelectorAll('pre code.language-mermaid').forEach((block) => {
      const code = block.textContent;
      const mermaidContainer = document.createElement('div');
      mermaidContainer.className = 'mermaid';
      mermaidContainer.textContent = code;
      block.parentElement.replaceWith(mermaidContainer);
    });
    if (window.mermaid) {
      try {
        mermaid.run({ nodes: markdownBody.querySelectorAll('.mermaid') });
      } catch (e) {
        console.error('Mermaid render error:', e);
      }
    }

    markdownBody.querySelectorAll('pre code.language-plantuml, pre code.language-puml').forEach((block) => {
      let code = block.textContent.trim();
      if (!code.startsWith('@start')) {
        code = `@startuml\n${code}\n@enduml`;
      }
      const encoded = encodePlantUML(code);
      if (encoded) {
        const baseUrl = plantumlServer.replace(/\/+$/, '');
        const imgUrl = `${baseUrl}/plantuml/svg/${encoded}`;
        
        const container = document.createElement('div');
        container.className = 'plantuml-container';
        
        const img = document.createElement('img');
        img.className = 'plantuml-diagram';
        img.src = imgUrl;
        img.alt = 'PlantUML Diagram';
        img.loading = 'lazy';
        img.onerror = () => {
          container.innerHTML = `<div class="plantuml-error">PlantUML ダイアグラムの描画に失敗しました。サーバー設定 (<code>${escapeHtml(baseUrl)}</code>) を確認してください。</div>`;
        };
        
        container.appendChild(img);
        block.parentElement.replaceWith(container);
      }
    });

    if (window.renderMathInElement) {
      renderMathInElement(markdownBody, {
        delimiters: [
          {left: '$$', right: '$$', display: true},
          {left: '$', right: '$', display: false},
          {left: '\\(', right: '\\)', display: false},
          {left: '\\[', right: '\\]', display: true}
        ]
      });
    }

    // ドキュメント読み込み時に現在検索中のキーワードがあれば本文内をハイライト
    if (currentQuery) {
      highlightInElement(markdownBody, currentQuery);
    }

    generateTOC();
  }

  function generateTOC() {
    tocList.innerHTML = '';
    const headings = markdownBody.querySelectorAll('h1, h2, h3');
    headings.forEach((heading, idx) => {
      const id = `heading-${idx}`;
      heading.id = id;

      const link = document.createElement('a');
      link.href = `#${id}`;
      link.className = `toc-link depth-${heading.tagName.charAt(1)}`;
      link.textContent = heading.textContent;
      
      link.addEventListener('click', (e) => {
        e.preventDefault();
        heading.scrollIntoView({ behavior: 'smooth' });
      });

      tocList.appendChild(link);
    });
  }

  async function performSearch(query) {
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
      const results = await res.json();

      searchResults.innerHTML = '';
      if (!results || results.length === 0) {
        searchResults.innerHTML = '<div class="search-item-snippet">該当する結果はありませんでした</div>';
      } else {
        results.forEach(res => {
          const item = document.createElement('div');
          item.className = 'search-item';
          const highlightedPath = highlightText(res.path, query);
          const highlightedSnippet = res.snippet ? highlightText(res.snippet, query) : '';
          
          item.innerHTML = `
            <div class="search-item-file">${highlightedPath} ${res.lineNumber ? `(L${res.lineNumber})` : ''}</div>
            ${highlightedSnippet ? `<div class="search-item-snippet">${highlightedSnippet}</div>` : ''}
          `;
          item.addEventListener('click', () => {
            loadDocument(res.path);
          });
          searchResults.appendChild(item);
        });
      }

      fileTree.classList.add('hidden');
      searchResults.classList.remove('hidden');
    } catch (err) {
      console.error('Search failed:', err);
    }
  }

  // テキストハイライト補助関数
  function highlightText(text, query) {
    if (!query) return escapeHtml(text);
    const escapedText = escapeHtml(text);
    const escapedQuery = escapeRegExp(query);
    const regex = new RegExp(`(${escapedQuery})`, 'gi');
    return escapedText.replace(regex, '<mark class="search-highlight">$1</mark>');
  }

  // DOM内テキストノードハイライト関数
  function highlightInElement(element, query) {
    removeHighlights(element);
    if (!query) return;

    const regex = new RegExp(`(${escapeRegExp(query)})`, 'gi');
    const walker = document.createTreeWalker(element, NodeFilter.SHOW_TEXT, {
      acceptNode: (node) => {
        if (!node.nodeValue.trim()) return NodeFilter.FILTER_REJECT;
        if (node.parentElement && (node.parentElement.tagName === 'SCRIPT' || node.parentElement.tagName === 'STYLE' || node.parentElement.classList.contains('search-highlight'))) {
          return NodeFilter.FILTER_REJECT;
        }
        return NodeFilter.FILTER_ACCEPT;
      }
    });

    const nodesToReplace = [];
    while (walker.nextNode()) {
      const node = walker.currentNode;
      if (regex.test(node.nodeValue)) {
        nodesToReplace.push(node);
      }
    }

    let firstMatch = null;
    nodesToReplace.forEach(node => {
      const parent = node.parentNode;
      const fragment = document.createDocumentFragment();
      const parts = node.nodeValue.split(regex);

      parts.forEach(part => {
        if (part.toLowerCase() === query.toLowerCase()) {
          const mark = document.createElement('mark');
          mark.className = 'search-highlight';
          mark.textContent = part;
          fragment.appendChild(mark);
          if (!firstMatch) firstMatch = mark;
        } else {
          fragment.appendChild(document.createTextNode(part));
        }
      });

      parent.replaceChild(fragment, node);
    });

    if (firstMatch) {
      firstMatch.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }

  function removeHighlights(element) {
    element.querySelectorAll('mark.search-highlight').forEach(mark => {
      const parent = mark.parentNode;
      parent.replaceChild(document.createTextNode(mark.textContent), mark);
      parent.normalize();
    });
  }

  function escapeRegExp(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  }

  function escapeHtml(str) {
    return str.replace(/[&<>"']/g, (m) => ({
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#039;'
    }[m]));
  }

  const eventSource = new EventSource('/api/events');
  eventSource.onmessage = (e) => {
    if (e.data === 'reload') {
      loadTree();
      if (currentPath) {
        loadDocument(currentPath);
      }
    }
  };
});
