<script lang="ts">
  import { encryptionService } from "$lib/encryption";
  import { addToast } from "$lib/stores/toast";
  import { fade, slide } from "svelte/transition";

  let note = "";
  let views = 1;
  let expiration = 43200; // Default 30 days
  let password = "";
  let showOptions = false;
  let link = "";
  let loading = false;

  async function createNote() {
    loading = true;
    link = "";

    try {
      const shortKey = await encryptionService.generateKey();
      const encryptedData = await encryptionService.encrypt(
        note,
        shortKey,
        password
      );

      const response = await fetch("/api/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
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
      const url = new URL(window.location.href);
      url.pathname = `/n/${data.id}`;
      url.hash = shortKey;
      link = url.toString();

      addToast("Note created successfully!", "success");
    } catch (e: any) {
      console.error(e);
      addToast(e.message || "An error occurred", "error");
    } finally {
      loading = false;
    }
  }

  function copyLink() {
    navigator.clipboard.writeText(link);
    addToast("Link copied to clipboard!", "success");
  }

  function generatePassword() {
    const chars =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*";
    let pass = "";
    const len = Math.floor(Math.random() * (8 - 4 + 1)) + 4; // 4 to 8 chars
    for (let i = 0; i < len; i++) {
      pass += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    password = pass;
    addToast("Random password generated", "info");
  }
</script>

<div
  class="min-h-screen flex flex-col bg-gradient-to-br from-slate-50 to-indigo-50"
>
  <div
    class="flex-grow py-12 px-4 sm:px-6 lg:px-8 flex items-center justify-center"
  >
    <div class="w-full max-w-2xl">
      <div class="text-center mb-10">
        <h1
          class="text-4xl font-extrabold text-slate-900 tracking-tight sm:text-5xl mb-2"
        >
          Safe<span class="text-indigo-600">Note</span>
        </h1>
        <p class="text-lg text-slate-600">
          Share encrypted notes that self-destruct.
        </p>
      </div>

      <div
        class="bg-white rounded-2xl shadow-xl overflow-hidden border border-slate-100 transition-all duration-300 hover:shadow-2xl"
      >
        <div class="p-8">
          {#if link}
            <div in:fade class="space-y-6">
              <div
                class="bg-green-50 border border-green-200 rounded-xl p-6 text-center"
              >
                <div
                  class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4"
                >
                  <svg
                    class="h-6 w-6 text-green-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
                <h3 class="text-xl font-semibold text-green-900">
                  Note Created Successfully!
                </h3>
                <p class="mt-2 text-green-700">
                  Copy the secure link below. It contains the decryption key and
                  will not be shown again.
                </p>
              </div>

              <div>
                <label
                  for="link"
                  class="block text-sm font-medium text-slate-700 mb-2"
                  >Secure Link</label
                >
                <div class="flex rounded-lg shadow-sm">
                  <input
                    type="text"
                    id="link"
                    readonly
                    value={link}
                    class="flex-1 block w-full rounded-l-lg border-slate-300 bg-slate-50 text-slate-600 sm:text-sm focus:ring-indigo-500 focus:border-indigo-500 p-3 border"
                  />
                  <button
                    on:click={copyLink}
                    class="inline-flex items-center px-6 py-3 border border-l-0 border-slate-300 rounded-r-lg bg-indigo-50 text-indigo-700 font-medium hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
                  >
                    Copy
                  </button>
                </div>
              </div>

              <div class="text-center pt-4">
                <button
                  on:click={() => {
                    link = "";
                    note = "";
                    password = "";
                    showOptions = false;
                  }}
                  class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors"
                >
                  Create another note
                </button>
              </div>
            </div>
          {:else}
            <form
              on:submit|preventDefault={createNote}
              class="space-y-6"
              in:fade
            >
              <div>
                <label
                  for="note"
                  class="block text-sm font-medium text-slate-700 mb-2"
                  >Note Content</label
                >
                <textarea
                  id="note"
                  bind:value={note}
                  rows="6"
                  class="shadow-sm block w-full sm:text-sm border-slate-300 rounded-xl focus:ring-indigo-500 focus:border-indigo-500 p-4 border resize-none transition-shadow focus:shadow-md"
                  placeholder="Write your secret note here..."
                  required
                ></textarea>
              </div>

              <div class="flex justify-end">
                <button
                  type="button"
                  class="text-sm font-medium text-indigo-600 hover:text-indigo-800 flex items-center transition-colors"
                  on:click={() => (showOptions = !showOptions)}
                >
                  <svg
                    class="w-4 h-4 mr-1 transform transition-transform duration-200 {showOptions
                      ? 'rotate-180'
                      : ''}"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                  {showOptions ? "Less Options" : "More Options"}
                </button>
              </div>

              {#if showOptions}
                <div
                  transition:slide
                  class="bg-slate-50 p-6 rounded-xl border border-slate-100 space-y-6"
                >
                  <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
                    <div>
                      <label
                        for="views"
                        class="block text-sm font-medium text-slate-700 mb-1"
                        >Views Limit</label
                      >
                      <input
                        type="number"
                        id="views"
                        bind:value={views}
                        min="1"
                        class="block w-full sm:text-sm border-slate-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 p-2.5 border"
                      />
                    </div>

                    <div>
                      <label
                        for="expiration"
                        class="block text-sm font-medium text-slate-700 mb-1"
                        >Expiration</label
                      >
                      <select
                        id="expiration"
                        bind:value={expiration}
                        class="block w-full sm:text-sm border-slate-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 p-2.5 border bg-white"
                      >
                        <option value={60}>1 Hour</option>
                        <option value={1440}>24 Hours</option>
                        <option value={10080}>7 Days</option>
                        <option value={43200}>30 Days</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label
                      for="password"
                      class="block text-sm font-medium text-slate-700 mb-1"
                      >Password Protection (Optional)</label
                    >
                    <div class="flex rounded-lg shadow-sm">
                      <input
                        type="text"
                        id="password"
                        bind:value={password}
                        placeholder="Optional password"
                        class="flex-1 block w-full rounded-l-lg border-slate-300 sm:text-sm focus:ring-indigo-500 focus:border-indigo-500 p-2.5 border"
                      />
                      <button
                        type="button"
                        on:click={generatePassword}
                        class="inline-flex items-center px-4 py-2 border border-l-0 border-slate-300 rounded-r-lg bg-slate-100 text-slate-600 text-sm font-medium hover:bg-slate-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
                        title="Generate Random Password"
                      >
                        <svg
                          class="w-5 h-5"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                          />
                        </svg>
                      </button>
                    </div>
                    <p class="mt-1 text-xs text-slate-500">
                      If set, the recipient must enter this password to decrypt
                      the note.
                    </p>
                  </div>
                </div>
              {/if}

              <button
                type="submit"
                disabled={loading || !note}
                class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:-translate-y-0.5"
              >
                {#if loading}
                  <svg
                    class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Creating...
                {:else}
                  Create Secure Note
                {/if}
              </button>
            </form>
          {/if}
        </div>
      </div>

      <div
        class="max-w-3xl mx-auto mt-12 text-slate-500 text-sm space-y-4 text-center"
      >
        <p>
          SafeNote is a free, web-based service that allows you to send
          top-secret notes over the internet. It's fast, easy, and requires no
          password or user registration at all.
        </p>
      </div>
    </div>
  </div>

  <footer class="bg-white border-t border-slate-200 mt-auto">
    <div class="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <nav
        class="flex flex-wrap justify-center gap-x-8 gap-y-4"
        aria-label="Footer"
      >
        <a
          href="/about"
          class="text-base text-slate-500 hover:text-indigo-600 transition-colors"
          >About</a
        >
        <a
          href="/privacy"
          class="text-base text-slate-500 hover:text-indigo-600 transition-colors"
          >Privacy Policy</a
        >
        <a
          href="mailto:admin@utux.ir"
          class="text-base text-slate-500 hover:text-indigo-600 transition-colors"
          >Contact</a
        >
      </nav>
      <p class="mt-8 text-center text-sm text-slate-400">
        &copy; {new Date().getFullYear()} SafeNote. All rights reserved.
      </p>
    </div>
  </footer>
</div>
