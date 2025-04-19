<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import LogTable from '$lib/components/LogTable.svelte';
	import { fade } from 'svelte/transition';
	import { queryHistory, addQuery } from '$lib/queryHistory';
	import QueryDrawer from '$lib/components/QueryDrawer.svelte';

	let drawerOpen = false;
	let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10';
	let results: any[] = [];
	let error: string | null = null;
	let success = false;
	let durationMs: number | null = null;

	async function fetchQuery() {
		error = null;
		success = false;
		durationMs = null;

		const start = performance.now();

		try {
			const res = await fetch(`/query?q=${encodeURIComponent(query)}`);
			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Unknown error');
			}
			results = await res.json().then((r) => r?.data || []);
			addQuery({ query, ok: true, timestamp: Date.now() });
			durationMs = performance.now() - start;
			success = true;
			setTimeout(() => (success = false), 2500);
		} catch (err) {
			error = err.message;
			addQuery({ query, ok: false, timestamp: Date.now() });
		}
	}

	let initialVisibility = {
		timestamp: true,
		trace_id: true,
		level: true,
		message: true,
		raw: true
	};
</script>

<div class="mx-auto max-w-screen-xl space-y-4 p-4">
	<h2 class="text-xl font-bold">{m.query_logs()}</h2>

	<textarea
		bind:value={query}
		class="w-full border border-gray-600 bg-gray-800 p-2 font-mono text-sm"
		rows={4}
	></textarea>

	<div class="flex items-center gap-2">
		<button
			on:click={fetchQuery}
			class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-500"
		>
			{m.run_query()}
		</button>

		<button
			on:click={() => (drawerOpen = true)}
			class="ml-auto rounded bg-gray-700 px-3 py-1 text-sm hover:bg-gray-600"
		>
			{m.bad_kind_flamingo_nurture()}
			{#if $queryHistory.length > 0}
				({$queryHistory.length})
			{/if}
		</button>
	</div>

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
			{#if durationMs !== null}
				<span class="ml-2 text-xs text-gray-300">({Math.round(durationMs)}ms)</span>
			{/if}
		</div>
	{/if}

	<QueryDrawer
		bind:open={drawerOpen}
		onSelect={(q) => {
			query = q;
			drawerOpen = false;
		}}
	/>

	<LogTable logs={results} {initialVisibility} />
</div>
