<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { format } from 'date-fns';
	import { presets } from '$lib/time-range-presets';
	import type { TimeRangeConfig } from '$lib/types';
	import { m } from '$lib/paraglide/messages';

	export let value: TimeRangeConfig;
	export let onChange: (range: TimeRangeConfig) => void;

	let interval: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		if (value.live && value.durationMs) {
			startLive();
		}
	});

	onDestroy(() => {
		stopLive();
	});

	function selectRange(preset: TimeRangeConfig) {
		stopLive();
		const now = new Date();
		onChange({
			...preset,
			from: new Date(now.getTime() - (preset.durationMs ?? 0)),
			to: now,
			live: preset.live || value.live
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
		onChange({ ...value, from, to, relative: false, live: false });
	}

	function stepForward() {
		stopLive();
		const duration = value.to.getTime() - value.from.getTime();
		const from = new Date(value.to.getTime());
		const to = new Date(value.to.getTime() + duration);
		onChange({ ...value, from, to, relative: false, live: false });
	}

	function toggleLive() {
		stopLive();
		if (!value.live && value.durationMs) {
			startLive();
		}
		onChange({
			...value,
			live: !value.live
			// from: !value.live && value.durationMs ? new Date(Date.now() - value.durationMs) : value.from,
			// to: !value.live ? new Date() : value.to
		});
	}

	function reload() {
		if (value.relative) {
			const now = new Date();
			onChange({
				...value,
				from: new Date(now.getTime() - (value.durationMs ?? 0)),
				to: now
				// live: true
			});
		} else {
			onChange({
				...value
			});
		}
	}

	function startLive() {
		stopLive();
		interval = setInterval(() => {
			reload();
		}, value.refreshMs ?? 5000);
	}

	function stopLive() {
		if (interval) {
			clearInterval(interval);
			interval = null;
		}
	}
</script>

<div class="flex gap-4">
	<div class="flex items-center gap-2">
		<div class="relative">
			<select
				bind:value={value.label}
				on:change={(e) => {
					const preset = presets.find((p) => p.label === e.target.value);
					if (preset) selectRange(preset);
				}}
				class="rounded border bg-gray-800 px-2 py-1 text-sm"
			>
				{#if !value.relative}
					<option value={value.label}>
						{format(value.from, 'MMM d, h:mm a')} – {format(value.to, 'MMM d, h:mm a')}
					</option>
				{/if}

				{#each presets as preset}
					<option value={preset.label}>{preset.label}</option>
				{/each}
			</select>
		</div>

		{#if value.live}
			<span class="rounded bg-green-700 px-2 py-1 text-xs font-semibold text-white"
				>{m.bland_any_dingo_empower()}</span
			>
		{/if}
	</div>

	<div class="flex items-center gap-2">
		<button on:click={stepBack} class="rounded bg-gray-700 px-2 py-1">«</button>
		<button on:click={toggleLive} class="rounded bg-gray-700 px-2 py-1">
			{value.live ? '⏸' : '⏯'}
		</button>
		<button on:click={stepForward} class="rounded bg-gray-700 px-2 py-1">»</button>
		<button on:click={reload} class="rounded bg-gray-700 px-2 py-1">↻</button>
	</div>
</div>
