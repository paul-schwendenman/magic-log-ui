<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import { isConnected } from '$lib/useWebSocket';
	import { liveLogs, paused } from '$lib/stores/liveLogs';
	import LogLineSimple from './LogLineSimple.svelte';
	import { createPaginationStore } from '$lib/stores/paginatedStore';

	const logs = createPaginationStore(liveLogs, 10);
</script>

<div>
	<h2 class="my-4 flex text-xl font-bold">
		<span>
			{m.live_logs()}
		</span>
		{#if !$isConnected}
			<div class="mb-2 ml-auto inline-block rounded bg-red-800/50 px-3 py-1 text-sm text-red-200">
				{m.known_awake_rooster_mend()}
			</div>
		{:else if $paused}
			<div
				class="mb-2 ml-auto inline-block rounded bg-yellow-900 px-3 py-1 text-sm text-yellow-300"
			>
				Stream is paused
			</div>
		{/if}
	</h2>
	<div class="space-y-2 lg:max-h-[90vh] lg:overflow-y-auto">
		{#each $logs.items as log (log)}
			<LogLineSimple {log} />
		{/each}
	</div>
</div>
