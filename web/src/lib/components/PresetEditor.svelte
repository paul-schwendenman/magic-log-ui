<script lang="ts">
	export let title: string;
	export let presets: Record<string, string>;

	let newKey = '';
	let newValue = '';

	function addPreset() {
		if (!newKey.trim() || newKey in presets) return;
		presets[newKey] = newValue;
		newKey = '';
		newValue = '';
	}

	function removePreset(key: string) {
		delete presets[key];
	}
</script>

<section class="mb-6">
	<h2 class="mb-2 text-xl font-semibold">{title}</h2>

	<!-- Add preset -->
	<div class="mb-4 flex gap-2">
		<input class="w-1/3 rounded border p-2" placeholder="New preset name" bind:value={newKey} />
		<input class="flex-1 rounded border p-2" placeholder="Value" bind:value={newValue} />
		<button
			class="rounded bg-green-500 px-3 py-2 text-white disabled:opacity-50"
			on:click={addPreset}
			disabled={!newKey || newKey in presets}
		>
			Add
		</button>
	</div>

	<!-- Existing presets -->
	{#each Object.entries(presets) as [key, value]}
		<div class="mb-2 flex items-center gap-2">
			<input class="w-1/3 rounded border p-2" value={key} disabled />
			<input class="flex-1 rounded border p-2" bind:value={presets[key]} />
			<button class="rounded bg-red-500 px-2 py-1 text-white" on:click={() => removePreset(key)}>
				Remove
			</button>
		</div>
	{/each}
</section>
