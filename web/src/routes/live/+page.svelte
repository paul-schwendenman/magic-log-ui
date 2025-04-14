<script lang="ts">
	import LogTable from '$lib/components/LogTable.svelte';
	import type { LogEntry } from '$lib/types';
	import { onMount } from 'svelte';

	let logs: LogEntry[] = [];

	const maxLogs = 500;

    onMount(() => {
	    const ws = new WebSocket(`ws://${location.host}/ws`);

        ws.addEventListener('message', (event) => {
            const entry: LogEntry = JSON.parse(event.data);
            logs = [entry, ...logs].slice(0, maxLogs);
        });

        ws.addEventListener('close', () => {
            console.warn('WebSocket closed');
        });

        ws.addEventListener('error', (err) => {
            console.error('WebSocket error:', err);
        });

        return () => ws.close();
    });
</script>

<h2 class="mb-2 text-xl font-bold">Live Logs</h2>
<LogTable {logs} />
