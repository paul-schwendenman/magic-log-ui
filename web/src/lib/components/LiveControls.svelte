<script lang="ts">
	import { paused, buffer, clearLogs, clearBuffer, bufferFull } from '$lib/stores/liveLogs';

	const bufferSize = $derived($buffer.length);

	function togglePause() {
		paused.update((p) => {
			return !p;
		});
	}
</script>

<div class="mb-4 flex items-center gap-2">
	<button
		onclick={togglePause}
		class="rounded bg-gray-700 px-4 py-1 text-sm text-white hover:bg-gray-600"
	>
		{$paused ? 'Resume' : 'Pause'}
	</button>

	<button
		onclick={clearLogs}
		class="rounded bg-red-700 px-4 py-1 text-sm text-white hover:bg-red-600"
	>
		Clear
	</button>
	{#if $paused && bufferSize}
		<button
			onclick={clearBuffer}
			class="rounded bg-red-700 px-4 py-1 text-sm text-white hover:bg-red-600"
		>
			Flush
		</button>
		<span class="text-xs text-yellow-400">+{bufferSize} buffered</span>
	{/if}
</div>
