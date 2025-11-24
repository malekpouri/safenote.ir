<script lang="ts">
  import { encryptionService } from "$lib/encryption";
  import { onMount } from "svelte";

  let noteText = "";
  let viewsRemaining = 1;
  let expiration = 1440; // 24 hours
  let createdLink = "";
  let error = "";
  let loading = false;

  onMount(async () => {
    try {
      await encryptionService.init();
    } catch (e) {
      console.error("Failed to initialize WASM:", e);
      error = "Failed to initialize encryption service.";
    }
  });

  async function createNote() {
    loading = true;
    error = "";
    createdLink = "";

    try {
      if (!noteText) {
        throw new Error("Note text is required");
      }

      const key = await encryptionService.generateKey();
      const encryptedData = await encryptionService.encrypt(noteText, key);

      const response = await fetch("/api/notes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          encrypted_data: encryptedData,
          views_remaining: viewsRemaining,
          expiration: expiration,
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || "Failed to create note");
      }

      const data = await response.json();
      const url = new URL(window.location.href);
      url.pathname = `/n/${data.id}`;
      url.hash = key;
      createdLink = url.toString();
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function copyLink() {
    navigator.clipboard.writeText(createdLink);
    alert("Link copied to clipboard!");
  }
</script>

<div class="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
  <div
    class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl p-6"
  >
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900">SafeNote</h1>
      <p class="mt-2 text-gray-600">Share encrypted notes securely.</p>
    </div>

    {#if createdLink}
      <div class="bg-green-50 border border-green-200 rounded-md p-4 mb-6">
        <h3 class="text-lg font-medium text-green-900">Note Created!</h3>
        <p class="mt-2 text-sm text-green-700">
          Share this link with the recipient. It contains the decryption key.
        </p>
        <div class="mt-3 flex items-center space-x-2">
          <input
            type="text"
            readonly
            value={createdLink}
            class="flex-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
          />
          <button
            on:click={copyLink}
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Copy
          </button>
        </div>
        <button
          on:click={() => (createdLink = "")}
          class="mt-4 text-sm text-indigo-600 hover:text-indigo-500"
        >
          Create another note
        </button>
      </div>
    {:else}
      <form on:submit|preventDefault={createNote} class="space-y-6">
        <div>
          <label for="note" class="block text-sm font-medium text-gray-700"
            >Note Content</label
          >
          <div class="mt-1">
            <textarea
              id="note"
              bind:value={noteText}
              rows="5"
              class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
              placeholder="Write your secret note here..."
            ></textarea>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="views" class="block text-sm font-medium text-gray-700"
              >Views Limit</label
            >
            <input
              type="number"
              id="views"
              bind:value={viewsRemaining}
              min="1"
              class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
            />
          </div>
          <div>
            <label
              for="expiration"
              class="block text-sm font-medium text-gray-700">Expiration</label
            >
            <select
              id="expiration"
              bind:value={expiration}
              class="mt-1 block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            >
              <option value={60}>1 Hour</option>
              <option value={1440}>24 Hours</option>
              <option value={10080}>7 Days</option>
              <option value={43200}>30 Days</option>
            </select>
          </div>
        </div>

        {#if error}
          <div class="text-red-600 text-sm">{error}</div>
        {/if}

        <div>
          <button
            type="submit"
            disabled={loading}
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            {loading ? "Encrypting..." : "Create Secure Note"}
          </button>
        </div>
      </form>
    {/if}
  </div>
</div>
