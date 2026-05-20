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
  let showHelp = false

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
      case '?':
        showHelp = !showHelp
        e.preventDefault()
        return
      case 'Escape':
        if (showHelp) {
          showHelp = false
          e.preventDefault()
          return
        }
        return
      case 'q':
        if (showHelp) return
        Quit()
        break
      case 'j':
      case ']':
        if (showHelp) return
        selectFile(currentFileIndex() + 1)
        break
      case 'k':
      case '[':
        if (showHelp) return
        selectFile(currentFileIndex() - 1)
        break
      case 'n':
        if (showHelp) return
        focusHunk(currentHunkIndex + 1)
        break
      case 'p':
        if (showHelp) return
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

{#if showHelp}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="help-overlay" on:click={() => showHelp = false}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="help-content" on:click|stopPropagation>
      <h2>Keyboard Shortcuts</h2>
      <table>
        <tr>
          <td><kbd>j</kbd> / <kbd>]</kbd></td>
          <td>Next file</td>
        </tr>
        <tr>
          <td><kbd>k</kbd> / <kbd>[</kbd></td>
          <td>Previous file</td>
        </tr>
        <tr>
          <td><kbd>n</kbd></td>
          <td>Next hunk</td>
        </tr>
        <tr>
          <td><kbd>p</kbd></td>
          <td>Previous hunk</td>
        </tr>
        <tr>
          <td><kbd>q</kbd></td>
          <td>Quit application</td>
        </tr>
        <tr>
          <td><kbd>?</kbd></td>
          <td>Toggle this help overlay</td>
        </tr>
        <tr>
          <td><kbd>Esc</kbd></td>
          <td>Close help overlay</td>
        </tr>
      </table>
      <p class="help-footer">Press <kbd>?</kbd> or <kbd>Esc</kbd> to close</p>
    </div>
  </div>
{/if}

<style>
  :global(body) {
    margin: 0;
    font-family: 'Courier New', Courier, monospace;
    color: var(--text-primary);
    background: var(--bg-primary);
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
    background: var(--bg-secondary);
    border-right: 1px solid var(--border-color);
    overflow-y: auto;
    padding: 8px;
  }

  .sidebar h2 {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-muted);
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
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 13px;
    text-align: left;
    font-family: inherit;
    border-radius: 3px;
  }

  .file-item:hover {
    background: var(--bg-hover);
  }

  .file-item.active {
    background: var(--bg-active);
  }

  .status {
    font-size: 10px;
    font-weight: bold;
    padding: 1px 4px;
    border-radius: 3px;
    flex-shrink: 0;
  }

  .status-m { background: var(--status-m-bg); color: var(--status-m-text); }
  .status-a { background: var(--status-a-bg); color: var(--status-a-text); }
  .status-d { background: var(--status-d-bg); color: var(--status-d-text); }
  .status-r { background: var(--status-r-bg); color: var(--status-r-text); }

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
    color: var(--text-header);
    padding: 4px 8px;
    background: var(--bg-primary);
    border-bottom: 1px solid var(--border-color);
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .hunk-active .hunk-header {
    background: var(--bg-hunk-active);
    border-color: var(--text-header);
  }

  .hunk-index {
    font-size: 10px;
    color: var(--text-muted);
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
    color: var(--text-muted);
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 2px;
    margin-bottom: 2px;
  }

  .line-num-left,
  .line-num-right {
    width: 48px;
    min-width: 48px;
    text-align: right;
    padding: 0 8px;
    color: var(--text-line-num);
    user-select: none;
  }

  .line-content {
    flex: 1;
    white-space: pre;
    padding: 0 4px;
  }

  .diff-row.added {
    background: var(--bg-added);
  }

  .diff-row.added .line-num-right {
    color: var(--text-added);
  }

  .diff-row.removed {
    background: var(--bg-removed);
  }

  .diff-row.removed .line-num-left {
    color: var(--text-removed);
  }

  .stub-msg {
    color: var(--text-stub);
    font-size: 14px;
    font-weight: bold;
    text-align: center;
    padding: 40px;
  }

  .status-msg {
    color: var(--text-muted);
    text-align: center;
    padding: 40px;
  }

  .error-msg {
    color: var(--text-removed);
    padding: 16px;
  }

  .empty {
    color: var(--text-muted);
    font-size: 13px;
    padding: 8px;
  }

  .help-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: var(--overlay-bg);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }

  .help-content {
    background: var(--overlay-content-bg);
    border: 1px solid var(--overlay-border);
    border-radius: 8px;
    padding: 24px 32px;
    min-width: 320px;
    box-shadow: 0 4px 24px rgba(0,0,0,0.3);
  }

  .help-content h2 {
    margin: 0 0 16px;
    font-size: 14px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-muted);
  }

  .help-content table {
    width: 100%;
    border-collapse: collapse;
  }

  .help-content tr {
    border-bottom: 1px solid var(--border-color);
  }

  .help-content tr:last-child {
    border-bottom: none;
  }

  .help-content td {
    padding: 8px 8px;
    font-size: 13px;
  }

  .help-content td:first-child {
    white-space: nowrap;
    padding-right: 24px;
    color: var(--text-primary);
  }

  .help-content td:last-child {
    color: var(--text-muted);
  }

  kbd {
    display: inline-block;
    padding: 2px 6px;
    font-size: 12px;
    font-family: 'Courier New', Courier, monospace;
    background: var(--key-bg);
    border: 1px solid var(--key-border);
    border-radius: 3px;
    color: var(--text-primary);
  }

  .help-footer {
    margin: 16px 0 0;
    font-size: 11px;
    color: var(--text-muted);
    text-align: center;
  }
</style>
