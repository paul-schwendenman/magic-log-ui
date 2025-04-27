<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import { fade } from 'svelte/transition';
	import LogLineSimple from './LogLineSimple.svelte';

	let query =
		'SELECT level, message, trace_id, timestamp FROM logs ORDER BY timestamp DESC LIMIT 10';
	let results: any[] = [];
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
			results = await res.json().then((resp) => resp?.data || []);
			setTimeout(() => {
				success = true;
				setTimeout(() => (success = false), 2500);
			}, 100);
		} catch (err) {
			error = err.message;
		}
	}
</script>

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
	<div class="space-y-2 lg:max-h-[75vh] lg:overflow-y-auto">
		{#each results as log (log)}
			<LogLineSimple {log} />
		{/each}
	</div>
</div>
