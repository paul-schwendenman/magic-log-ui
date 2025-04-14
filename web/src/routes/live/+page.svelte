<script lang="ts">
	import LogTable from '$lib/components/LogTable.svelte';
	import { liveLogs } from '$lib/stores/liveLogs';
</script>

<h2 class="mb-2 text-xl font-bold">Live Logs</h2>
<LogTable logs={$liveLogs} />
