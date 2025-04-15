<script>
	import { m } from '$lib/paraglide/messages.js';
	import { queryHistory, clearHistory } from '$lib/queryHistory';
	import { formatRelativeTime } from '$lib/formatRelative';
	import { fly, fade } from 'svelte/transition';
	export let open = false;
	export let onSelect = (q) => {};
</script>

{#if open}
	<div class="fixed inset-0 z-40 bg-black/30" on:click={() => (open = false)} in:fade out:fade />
	<aside
		class="fixed top-0 right-0 z-50 flex h-full w-72 flex-col border-l border-gray-700 bg-gray-900 p-4 shadow-lg"
		in:fly={{ x: 300 }}
		out:fly={{ x: 300 }}
	>
		<h2 class="my-2 text-lg font-bold text-white">Query History</h2>

		<div class="mb-3 flex justify-between text-sm text-gray-400">
			<span>{$queryHistory.length} queries</span>
			<button on:click={clearHistory} class="text-xs hover:underline">{m.clear()}</button>
		</div>

		<ul class="flex-1 space-y-1 overflow-y-auto pr-1 text-sm text-gray-300">
			{#each $queryHistory as entry}
				<li>
					<button
						class="flex w-full flex-col items-start gap-0.5 text-left hover:text-blue-400"
						on:click={() => onSelect(entry.query)}
					>
						<div class="flex w-full items-center gap-2 truncate">
							<span class={entry.ok ? 'text-green-400' : 'text-red-400'}>
								{entry.ok ? '✅' : '❌'}
							</span>
							<span class="truncate">{entry.query}</span>
						</div>
						<span class="ml-6 text-xs text-gray-500">
							{formatRelativeTime(entry.timestamp)}
						</span>
					</button>
				</li>
			{/each}
		</ul>
	</aside>
{/if}
