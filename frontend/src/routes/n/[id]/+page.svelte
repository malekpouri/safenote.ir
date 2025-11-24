<script lang="ts">
  import { page } from "$app/stores";
  import { encryptionService } from "$lib/encryption";
  import { onMount } from "svelte";

  let id = $page.params.id;
  let key = "";
  let noteContent = "";
  let error = "";
  let loading = false;
  let confirmed = false;
  let meta: any = null;
  let password = "";
  let showPasswordPrompt = false;

  onMount(async () => {
    const hash = window.location.hash;
    if (hash.startsWith("#")) {
      key = hash.substring(1);
    }

    if (!key) {
      error = "Decryption key missing from URL.";
      return;
    }

    try {
      await encryptionService.init();
    } catch (e) {
      console.error("Failed to initialize WASM:", e);
      error = "Failed to initialize encryption service.";
    }
  });

  async function viewNote() {
    loading = true;
    error = "";

    try {
      // If we don't have metadata yet, fetch it
      if (!meta) {
        const response = await fetch(`/api/notes/${id}`);
        if (!response.ok) {
          const data = await response.json().catch(() => ({}));
          throw new Error(data.error || "Note not found or expired");
        }
        meta = await response.json();
      }

      // Check for password protection
      if (meta.is_password_protected && !password) {
        showPasswordPrompt = true;
        loading = false;
        return;
      }

      // Decrypt
      noteContent = await encryptionService.decrypt(
        meta.encrypted_data,
        key,
        password
      );
      confirmed = true;
      showPasswordPrompt = false;
    } catch (e: any) {
      error = e.message;
      // If decryption failed, it might be wrong password
      if (meta && meta.is_password_protected) {
        error = "Decryption failed. Wrong password?";
      }
    } finally {
      loading = false;
    }
  }
</script>

<div class="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
  <div
    class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl p-6"
  >
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">SafeNote</h1>
      <p class="mt-2 text-gray-600">Secure Note Viewer</p>
    </div>

    {#if error}
      <div class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
        <h3 class="text-lg font-medium text-red-900">Error</h3>
        <p class="mt-2 text-sm text-red-700">{error}</p>
      </div>
    {:else if confirmed}
      <div class="space-y-4">
        <div class="bg-yellow-50 border border-yellow-200 rounded-md p-4">
          <p class="text-sm text-yellow-800">
            This note has been destroyed from the server (if view limit was 1).
            Copy it now if you need to keep it.
          </p>
        </div>

        <div>
          <label for="content" class="block text-sm font-medium text-gray-700"
            >Note Content</label
          >
          <div class="mt-1">
            <textarea
              id="content"
              readonly
              value={noteContent}
              rows="10"
              class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border bg-gray-50"
            ></textarea>
          </div>
        </div>

        <div class="text-xs text-gray-500 text-right">
          Created at: {new Date(meta.created_at).toLocaleString()}
        </div>

        <div class="text-center">
          <a href="/" class="text-indigo-600 hover:text-indigo-500 text-sm"
            >Create a new note</a
          >
        </div>
      </div>
    {:else if showPasswordPrompt}
      <div class="space-y-6">
        <p class="text-gray-700 text-center">
          This note is password protected. Please enter the password to view it.
        </p>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700"
            >Password</label
          >
          <input
            type="password"
            id="password"
            bind:value={password}
            class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
          />
        </div>

        <button
          on:click={viewNote}
          disabled={loading || !password}
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? "Decrypting..." : "Unlock Note"}
        </button>
      </div>
    {:else}
      <div class="text-center space-y-6">
        <p class="text-gray-700">
          You are about to view a secure note. If this note has a view limit,
          viewing it will consume a view and might destroy the note.
        </p>

        <button
          on:click={viewNote}
          disabled={loading || !key}
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? "Decrypting..." : "Yes, show me the note"}
        </button>
      </div>
    {/if}
  </div>
</div>
