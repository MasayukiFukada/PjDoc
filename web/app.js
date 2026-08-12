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

  if (window.mermaid) {
    mermaid.initialize({ startOnLoad: false, theme: 'dark' });
  }

  loadTree();

  themeToggle.addEventListener('click', () => {
    const isLight = document.body.getAttribute('data-theme') === 'light';
    if (isLight) {
      document.body.removeAttribute('data-theme');
    } else {
      document.body.setAttribute('data-theme', 'light');
    }
  });

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
