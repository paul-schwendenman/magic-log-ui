<script lang="ts">
	import LiveControls from '$lib/components/LiveControls.svelte';
	import LogTable from '$lib/components/LogTable.svelte';
	import { liveFilter, filteredLiveLogs } from '$lib/stores/liveLogs';

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

<LiveControls />

<LogTable logs={$filteredLiveLogs} {initialVisibility} />
