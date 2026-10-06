<script lang="ts">
	let className = 'h-7 w-7';
	export { className as class };
</script>

<svg class={className} viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
	<defs>
		<linearGradient id="logo-g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
			<stop stop-color="#6366f1" />
			<stop offset="1" stop-color="#a855f7" />
		</linearGradient>
	</defs>
	<rect width="32" height="32" rx="9" fill="url(#logo-g)" />
	<path d="M11 14.5V12a5 5 0 0 1 10 0v2.5" stroke="white" stroke-width="2.2" stroke-linecap="round" />
	<rect x="8.5" y="14.5" width="15" height="10.5" rx="2.5" fill="white" />
	<circle cx="16" cy="19.2" r="1.7" fill="#6366f1" />
	<path d="M16 20.5v2" stroke="#6366f1" stroke-width="1.6" stroke-linecap="round" />
</svg>
