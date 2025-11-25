<script lang="ts">
  import { onMount } from "svelte";
  import { t } from "$lib/i18n";
  import { fade } from "svelte/transition";

  let stats = {
    active_notes: 0,
    total_notes_created: 0,
  };
  let loading = true;

  async function fetchStats() {
    loading = true;
    try {
      const response = await fetch("/api/admin/stats");
      if (response.ok) {
        stats = await response.json();
      }
    } catch (e) {
      console.error("Failed to fetch stats", e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    fetchStats();
  });
</script>

<svelte:head>
  <title>{$t.admin.title} - {$t.app.title}</title>
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
          class="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight lg:text-5xl mb-2"
        >
          Safe<span class="text-indigo-600">Note</span>
        </h1>
        <p class="text-base sm:text-lg text-slate-600">{$t.admin.subtitle}</p>
      </div>

      <div
        class="bg-white rounded-2xl shadow-xl overflow-hidden border border-slate-100 transition-all duration-300 hover:shadow-2xl"
      >
        <div class="p-6 sm:p-8">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <div
              class="bg-indigo-50 rounded-xl p-6 text-center border border-indigo-100"
            >
              <div class="text-indigo-600 font-medium mb-2">
                {$t.admin.active_notes}
              </div>
              <div class="text-4xl font-bold text-slate-900">
                {#if loading}
                  ...
                {:else}
                  {stats.active_notes}
                {/if}
              </div>
            </div>

            <div
              class="bg-green-50 rounded-xl p-6 text-center border border-green-100"
            >
              <div class="text-green-600 font-medium mb-2">
                {$t.admin.total_notes}
              </div>
              <div class="text-4xl font-bold text-slate-900">
                {#if loading}
                  ...
                {:else}
                  {stats.total_notes_created}
                {/if}
              </div>
            </div>
          </div>

          <div class="mt-8 flex justify-center">
            <button
              on:click={fetchStats}
              disabled={loading}
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 font-medium transition-colors disabled:opacity-50"
            >
              {#if loading}
                <svg
                  class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
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
              {$t.admin.refresh}
            </button>
          </div>

          <div class="mt-6 text-center">
            <a
              href="/"
              class="text-sm text-slate-500 hover:text-indigo-600 transition-colors"
            >
              {$t.app.title}
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
