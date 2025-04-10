<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<main class="p-4 bg-gray-900 text-white min-h-screen font-mono">
	{@render children()}
</main>
