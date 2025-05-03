<script lang="ts">
	import PresetEditor from '$lib/components/PresetEditor.svelte';
	import { onMount } from 'svelte';
	let config: any = null;
	let loading = true;
	let saving = false;
	let error = '';

	onMount(async () => {
		try {
			const res = await fetch('/api/config');
			config = await res.json();
		} catch (e) {
			error = 'Failed to load config.';
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		error = '';
		try {
			const res = await fetch('/api/config', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(config)
			});
			if (!res.ok) throw new Error('Save failed');
		} catch (e) {
			error = 'Failed to save config.';
		} finally {
			saving = false;
		}
	}
</script>

{#if loading}
	<p class="p-4">Loading...</p>
{:else}
	<div class="mx-auto max-w-3xl space-y-6 p-4">
		<h1 class="text-2xl font-bold">App Settings</h1>

		{#if error}
			<div class="rounded bg-red-100 p-2 text-red-800">{error}</div>
		{/if}

		<!-- Defaults -->
		{#if config?.defaults}
			<section>
				<h2 class="mb-2 text-xl font-semibold">Defaults</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					{#each Object.entries(config.defaults) as [key, value]}
						<div>
							<label class="block font-medium capitalize">{key}</label>
							<input class="w-full rounded border p-2" bind:value={config.defaults[key]} />
						</div>
					{/each}
				</div>
			</section>
		{/if}

		<!-- Regex Presets -->
		{#if config?.regex_presets}
			<PresetEditor title="Regex Presets" presets={config.regex_presets} />
		{/if}

		<!-- JQ Presets -->
		{#if config?.jq_presets}
			<PresetEditor title="JQ Presets" presets={config.jq_presets} />
		{/if}

		<button
			on:click={save}
			class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
			disabled={saving}
		>
			{saving ? 'Saving...' : 'Save Changes'}
		</button>
	</div>
{/if}
