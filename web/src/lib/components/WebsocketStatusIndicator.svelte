<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import { wsStatus } from '$lib/useWebSocket';
	import { paused } from '$lib/stores/liveLogs';
	import { fade } from 'svelte/transition';

	export let showConnected = false;

	$: statusKey =
		$wsStatus === 'connecting'
			? 'connecting'
			: $wsStatus === 'error'
				? 'error'
				: $wsStatus === 'closed'
					? 'closed'
					: $wsStatus === 'open' && $paused
						? 'paused'
						: $wsStatus === 'open' && showConnected
							? 'connected'
							: null;

	$: statusClasses =
		{
			connecting: 'bg-blue-800/50 text-blue-200',
			error: 'bg-red-900 text-red-300',
			closed: 'bg-red-800/50 text-red-200',
			paused: 'bg-yellow-900 text-yellow-300',
			connected: 'bg-green-900/50 text-green-300'
		}[statusKey] ?? '';
</script>

{#if statusKey}
	{#key statusKey}
		<div
			in:fade={{ duration: 400 }}
			out:fade={{ duration: 400 }}
			class="mb-2 ml-auto inline-block rounded px-3 py-1 text-sm {statusClasses}"
		>
			{#if statusKey === 'connecting'}
				{m.zesty_same_kitten_view()}
			{:else if statusKey === 'error'}
				{m.jumpy_elegant_scallop_cherish()}
			{:else if statusKey === 'closed'}
				{m.known_awake_rooster_mend()}
			{:else if statusKey === 'paused'}
				{m.novel_cozy_goose_stab()}
			{:else if statusKey === 'connected'}
				{m.wise_gaudy_weasel_drum()}
			{/if}
		</div>
	{/key}
{/if}
