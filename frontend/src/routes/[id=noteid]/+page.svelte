<script lang="ts">
	import { page } from '$app/stores';
	import NoteViewer from '$lib/components/NoteViewer.svelte';
</script>

<NoteViewer id={$page.params.id ?? ''} />
