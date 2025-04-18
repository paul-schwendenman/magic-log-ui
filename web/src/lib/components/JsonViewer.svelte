<script lang="ts">
	import { onMount } from 'svelte';

	export let data: unknown;
	export let level: number = 0;
	export let collapsed: boolean = false;

	let isCollapsed = collapsed;
	let keys: string[] = [];

	onMount(() => {
		if (typeof data === 'object' && data !== null && !Array.isArray(data)) {
			keys = Object.keys(data as Record<string, unknown>);
		}
	});

	function toggle() {
		isCollapsed = !isCollapsed;
	}

	const indent = `pl-${Math.min(level * 4, 48)}`; // tailwind safe cap
</script>

{#if typeof data === 'object' && data !== null}
	<div class={`font-mono text-sm ${indent}`}>
		{#if Array.isArray(data)}
			<div>
				<button on:click={toggle} class="text-gray-400 select-none">
					{isCollapsed ? '▶' : '▼'} [
				</button>
				{#if !isCollapsed}
					{#each data as item, i}
						<svelte:self data={item} level={level + 1} />
					{/each}
				{/if}
				<span class="text-gray-400">]</span>
			</div>
		{:else}
			<div>
				<button on:click={toggle} class="text-gray-400 select-none">
					{isCollapsed ? '▶' : '▼'} {'{'}
				</button>
				{#if !isCollapsed}
					{#each keys as key}
						<div class="flex gap-1">
							<span class="text-teal-500">"{key}"</span>:
							<svelte:self data={(data as any)[key]} level={level + 1} />
						</div>
					{/each}
				{/if}
				<span class="text-gray-400">{'}'}</span>
			</div>
		{/if}
	</div>
{:else}
	<span class={`font-mono text-sm ${indent}`}>
		{#if typeof data === 'string'}
			<span class="text-green-600">"{data}"</span>
		{:else if typeof data === 'number'}
			<span class="text-blue-500">{data}</span>
		{:else if typeof data === 'boolean'}
			<span class="text-rose-600">{data ? 'true' : 'false'}</span>
		{:else if data === null}
			<span class="text-gray-500 italic">null</span>
		{:else}
			<span class="text-gray-400">{String(data)}</span>
		{/if}
	</span>
{/if}
