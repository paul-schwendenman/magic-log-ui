<script lang="ts">
	import LogTable from '$lib/components/LogTable.svelte';
	import { paused, liveFilter, filteredLiveLogs, buffer } from '$lib/stores/liveLogs';

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
<button
	on:click={togglePause}
	class="mb-4 rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
>
	{$paused ? 'Resume Live Logs' : 'Pause'}
</button>
<LogTable logs={$filteredLiveLogs} {initialVisibility} />
