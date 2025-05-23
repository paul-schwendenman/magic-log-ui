<script lang="ts">
	export let title: string;
	export let presets: Record<string, string>;
	export let validateValue: (value: string) => string | null;

	let newKey = '';
	let newValue = '';
	let error = '';

	function addPreset() {
		error = '';
		if (!newKey.trim()) {
			error = 'Preset name is required.';
			return;
		}
		if (newKey in presets) {
			error = 'Preset name already exists.';
			return;
		}
		const validationError = validateValue?.(newValue);
		if (validationError) {
			error = validationError;
			return;
		}
		presets[newKey] = newValue;
		newKey = '';
		newValue = '';
	}

	function removePreset(key: string) {
		const updated = { ...presets };
		delete updated[key];
		presets = updated;
	}
</script>

<section class="mb-6">
	<h2 class="mb-2 text-xl font-semibold">{title}</h2>

	{#if error}
		<div class="mb-2 rounded bg-red-100 p-2 text-red-700">{error}</div>
	{/if}

	<!-- Existing presets -->
	{#each Object.entries(presets) as [key, value]}
		<div class="mb-2 flex items-center gap-2">
			<input class="w-1/3 rounded border p-2" value={key} disabled />
			<input class="flex-1 rounded border p-2" bind:value={presets[key]} />
			<button
				class="rounded bg-red-500 p-2 text-white"
				aria-label="trash"
				on:click={() => removePreset(key)}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0"
					/>
				</svg>
			</button>
		</div>
	{/each}

	<!-- Add new preset -->
	<div class="mb-4 flex gap-2">
		<input class="w-1/3 rounded border p-2" placeholder="Name" bind:value={newKey} />
		<input class="flex-1 rounded border p-2" placeholder="Value" bind:value={newValue} />
		<button
			class="rounded bg-green-500 p-2 text-white disabled:opacity-50"
			on:click={addPreset}
			disabled={!newKey || newKey in presets}
			aria-label="Add"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="size-6"
			>
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
			</svg>
		</button>
	</div>
</section>
