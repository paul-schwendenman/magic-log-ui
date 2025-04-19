<script lang="ts">
	import { page } from '$app/stores';
	import NavBar from '$lib/components/NavBar.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>{$page.data.title}</title>
	<meta name="title" content={$page.data.title} />

	<!-- Open Graph / Facebook -->
	<meta property="og:type" content="website" />
	<meta property="og:title" content={$page.data.title} />

	<!-- Twitter -->
	<meta property="twitter:card" content="summary" />
	<meta property="twitter:title" content={$page.data.title} />
</svelte:head>

<main class="min-h-screen bg-gray-900 p-4 font-mono text-white">
	<NavBar />
	{@render children()}
</main>
