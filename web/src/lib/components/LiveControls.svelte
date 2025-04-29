<script lang="ts">
	import { m } from '$lib/paraglide/messages';
	import { paused, liveLogs } from '$lib/stores/liveLogs';
	const bufferSize = liveLogs.bufferSize;
	const isBufferFull = liveLogs.isBufferFull;

	function togglePause() {
		paused.update((p) => {
			return !p;
		});
	}
</script>

<div class="mb-4 flex flex-wrap items-center justify-between gap-4">
	<!-- Primary controls -->
	<div class="flex items-center gap-2">
		<button
			onclick={togglePause}
			class="rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
		>
			{$paused ? m.mild_good_skunk_buzz() : m.proof_spare_reindeer_praise()}
		</button>

		<button
			onclick={liveLogs.clearLogs}
			class="rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
		>
			{m.teary_nimble_crossbill_taste()}
		</button>

		{#if $paused}
			<button
				onclick={liveLogs.clearBuffer}
				class="rounded bg-yellow-700 px-4 py-1 text-sm text-white hover:bg-yellow-600"
			>
				{m.calm_ideal_bee_peek()}
			</button>
		{/if}
	</div>

	{#if $paused && $bufferSize > 0}
		<div
			class="rounded bg-yellow-900 px-2 py-1 text-xs text-yellow-400"
			class:font-bold={$isBufferFull}
		>
			{#if $isBufferFull}
				{m.plane_calm_oryx_earn()}
			{:else}
				{m.sweet_sound_toucan_bloom({ bufferSize: $bufferSize })}
			{/if}
		</div>
	{/if}
</div>
