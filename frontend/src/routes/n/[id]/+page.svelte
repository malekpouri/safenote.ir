<script lang="ts">
  import { page } from "$app/stores";
  import { encryptionService } from "$lib/encryption";
  import { onMount } from "svelte";
  import { addToast } from "$lib/stores/toast";
  import { fade, scale } from "svelte/transition";

  let id = $page.params.id;
  let key = "";
  let noteContent = "";
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
      addToast("Decryption key missing from URL.", "error");
      return;
    }

    try {
      await encryptionService.init();
    } catch (e) {
      console.error("Failed to initialize WASM:", e);
      addToast("Failed to initialize encryption service.", "error");
    }
  });

  async function viewNote() {
    loading = true;

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
      addToast("Note decrypted successfully", "success");
    } catch (e: any) {
      let msg = e.message;
      // If decryption failed, it might be wrong password
      if (meta && meta.is_password_protected && msg.includes("cipher")) {
        msg = "Decryption failed. Wrong password?";
      }
      addToast(msg, "error");
    } finally {
      loading = false;
    }
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
        <p class="text-lg text-slate-600">Secure Note Viewer</p>
      </div>

      <div
        class="bg-white rounded-2xl shadow-xl overflow-hidden border border-slate-100 transition-all duration-300 hover:shadow-2xl"
      >
        <div class="p-8">
          {#if confirmed}
            <div in:fade class="space-y-6">
              <div
                class="bg-yellow-50 border border-yellow-200 rounded-xl p-4 flex items-start"
              >
                <svg
                  class="h-6 w-6 text-yellow-600 mr-3 mt-0.5 flex-shrink-0"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
                <p class="text-sm text-yellow-800">
                  This note has been destroyed from the server (if view limit
                  was 1).
                  <span class="font-bold">Copy it now</span> if you need to keep
                  it.
                </p>
              </div>

              <div>
                <label
                  for="content"
                  class="block text-sm font-medium text-slate-700 mb-2"
                  >Note Content</label
                >
                <textarea
                  id="content"
                  readonly
                  value={noteContent}
                  rows="10"
                  class="shadow-sm block w-full sm:text-sm border-slate-300 rounded-xl focus:ring-indigo-500 focus:border-indigo-500 p-4 border bg-slate-50 resize-none font-mono text-slate-800"
                ></textarea>
              </div>

              <div
                class="flex justify-between items-center text-xs text-slate-400 pt-2"
              >
                <span
                  >Created: {new Date(meta.created_at).toLocaleString()}</span
                >
                <a
                  href="/"
                  class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors"
                  >Create new note &rarr;</a
                >
              </div>
            </div>
          {:else if showPasswordPrompt}
            <div in:scale class="space-y-6">
              <div class="text-center">
                <div
                  class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-indigo-100 mb-4"
                >
                  <svg
                    class="h-6 w-6 text-indigo-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                    />
                  </svg>
                </div>
                <h3 class="text-lg font-medium text-slate-900">
                  Password Protected
                </h3>
                <p class="mt-1 text-slate-500">
                  Please enter the password to unlock this note.
                </p>
              </div>

              <div>
                <label for="password" class="sr-only">Password</label>
                <input
                  type="password"
                  id="password"
                  bind:value={password}
                  class="block w-full sm:text-sm border-slate-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 p-3 border text-center tracking-widest"
                  placeholder="Enter Password"
                />
              </div>

              <button
                on:click={viewNote}
                disabled={loading || !password}
                class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-all transform hover:-translate-y-0.5"
              >
                {loading ? "Decrypting..." : "Unlock Note"}
              </button>
            </div>
          {:else}
            <div in:fade class="text-center space-y-8 py-4">
              <div>
                <div
                  class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-indigo-50 mb-6"
                >
                  <svg
                    class="h-8 w-8 text-indigo-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                </div>
                <h3 class="text-xl font-bold text-slate-900">Ready to View?</h3>
                <p class="mt-2 text-slate-600 max-w-sm mx-auto">
                  You are about to view a secure note. Viewing it may
                  permanently destroy it from the server.
                </p>
              </div>

              <button
                on:click={viewNote}
                disabled={loading || !key}
                class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-all transform hover:-translate-y-0.5"
              >
                {loading ? "Decrypting..." : "Yes, show me the note"}
              </button>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
