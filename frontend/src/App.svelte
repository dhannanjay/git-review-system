<script>
  import { onMount } from 'svelte'
  import { ListChanges, LoadPatch } from '../wailsjs/go/main/App.js'
  import { Quit } from '../wailsjs/runtime/runtime.js'

  let files = []
  let selectedFile = ''
  let patch = null
  let filesLoaded = false
  let patchLoading = false
  let error = ''
  let currentHunkIndex = 0
  let diffAreaEl

  onMount(async () => {
    try {
      files = await ListChanges()
    } catch (e) {
      error = String(e)
    }
    filesLoaded = true
    if (files.length > 0) {
      selectedFile = files[0].path
      await loadPatch(selectedFile)
    }
  })

  async function loadPatch(path) {
    selectedFile = path
    patchLoading = true
    error = ''
    currentHunkIndex = 0
    try {
      patch = await LoadPatch(path)
    } catch (e) {
      error = String(e)
      patch = null
    } finally {
      patchLoading = false
    }
  }

  function selectFile(index) {
    if (index >= 0 && index < files.length) {
      loadPatch(files[index].path)
    }
  }

  function currentFileIndex() {
    return files.findIndex(f => f.path === selectedFile)
  }

  function focusHunk(index) {
    if (!patch || index < 0 || index >= patch.hunks.length) return
    currentHunkIndex = index
    const hunkEls = diffAreaEl.querySelectorAll('.hunk-container')
    if (hunkEls[index]) {
      hunkEls[index].scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }

  function handleKeydown(e) {
    const tagName = e.target.tagName.toLowerCase()
    if (tagName === 'input' || tagName === 'textarea' || tagName === 'select') return

    switch (e.key) {
      case 'q':
        Quit()
        break
      case 'j':
      case ']':
        selectFile(currentFileIndex() + 1)
        break
      case 'k':
      case '[':
        selectFile(currentFileIndex() - 1)
        break
      case 'n':
        focusHunk(currentHunkIndex + 1)
        break
      case 'p':
        focusHunk(currentHunkIndex - 1)
        break
      default:
        return
    }
    e.preventDefault()
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<main>
  <div class="sidebar" id="sidebar">
    <h2>Files</h2>
    {#if !filesLoaded}
      <p class="empty">Loading files...</p>
    {:else if files.length === 0}
      <p class="empty">No changed files</p>
    {:else}
      {#each files as f}
        <button
          class="file-item"
          class:active={f.path === selectedFile}
          on:click={() => loadPatch(f.path)}
        >
          <span class="status status-{f.status.toLowerCase()}">{f.status}</span>
          <span class="path">
            {#if f.oldPath}
              <span class="old-path">{f.oldPath}</span>
              <span class="arrow">→</span>
            {/if}
            {f.path}
          </span>
        </button>
      {/each}
    {/if}
  </div>

  <div class="diff-area" id="diff-area" bind:this={diffAreaEl}>
    {#if patchLoading}
      <p class="status-msg">Loading diff...</p>
    {:else if error}
      <p class="error-msg">{error}</p>
    {:else if patch}
      {#if patch.isBinary}
        <div class="stub-msg">Binary file changed</div>
      {:else if patch.isSubmodule}
        <div class="stub-msg">Submodule changed</div>
      {:else if patch.hunks.length > 0}
        {#each patch.hunks as hunk, i}
          <div class="hunk-container" class:hunk-active={i === currentHunkIndex}>
            <div class="hunk-header">
              {hunk.header}
              {#if patch.hunks.length > 1}
                <span class="hunk-index">({i + 1}/{patch.hunks.length})</span>
              {/if}
            </div>
            <div class="diff-table">
              <div class="diff-row header-row">
                <div class="line-num-left">Old</div>
                <div class="line-num-right">New</div>
                <div class="line-content">Content</div>
              </div>
              {#each hunk.lines as line}
                <div
                  class="diff-row"
                  class:added={line.type === 1}
                  class:removed={line.type === 2}
                >
                  <div class="line-num-left">{line.oldNum || ''}</div>
                  <div class="line-num-right">{line.newNum || ''}</div>
                  <div class="line-content">{line.content}</div>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      {:else}
        <p class="status-msg">No textual diff available</p>
      {/if}
    {:else}
      <p class="status-msg">Select a file to view its diff</p>
    {/if}
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    color: #d4d4d4;
    background: #1e1e1e;
  }

  :global(#app) {
    height: 100vh;
    display: flex;
  }

  main {
    display: flex;
    height: 100vh;
    width: 100%;
  }

  .sidebar {
    width: 240px;
    min-width: 240px;
    background: #252526;
    border-right: 1px solid #3c3c3c;
    overflow-y: auto;
    padding: 8px;
  }

  .sidebar h2 {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: #888;
    margin: 0 0 8px;
    padding: 0 8px;
  }

  .file-item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 6px 8px;
    border: none;
    background: none;
    color: #ccc;
    cursor: pointer;
    font-size: 13px;
    text-align: left;
    font-family: inherit;
    border-radius: 3px;
  }

  .file-item:hover {
    background: #2a2d2e;
  }

  .file-item.active {
    background: #37373d;
  }

  .status {
    font-size: 10px;
    font-weight: bold;
    padding: 1px 4px;
    border-radius: 3px;
    flex-shrink: 0;
  }

  .status-m { background: #4ec9b0; color: #000; }
  .status-a { background: #6a9955; color: #fff; }
  .status-d { background: #f14c4c; color: #fff; }
  .status-r { background: #dcdcaa; color: #000; }

  .path {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .old-path {
    opacity: 0.7;
  }

  .arrow {
    margin: 0 2px;
    opacity: 0.7;
  }

  .diff-area {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  .hunk-header {
    font-size: 12px;
    color: #569cd6;
    padding: 4px 8px;
    background: #1e1e1e;
    border-bottom: 1px solid #3c3c3c;
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .hunk-active .hunk-header {
    background: #2a2d50;
    border-color: #569cd6;
  }

  .hunk-index {
    font-size: 10px;
    color: #888;
    font-weight: normal;
  }

  .hunk-container {
    scroll-margin-top: 8px;
  }

  .diff-table {
    font-family: 'Courier New', Courier, monospace;
    font-size: 13px;
    line-height: 1.5;
  }

  .diff-row {
    display: flex;
    min-height: 20px;
  }

  .diff-row.header-row {
    font-weight: bold;
    color: #888;
    border-bottom: 1px solid #3c3c3c;
    padding-bottom: 2px;
    margin-bottom: 2px;
  }

  .line-num-left,
  .line-num-right {
    width: 48px;
    min-width: 48px;
    text-align: right;
    padding: 0 8px;
    color: #858585;
    user-select: none;
  }

  .line-content {
    flex: 1;
    white-space: pre;
    padding: 0 4px;
  }

  .diff-row.added {
    background: #1e3a1e;
  }

  .diff-row.added .line-num-right {
    color: #6a9955;
  }

  .diff-row.removed {
    background: #3a1e1e;
  }

  .diff-row.removed .line-num-left {
    color: #f14c4c;
  }

  .stub-msg {
    color: #dcdcaa;
    font-size: 14px;
    font-weight: bold;
    text-align: center;
    padding: 40px;
  }

  .status-msg {
    color: #888;
    text-align: center;
    padding: 40px;
  }

  .error-msg {
    color: #f14c4c;
    padding: 16px;
  }

  .empty {
    color: #888;
    font-size: 13px;
    padding: 8px;
  }
</style>
