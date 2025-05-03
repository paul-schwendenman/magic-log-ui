<script lang="ts">
	import { page } from '$app/stores';
	import { m } from '$lib/paraglide/messages';

	const links = [
		{ href: '/', label: m.ok_neat_mallard_fade() },
		{ href: '/live', label: m.top_slow_shad_nudge() },
		{ href: '/search', label: m.best_silly_antelope_vent() },
		{ href: '/config', label: m.great_zesty_buzzard_savor() }
		// { href: '/table', label: 'Table (alpha)' }
		// add more pages here
	];
</script>

<nav class="flex items-center gap-6 border-b border-gray-700 bg-gray-900 px-6 py-2 text-sm">
	<a href="/" class="flex items-center gap-2 leading-none font-bold text-white">
		<img src="/logo.png" alt="Magic Log logo" class="h-8 w-auto align-middle" />
	</a>

	{#each links as { href, label }}
		<a {href} class="text-gray-300 hover:text-white" class:selected={$page.url.pathname === href}>
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
