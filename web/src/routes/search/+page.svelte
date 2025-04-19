<script lang="ts">
	import QueryInput from '../../lib/components/QueryInput.svelte';

	import { m } from '$lib/paraglide/messages.js';
	import { fade } from 'svelte/transition';
	import { queryHistory, addQuery } from '$lib/queryHistory';
	import QueryDrawer from '$lib/components/QueryDrawer.svelte';
	import LogLine from '$lib/components/LogLine.svelte';
	import { createQueryStore } from '$lib/stores/queryStore';
	import { onMount } from 'svelte';

	const initialLimit = 10;
	const initialQuery = 'SELECT * FROM logs ORDER BY timestamp DESC';

	let drawerOpen = $state(false);
	let query = $state(initialQuery);
	let showSuccess = $state(false);
	let limit = $state(initialLimit);

	const store = createQueryStore({ query: initialQuery, limit: initialLimit });
	const page = $derived($store.page);

	onMount(() =>
		store.subscribe((state) => {
			if (state.error) {
				addQuery({ query, ok: false, timestamp: Date.now() });
			} else {
				showSuccess = true;
				setTimeout(() => (showSuccess = false), 2500);
				addQuery({ query, ok: true, timestamp: Date.now() });
			}
		})
	);

	function fetchQuery() {
		return store.setQuery(query);
	}
</script>

<div class="mx-auto max-w-screen-xl space-y-4 p-4">
	<QueryInput bind:query onQuery={fetchQuery} />

	<div class="flex items-center gap-2">
		<div class="flex gap-2">
			<button
				onclick={fetchQuery}
				class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-500"
			>
				{m.run_query()}
			</button>
			{#if $store.error}
				<div
					class="rounded border border-red-500 bg-red-900/30 p-2 text-sm text-red-400"
					in:fade={{ duration: 300 }}
					out:fade={{ duration: 100 }}
				>
					<b>{m.lower_noisy_eel_treasure()}:</b>
					{$store.error}
				</div>
			{:else if showSuccess}
				<div
					class="rounded border border-green-500 bg-green-900/30 p-2 text-sm text-green-400"
					in:fade={{ duration: 300 }}
					out:fade={{ duration: 200 }}
				>
					{m.mellow_many_quail_imagine()}
					{#if $store.durationMs !== null}
						<span class="ml-2 text-xs text-gray-300">({Math.round($store.durationMs || 0)}ms)</span>
					{/if}
				</div>
			{/if}
		</div>

		<button
			onclick={() => (drawerOpen = true)}
			class="ml-auto rounded bg-gray-700 px-3 py-1 text-sm hover:bg-gray-600"
		>
			{m.bad_kind_flamingo_nurture()}
			{#if $queryHistory.length > 0}
				({$queryHistory.length})
			{/if}
		</button>
	</div>

	<QueryDrawer
		bind:open={drawerOpen}
		onSelect={(q) => {
			query = q;
			drawerOpen = false;
		}}
	/>

	<div class="space-y-2 lg:max-h-[75vh] lg:overflow-y-auto">
		{#each $store.results as log (log)}
			<LogLine {log} />
		{/each}
	</div>
	<div
		class="mt-4 flex flex-wrap items-center justify-between gap-4 border-t pt-4 text-sm text-gray-300"
	>
		<div class="flex items-center gap-2">
			<label for="limit" class="text-gray-400">Rows per page:</label>
			<select
				id="limit"
				class="rounded border border-gray-600 bg-gray-800 p-1"
				bind:value={limit}
				onchange={() => store.setLimit(limit)}
			>
				<option value={10}>10</option>
				<option value={20}>20</option>
				<option value={50}>50</option>
				<option value={100}>100</option>
			</select>
		</div>

		<div class="flex items-center gap-2">
			<button
				onclick={store.prevPage}
				disabled={!$store.meta.hasPreviousPage}
				class="rounded bg-gray-700 px-3 py-1 hover:bg-gray-600 disabled:opacity-50"
			>
				Previous
			</button>

			<span class="text-xs">Page {page + 1}</span>

			<button
				onclick={store.nextPage}
				disabled={!$store.meta.hasNextPage}
				class="rounded bg-gray-700 px-3 py-1 hover:bg-gray-600 disabled:opacity-50"
			>
				Next
			</button>
		</div>
	</div>
</div>
