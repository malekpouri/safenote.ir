<script lang="ts">
  import { encryptionService } from "$lib/encryption";
  import { onMount } from "svelte";

  let note = "";
  let views = 1;
  let expiration = 43200; // Default 30 days
  let password = "";
  let showOptions = false;
  let link = "";
  let loading = false;
  let error = "";

  async function createNote() {
    loading = true;
    error = "";
    link = "";

    try {
      // 1. Generate Short Key (6 chars)
      const shortKey = await encryptionService.generateKey();

      // 2. Encrypt locally
      const encryptedData = await encryptionService.encrypt(
        note,
        shortKey,
        password
      );

      // 3. Send to backend
      const response = await fetch("/api/notes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          encrypted_data: encryptedData,
          views_remaining: views,
          expiration: expiration,
          is_password_protected: !!password,
        }),
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({}));
        throw new Error(data.error || "Failed to create note");
      }

      const data = await response.json();

      // 4. Construct Link
      const url = new URL(window.location.href);
      url.pathname = `/n/${data.id}`;
      url.hash = shortKey;
      link = url.toString();
    } catch (e: any) {
      console.error(e);
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function copyLink() {
    navigator.clipboard.writeText(link);
    alert("Link copied to clipboard!");
  }
</script>

<div class="min-h-screen bg-gray-100 flex flex-col">
  <div class="flex-grow py-12 px-4 sm:px-6 lg:px-8">
    <div
      class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl p-6"
    >
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">SafeNote</h1>
        <p class="mt-2 text-gray-600">Share encrypted notes securely.</p>
      </div>

      {#if link}
        <div class="space-y-4">
          <div class="bg-green-50 border border-green-200 rounded-md p-4">
            <h3 class="text-lg font-medium text-green-900">Note Created!</h3>
            <p class="mt-2 text-sm text-green-700">
              Here is your secure link. It contains the decryption key, so share
              it carefully.
            </p>
          </div>

          <div>
            <label for="link" class="block text-sm font-medium text-gray-700"
              >Secure Link</label
            >
            <div class="mt-1 flex rounded-md shadow-sm">
              <input
                type="text"
                id="link"
                readonly
                value={link}
                class="focus:ring-indigo-500 focus:border-indigo-500 flex-1 block w-full rounded-none rounded-l-md sm:text-sm border-gray-300 p-2 border bg-gray-50"
              />
              <button
                on:click={copyLink}
                class="-ml-px relative inline-flex items-center space-x-2 px-4 py-2 border border-gray-300 text-sm font-medium rounded-r-md text-gray-700 bg-gray-50 hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500"
              >
                Copy
              </button>
            </div>
          </div>

          <div class="text-center">
            <button
              on:click={() => {
                link = "";
                note = "";
                password = "";
              }}
              class="text-indigo-600 hover:text-indigo-500 text-sm"
            >
              Create another note
            </button>
          </div>
        </div>
      {:else}
        <form on:submit|preventDefault={createNote} class="space-y-6">
          {#if error}
            <div class="bg-red-50 border border-red-200 rounded-md p-4">
              <p class="text-sm text-red-700">{error}</p>
            </div>
          {/if}

          <div>
            <label for="note" class="block text-sm font-medium text-gray-700"
              >Note Content</label
            >
            <div class="mt-1">
              <textarea
                id="note"
                bind:value={note}
                rows="5"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
                placeholder="Write your secret note here..."
                required
              ></textarea>
            </div>
          </div>

          <div class="text-right">
            <button
              type="button"
              class="text-sm text-indigo-600 hover:text-indigo-500"
              on:click={() => (showOptions = !showOptions)}
            >
              {showOptions ? "Less Options" : "More Options"}
            </button>
          </div>

          {#if showOptions}
            <div
              class="grid grid-cols-1 gap-4 sm:grid-cols-2 bg-gray-50 p-4 rounded-md"
            >
              <div>
                <label
                  for="views"
                  class="block text-sm font-medium text-gray-700"
                  >Views Limit</label
                >
                <input
                  type="number"
                  id="views"
                  bind:value={views}
                  min="1"
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
                />
              </div>

              <div>
                <label
                  for="expiration"
                  class="block text-sm font-medium text-gray-700"
                  >Expiration</label
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

              <div class="sm:col-span-2">
                <label
                  for="password"
                  class="block text-sm font-medium text-gray-700"
                  >Password Protection (Optional)</label
                >
                <input
                  type="password"
                  id="password"
                  bind:value={password}
                  placeholder="Enter a password to encrypt the note"
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md p-2 border"
                />
                <p class="mt-1 text-xs text-gray-500">
                  If set, the recipient must enter this password to decrypt the
                  note.
                </p>
              </div>
            </div>
          {/if}

          <div>
            <button
              type="submit"
              disabled={loading || !note}
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              {loading ? "Creating..." : "Create Secure Note"}
            </button>
          </div>
        </form>
      {/if}
    </div>

    <div class="max-w-4xl mx-auto mt-12 text-gray-600 text-sm space-y-4">
      <p>
        If you need to send a password or some other form of simple but
        sensitive information to someone you can not send it over IM or email.
        These methods are not secure, anyone with little knowledge can intercept
        this information during transmission. Using the private note as the
        'middle man' you can safely and securely transfer this data to your
        recipient. The private note in other words self-destruct note is a
        web-based service that allows you to share a note or a message over the
        internet confidently. Write your private message, and the system will
        generate a secure link. Then you can copy that link into an email or
        instant message then you can send it to the person you want to read the
        note. When that person first clicks the link they will see the note in
        their browser and the note will immediately delete itself, which ensures
        that no one (even the same person) will read the note again. The link
        will no longer be accessible.
      </p>
    </div>
  </div>

  <footer class="bg-white border-t border-gray-200 mt-auto">
    <div class="max-w-7xl mx-auto py-6 px-4 overflow-hidden sm:px-6 lg:px-8">
      <nav
        class="-mx-5 -my-2 flex flex-wrap justify-center"
        aria-label="Footer"
      >
        <div class="px-5 py-2">
          <a href="/about" class="text-base text-gray-500 hover:text-gray-900">
            About
          </a>
        </div>
        <div class="px-5 py-2">
          <a
            href="/privacy"
            class="text-base text-gray-500 hover:text-gray-900"
          >
            Privacy Policy
          </a>
        </div>
        <div class="px-5 py-2">
          <a
            href="mailto:admin@utux.ir"
            class="text-base text-gray-500 hover:text-gray-900"
          >
            Contact
          </a>
        </div>
      </nav>
      <p class="mt-8 text-center text-base text-gray-400">
        &copy; {new Date().getFullYear()} SafeNote. All rights reserved.
      </p>
    </div>
  </footer>
</div>
