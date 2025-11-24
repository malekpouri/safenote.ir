<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import { onMount } from "svelte";

  export let message: string;
  export let type: "success" | "error" | "info" = "info";
  export let duration = 3000;
  export let onRemove: () => void;

  onMount(() => {
    const timer = setTimeout(() => {
      onRemove();
    }, duration);

    return () => clearTimeout(timer);
  });

  const colors = {
    success: "bg-green-500",
    error: "bg-red-500",
    info: "bg-indigo-500",
  };

  const icons = {
    success:
      '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>',
    error:
      '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>',
    info: '<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>',
  };
</script>

<div
  in:fly={{ y: 20, duration: 300 }}
  out:fade
  class="{colors[
    type
  ]} text-white px-6 py-4 rounded-lg shadow-lg flex items-center space-x-3 min-w-[300px] pointer-events-auto"
  role="alert"
>
  <span>{@html icons[type]}</span>
  <p class="font-medium">{message}</p>
</div>
