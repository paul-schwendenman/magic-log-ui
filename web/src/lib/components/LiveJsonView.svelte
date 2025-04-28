<script lang="ts">
	import WebsocketStatusIndicator from './WebsocketStatusIndicator.svelte';

	import { m } from '$lib/paraglide/messages.js';
	import { liveLogs } from '$lib/stores/liveLogs';
	import LogLineSimple from './LogLineSimple.svelte';
	import { createPaginationStore } from '$lib/stores/paginatedStore';

	const logs = createPaginationStore(liveLogs, 10);
</script>

<div>
	<h2 class="my-4 flex text-xl font-bold">
		<span>
			{m.live_logs()}
		</span>
		<WebsocketStatusIndicator />
	</h2>
	<div class="space-y-2 lg:max-h-[90vh] lg:overflow-y-auto">
		{#each $logs.items as log (log)}
			<LogLineSimple {log} />
		{/each}
	</div>
</div>
