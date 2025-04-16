<script lang="ts">
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
			{$paused ? 'Resume' : 'Pause'}
		</button>

		<button
			onclick={liveLogs.clearLogs}
			class="rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
		>
			Clear
		</button>

		{#if $paused }
			<button
				onclick={liveLogs.clearBuffer}
				class="rounded bg-yellow-700 px-4 py-1 text-sm text-white hover:bg-yellow-600"
			>
				Flush
			</button>
		{/if}
	</div>

	{#if $paused && $bufferSize > 0}
		<div
			class="rounded bg-yellow-900 px-2 py-1 text-xs text-yellow-400"
			class:font-bold={$isBufferFull}
		>
			{#if $isBufferFull}
				Buffer full
			{:else}
				+{$bufferSize} buffered
			{/if}
		</div>
	{/if}
</div>
