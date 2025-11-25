<script lang="ts">
  import { page } from "$app/stores";
  import { encryptionService } from "$lib/encryption";
  import { onMount } from "svelte";
  import { addToast } from "$lib/stores/toast";
  import { fade, scale } from "svelte/transition";
  import { t } from "$lib/i18n";

  let id = $page.params.id;
  let key = "";
  let noteContent = "";
  let loading = false;
  let confirmed = false;
  let meta: any = null;
  let password = "";
  let showPasswordPrompt = false;
  let deleting = false;

  onMount(async () => {
    const hash = window.location.hash;
    if (hash.startsWith("#")) {
      key = hash.substring(1);
    }

    if (!key) {
      addToast($t.view.toast_key_missing, "error");
      return;
    }

    try {
      await encryptionService.init();
    } catch (e) {
      console.error("Failed to initialize WASM:", e);
      addToast($t.view.toast_init_failed, "error");
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
          throw new Error(data.error || $t.view.error_not_found);
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
      addToast($t.view.toast_decrypted, "success");
    } catch (e: any) {
      let msg = e.message;
      // If decryption failed, it might be wrong password
      if (meta && meta.is_password_protected && msg.includes("cipher")) {
        msg = $t.view.error_wrong_password;
      }
      addToast(msg, "error");
    } finally {
      loading = false;
    }
  }

  async function deleteNote() {
    if (!confirm($t.view.confirm_delete)) return;

    // If password protected and we don't have it, we need to ask.
    // But here we assume if they are on this page, they might want to delete it.
    // If it's password protected, we need the password hash.

    // If we are in "confirmed" state, we have the password.
    // If we are in "showPasswordPrompt" state, we might have typed it but not verified.
    // If we are in initial state, we don't have it.

    // If we don't have meta, fetch it first to know if it's password protected.
    if (!meta) {
      try {
        const response = await fetch(`/api/notes/${id}`);
        if (response.ok) {
          meta = await response.json();
        }
      } catch (e) {}
    }

    if (meta && meta.is_password_protected && !password) {
      // If we don't have password, we must show prompt.
      // If we are not already showing prompt, show it.
      if (!showPasswordPrompt) {
        showPasswordPrompt = true;
        addToast($t.view.password_message, "info");
        return;
      }
      // If we are showing prompt and password is empty, tell them.
      if (!password) {
        addToast($t.view.password_placeholder, "error");
        return;
      }
    }

    deleting = true;
    try {
      let passwordHash = "";
      if (password) {
        passwordHash = await encryptionService.hashPassword(password);
      }

      const response = await fetch(`/api/notes/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          password_hash: passwordHash,
        }),
      });

      if (!response.ok) {
        const data = await response.json().catch(() => ({}));
        if (response.status === 401) {
          throw new Error($t.view.error_wrong_password);
        }
        throw new Error(data.error || "Failed to delete note");
      }

      addToast("Note deleted successfully", "success");
      // Redirect to home
      window.location.href = "/";
    } catch (e: any) {
      addToast(e.message, "error");
    } finally {
      deleting = false;
    }
  }

  function copyText() {
    navigator.clipboard.writeText(noteContent);
    addToast($t.view.toast_copied_text, "success");
  }
</script>

<svelte:head>
  <title>{$t.view.title} - {$t.view.subtitle}</title>
  <meta name="description" content={$t.app.description} />
  <meta property="og:title" content="{$t.view.title} - {$t.view.subtitle}" />
  <meta property="og:description" content={$t.app.description} />
</svelte:head>

<div
  class="min-h-screen flex flex-col bg-gradient-to-br from-slate-50 to-indigo-50"
>
  <div
    class="flex-grow py-8 px-4 sm:px-6 lg:px-8 flex items-center justify-center"
  >
    <div class="w-full max-w-2xl">
      <div class="text-center mb-8 sm:mb-10">
        <h1
          class="text-2xl sm:text-4xl font-extrabold text-slate-900 tracking-tight lg:text-5xl mb-2"
        >
          Safe<span class="text-indigo-600">Note</span>
        </h1>
        <p class="text-base sm:text-lg text-slate-600">{$t.view.subtitle}</p>
      </div>

      <div
        class="bg-white rounded-2xl shadow-xl overflow-hidden border border-slate-100 transition-all duration-300 hover:shadow-2xl"
      >
        <div class="p-6 sm:p-8">
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
                  {$t.view.warning_destroyed}
                  <span class="font-bold">{$t.view.warning_copy}</span>
                  {$t.view.warning_suffix}
                </p>
              </div>

              <div>
                <label
                  for="content"
                  class="block text-sm font-medium text-slate-700 mb-2"
                  >{$t.view.label_content}</label
                >
                <textarea
                  id="content"
                  readonly
                  value={noteContent}
                  rows="10"
                  class="shadow-sm block w-full text-sm border-slate-300 rounded-xl focus:ring-indigo-500 focus:border-indigo-500 p-4 border bg-slate-50 resize-none font-mono text-slate-800"
                ></textarea>
              </div>

              <div
                class="flex flex-col sm:flex-row justify-between items-center gap-4 pt-2"
              >
                <div class="flex gap-2">
                  <button
                    on:click={copyText}
                    class="inline-flex items-center px-4 py-2 border border-slate-300 rounded-lg bg-white text-slate-700 font-medium hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
                  >
                    {$t.view.button_copy_text}
                  </button>
                </div>

                <button
                  on:click={() => (window.location.href = "/")}
                  class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 font-medium transition-colors"
                >
                  {$t.view.button_new}
                </button>
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
                  {$t.view.password_title}
                </h3>
                <p class="mt-1 text-slate-500">
                  {$t.view.password_message}
                </p>
              </div>

              <div>
                <label for="password" class="sr-only">Password</label>
                <input
                  type="password"
                  id="password"
                  bind:value={password}
                  class="block w-full text-sm border-slate-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 p-3 border text-center tracking-widest"
                  placeholder={$t.view.password_placeholder}
                />
              </div>

              <div class="space-y-3">
                <button
                  on:click={viewNote}
                  disabled={loading || !password}
                  class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-all transform hover:-translate-y-0.5"
                >
                  {loading ? $t.view.button_decrypting : $t.view.button_unlock}
                </button>

                <button
                  on:click={deleteNote}
                  disabled={deleting || !password}
                  class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl text-base font-medium text-red-600 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 transition-all"
                >
                  {#if deleting}
                    <svg
                      class="animate-spin -ml-1 mr-2 h-5 w-5 text-red-600"
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
                  {/if}
                  {$t.view.button_delete}
                </button>
              </div>
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
                <h3 class="text-xl font-bold text-slate-900">
                  {$t.view.ready_title}
                </h3>
                <p class="mt-2 text-slate-600 max-w-sm mx-auto">
                  {$t.view.ready_message}
                </p>
              </div>

              <div class="space-y-3">
                <button
                  on:click={viewNote}
                  disabled={loading || !key}
                  class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 transition-all transform hover:-translate-y-0.5"
                >
                  {loading ? $t.view.button_decrypting : $t.view.button_show}
                </button>

                <button
                  on:click={deleteNote}
                  disabled={deleting}
                  class="w-full flex justify-center py-3 px-4 border border-transparent rounded-xl text-base font-medium text-red-600 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 transition-all"
                >
                  {#if deleting}
                    <svg
                      class="animate-spin -ml-1 mr-2 h-5 w-5 text-red-600"
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
                  {/if}
                  {$t.view.button_delete}
                </button>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
