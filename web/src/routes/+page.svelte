<script lang="ts">
	import { onMount } from 'svelte';
	import { m } from '$lib/paraglide/messages.js';

	let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10';
	let results: any[] = [];
	let wsLogs: any[] = [];
	let socket: WebSocket;

	const fetchQuery = async () => {
		const res = await fetch(`/query?q=${encodeURIComponent(query)}`);
		results = await res.json();
	};

	onMount(() => {
		socket = new WebSocket(`ws://${location.host}/ws`);
		socket.onmessage = (msg) => {
			wsLogs = [JSON.parse(msg.data), ...wsLogs].slice(0, 100);
		};
	});
</script>

<div class="space-y-4">
	<input bind:value={query} class="w-full rounded border border-gray-600 bg-gray-800 p-2" />
	<button on:click={fetchQuery} class="rounded bg-blue-600 px-4 py-2 hover:bg-blue-500"
		>{m.run_query()}</button
	>

	<h2 class="text-xl font-bold">{m.live_logs()}</h2>
	{#each wsLogs as log}
		<pre class="mb-1 rounded bg-gray-800 p-2 text-sm">{JSON.stringify(log, null, 2)}</pre>
	{/each}

	<h2 class="mt-6 text-xl font-bold">{m.query_results()}</h2>
	{#each results as row}
		<pre class="mb-1 rounded bg-gray-700 p-2 text-sm">{JSON.stringify(row, null, 2)}</pre>
	{/each}
</div>
