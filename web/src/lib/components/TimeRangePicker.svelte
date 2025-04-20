<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { format } from 'date-fns';
	import { presets } from '$lib/time-range-presets';
	import type { TimeRange } from '$lib/types';

	export let value: TimeRange;
	export let onChange: (range: TimeRange) => void;

	let interval: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		if (value.live && value.durationMs) {
			startLive();
		}
	});

	onDestroy(() => {
		stopLive();
	});

	function selectRange(preset: TimeRange) {
		stopLive();
		const now = new Date();
		onChange({
			...preset,
			from: new Date(now.getTime() - (preset.durationMs ?? 0)),
			to: now
		});
		if (preset.live && preset.durationMs) {
			startLive();
		}
	}

	function stepBack() {
		stopLive();
		const duration = value.to.getTime() - value.from.getTime();
		const to = new Date(value.from.getTime());
		const from = new Date(value.from.getTime() - duration);
		onChange({ ...value, from, to, live: false });
	}

	function stepForward() {
		stopLive();
		const duration = value.to.getTime() - value.from.getTime();
		const from = new Date(value.to.getTime());
		const to = new Date(value.to.getTime() + duration);
		onChange({ ...value, from, to, live: false });
	}

	function startLive() {
		stopLive();
		interval = setInterval(() => {
			const now = new Date();
			onChange({
				...value,
				from: new Date(now.getTime() - (value.durationMs ?? 0)),
				to: now,
				live: true
			});
		}, 5000); // every 5 seconds
	}

	function stopLive() {
		if (interval) {
			clearInterval(interval);
			interval = null;
		}
	}
</script>

<div class="flex items-center gap-2">
	<button on:click={stepBack} class="rounded bg-gray-700 px-2 py-1">⏪</button>
	<div class="relative">
		<select
			bind:value={value.label}
			on:change={(e) => {
				const preset = presets.find((p) => p.label === e.target.value);
				if (preset) selectRange(preset);
			}}
			class="rounded border bg-gray-800 px-2 py-1 text-sm"
		>
			{#each presets as preset}
				<option value={preset.label}>{preset.label}</option>
			{/each}
		</select>
	</div>
	<span class="text-sm text-gray-400">
		{format(value.from, 'MMM d, h:mm a')} – {format(value.to, 'h:mm a')}
	</span>
	{#if value.live}
		<span class="rounded bg-green-700 px-2 py-1 text-xs font-semibold text-white">LIVE</span>
	{/if}
	<button on:click={stepForward} class="rounded bg-gray-700 px-2 py-1">⏩</button>
</div>
