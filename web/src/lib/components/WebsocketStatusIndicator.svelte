<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import { wsStatus } from '$lib/useWebSocket';
	import { paused } from '$lib/stores/liveLogs';
</script>

{#if $wsStatus === 'connecting'}
	<div class="mb-2 ml-auto inline-block rounded bg-blue-800/50 px-3 py-1 text-sm text-blue-200">
		{m.zesty_same_kitten_view()}
	</div>
{:else if $wsStatus === 'error'}
	<div class="mb-2 ml-auto inline-block rounded bg-red-900 px-3 py-1 text-sm text-red-300">
		{m.jumpy_elegant_scallop_cherish()}
	</div>
{:else if $wsStatus === 'closed'}
	<div class="mb-2 ml-auto inline-block rounded bg-red-800/50 px-3 py-1 text-sm text-red-200">
		{m.known_awake_rooster_mend()}
	</div>
{:else if $paused}
	<div class="mb-2 ml-auto inline-block rounded bg-yellow-900 px-3 py-1 text-sm text-yellow-300">
		{m.novel_cozy_goose_stab()}
	</div>
{/if}
