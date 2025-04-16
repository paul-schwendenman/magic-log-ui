<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import LogLine from '$lib/LogLine.svelte';
	import { fade } from 'svelte/transition';
	import { isConnected } from '$lib/useWebSocket';
	import { queryHistory, addQuery } from '$lib/queryHistory';
	import QueryDrawer from '$lib/QueryDrawer.svelte';
	import { liveLogs } from '$lib/stores/liveLogs';

	let drawerOpen = false;
	let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10';
	let results: any[] = [];
	let wsLogs: any[] = [];
	let error: string | null = null;
	let success = false;

	async function fetchQuery() {
		error = null;

		try {
			const res = await fetch(`/query?q=${encodeURIComponent(query)}`);
			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Unknown error');
			}
			results = await res.json();
			addQuery({ query, ok: true, timestamp: Date.now() });
			setTimeout(() => {
				success = true;
				setTimeout(() => (success = false), 2500);
			}, 100);
		} catch (err) {
			error = err.message;
			addQuery({ query, ok: false, timestamp: Date.now() });
		}
	}
</script>

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	<div class="space-y-4">
		<h2 class="my-4 text-xl font-bold">{m.query_logs()}</h2>
		<input
			bind:value={query}
			class="w-full rounded border border-gray-600 bg-gray-800 p-2"
			placeholder="SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100"
		/>
		<div class="space-between flex gap-1">
			<button on:click={fetchQuery} class="h-fit rounded bg-blue-600 px-4 py-2 hover:bg-blue-500">
				{m.run_query()}
			</button>
			<button
				on:click={() => (drawerOpen = true)}
				class="fixed top-4 right-4 z-50 rounded bg-gray-700 px-3 py-1 text-sm hover:bg-gray-600"
			>
				{m.bad_kind_flamingo_nurture()}
				{#if $queryHistory.length > 0}
					({$queryHistory.length})
				{/if}
			</button>

			{#if error}
				<div
					class="rounded border border-red-500 bg-red-900/30 p-2 text-sm text-red-400"
					in:fade={{ duration: 300 }}
					out:fade={{ duration: 100 }}
				>
					<b>{m.lower_noisy_eel_treasure()}:</b>
					{error}
				</div>
			{:else if success}
				<div
					class="rounded border border-green-500 bg-green-900/30 p-2 text-sm text-green-400"
					in:fade={{ duration: 300 }}
					out:fade={{ duration: 200 }}
				>
					{m.mellow_many_quail_imagine()}
				</div>
			{/if}
		</div>
		<QueryDrawer
			bind:open={drawerOpen}
			onSelect={(q) => {
				query = q;
				drawerOpen = false;
			}}
		/>

		<div class="space-y-2 lg:max-h-[75vh] lg:overflow-y-auto">
			{#each results as log (log)}
				<LogLine {log} />
			{/each}
		</div>
	</div>

	<div>
		<h2 class="my-4 text-xl font-bold">{m.live_logs()}</h2>
		{#if !$isConnected}
			<div class="mb-2 inline-block rounded bg-red-800/50 px-3 py-1 text-sm text-red-200">
				{m.known_awake_rooster_mend()}
			</div>
		{/if}
		<div class="space-y-2 lg:max-h-[90vh] lg:overflow-y-auto">
			{#each $liveLogs as log (log)}
				<LogLine {log} />
			{/each}
		</div>
	</div>
</div>
