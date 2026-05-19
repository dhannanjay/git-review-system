<script>
  import { onMount } from 'svelte'
  import { ListChanges, LoadPatch } from '../wailsjs/go/main/App.js'

  let files = []
  let selectedFile = ''
  let patch = null
  let filesLoaded = false
  let patchLoading = false
  let error = ''

  onMount(async () => {
    try {
      files = await ListChanges()
      if (files.length > 0) {
        selectedFile = files[0].path
        await loadPatch(selectedFile)
      }
    } catch (e) {
      error = String(e)
    } finally {
      filesLoaded = true
    }
  })

  async function loadPatch(path) {
    selectedFile = path
    patchLoading = true
    error = ''
    try {
      patch = await LoadPatch(path)
    } catch (e) {
      error = String(e)
      patch = null
    } finally {
      patchLoading = false
    }
  }
</script>

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
          <span class="path">{f.path}</span>
        </button>
      {/each}
    {/if}
  </div>

  <div class="diff-area" id="diff-area">
    {#if patchLoading}
      <p class="status-msg">Loading diff...</p>
    {:else if error}
      <p class="error-msg">{error}</p>
    {:else if patch}
      {#each patch.hunks as hunk}
        <div class="hunk-header">{hunk.header}</div>
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
      {/each}
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
  }

  .status-m { background: #4ec9b0; color: #000; }
  .status-a { background: #6a9955; color: #fff; }
  .status-d { background: #f14c4c; color: #fff; }

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
