<script lang="ts">
	import LogTable from '$lib/components/LogTable.svelte';
	import { paused, liveFilter, filteredLiveLogs, buffer, clearLogs } from '$lib/stores/liveLogs';
	const bufferSize = $derived($buffer.length);

	function togglePause() {
		paused.update((p) => {
			// if (p) resumeLogs();
			return !p;
		});
	}

	let initialVisibility = {
		timestamp: true,
		trace_id: true,
		level: true,
		message: true,
		raw: false
	};
</script>

<h2 class="my-2 text-xl font-bold">Live Logs</h2>
<input
	bind:value={$liveFilter}
	placeholder="Filter logs..."
	class="mb-4 w-full rounded border border-gray-600 bg-gray-800 p-2 text-sm"
/>

<div class="mb-4 flex items-center gap-2">
	<button
		onclick={togglePause}
		class="rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
	>
		{$paused ? 'Resume' : 'Pause'}
	</button>

	<button
		onclick={clearLogs}
		class="rounded bg-red-700 px-4 py-1 text-sm text-white hover:bg-red-600"
	>
		Clear
	</button>
	{#if $paused && bufferSize}
		<span class="text-xs text-yellow-400">+{bufferSize} buffered</span>
	{/if}
</div>

<LogTable logs={$filteredLiveLogs} {initialVisibility} />
