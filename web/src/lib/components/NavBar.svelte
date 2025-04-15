<script lang="ts">
	import { page } from '$app/stores';

	const links = [
		{ href: '/', label: 'Home' },
		{ href: '/live', label: 'Live' },
		// { href: '/table', label: 'Table' }
		// add more pages here
	];
</script>

<nav class="flex items-center gap-6 border-b border-gray-700 bg-gray-900 px-6 py-4 text-sm">
	{#each links as { href, label }}
		<a
			href={href}
			class="text-gray-300 hover:text-white"
			class:selected={$page.url.pathname === href}
		>
			{label}
		</a>
	{/each}
</nav>

<style>
	a.selected {
		font-weight: bold;
		color: white;
		border-bottom: 2px solid white;
	}
</style>
