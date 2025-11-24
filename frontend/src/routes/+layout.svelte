<script lang="ts">
  import "../app.css";
  import Toast from "$lib/components/Toast.svelte";
  import { toasts, removeToast } from "$lib/stores/toast";
</script>

<div
  class="min-h-screen bg-slate-50 font-inter text-slate-900 selection:bg-indigo-100 selection:text-indigo-700"
>
  <slot />

  <!-- Toast Container -->
  <div
    class="fixed bottom-4 right-4 z-50 flex flex-col space-y-2 pointer-events-none"
  >
    {#each $toasts as toast (toast.id)}
      <Toast
        message={toast.message}
        type={toast.type}
        onRemove={() => removeToast(toast.id)}
      />
    {/each}
  </div>
</div>

<style>
  :global(body) {
    font-family: "Inter", sans-serif;
  }
</style>
