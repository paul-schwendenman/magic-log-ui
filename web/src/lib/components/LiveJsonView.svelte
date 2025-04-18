<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import LogLine from '$lib/components/LogLine.svelte';
	import { isConnected } from '$lib/useWebSocket';
	import { liveLogs } from '$lib/stores/liveLogs';
</script>

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
