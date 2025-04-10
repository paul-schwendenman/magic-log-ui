<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<main class="min-h-screen bg-gray-900 p-4 font-mono text-white">
	{@render children()}
</main>
