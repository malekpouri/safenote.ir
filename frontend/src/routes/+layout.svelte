<script lang="ts">
  import "../app.css";
  import Toast from "$lib/components/Toast.svelte"; // Changed from Toasts to Toast
  import { toasts, removeToast } from "$lib/stores/toast"; // Changed from { Toasts } to { toasts, removeToast }
  import { locale, dir, t } from "$lib/i18n"; // Added dir, t
  import { onMount } from "svelte";

  // Initialize currentDir based on the store
  let currentDir: string;
  const unsubscribeDir = dir.subscribe((d) => {
    currentDir = d;
  });

  onMount(() => {
    // Check for saved language preference or browser default
    const savedLocale = localStorage.getItem("locale");
    if (savedLocale === "fa" || savedLocale === "en") {
      locale.set(savedLocale);
    } else {
      // Default to Farsi if browser language is Farsi
      const browserLang = navigator.language;
      if (browserLang.startsWith("fa")) {
        locale.set("fa");
      }
    }

    // Save preference whenever locale changes
    const unsubscribeLocale = locale.subscribe((l) => {
      localStorage.setItem("locale", l);
    });

    return () => {
      unsubscribeDir();
      unsubscribeLocale();
    };
  });

  function toggleLanguage() {
    locale.toggle();
  }
</script>

<div
  class="min-h-screen flex flex-col bg-gradient-to-br from-slate-50 to-indigo-50"
  dir={currentDir}
>
  <div class="absolute top-4 right-8 z-40">
    <button
      on:click={toggleLanguage}
      class="bg-white/80 backdrop-blur-sm border border-slate-200 rounded-lg px-3 py-1.5 text-sm font-medium text-slate-600 hover:text-indigo-600 hover:border-indigo-300 transition-all shadow-sm"
    >
      {$locale === "en" ? "فارسی" : "English"}
    </button>
  </div>

  <slot />

  <!-- Toast Container -->
  <div
    class="fixed bottom-4 right-4 z-50 flex flex-col space-y-2 pointer-events-none"
    class:left-4={currentDir === "rtl"}
    class:right-auto={currentDir === "rtl"}
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
  :global([dir="rtl"] body) {
    font-family: "Vazirmatn", sans-serif;
  }
</style>
