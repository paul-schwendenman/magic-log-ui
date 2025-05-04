<script lang="ts">
	import PresetEditor from '$lib/components/PresetEditor.svelte';
	import { onMount } from 'svelte';

	let config: any = null;
	let loading = true;
	let saving = false;

	let errorMessages: string[] = [];
	let successMessage = '';
	let successTimeout: ReturnType<typeof setTimeout> | null = null;

	const defaultsFieldTypes: Record<string, 'string' | 'number' | 'boolean' | 'preset'> = {
		db_file: 'string',
		port: 'number',
		launch: 'boolean',
		log_format: 'string',
		regex_preset: 'preset',
		regex: 'string',
		jq_filter: 'string',
		jq_preset: 'preset'
	};

	$: regexPresetOptions = config?.regex_presets ? Object.keys(config.regex_presets) : [];
	$: jqPresetOptions = config?.jq_presets ? Object.keys(config.jq_presets) : [];

	onMount(async () => {
		try {
			const res = await fetch('/api/config');
			config = await res.json();
		} catch (e) {
			errorMessages = ['Failed to load config.'];
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		errorMessages = [];
		successMessage = '';
		if (successTimeout) clearTimeout(successTimeout);

		try {
			const res = await fetch('/api/config', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(config)
			});

			if (!res.ok) {
				const contentType = res.headers.get('Content-Type');
				if (contentType?.includes('application/json')) {
					const data = await res.json();
					errorMessages = Array.isArray(data.errors) ? data.errors : [JSON.stringify(data)];
				} else {
					const text = await res.text();
					errorMessages = [text || 'Unknown error'];
				}
				return;
			}

			successMessage = '✅ Config saved successfully!';
			successTimeout = setTimeout(() => {
				successMessage = '';
			}, 4000);
		} catch (e) {
			errorMessages = ['Failed to save config.'];
		} finally {
			saving = false;
		}
	}

	function validateRegex(value: string): string | null {
		try {
			new RegExp(value);
			return null;
		} catch {
			return 'Invalid regular expression.';
		}
	}

	function validateJQ(value: string): string | null {
		if (!value.trim()) return 'JQ filter cannot be empty.';
		// For real JQ parsing/validation, do it on backend
		return null;
	}

	function coerceDefault(key: string, val: string): any {
		if (key === 'port') return parseInt(val, 10);
		if (key === 'launch') return val === 'true';
		return val;
	}
</script>

{#if loading}
	<p class="p-4">Loading...</p>
{:else}
	<div class="mx-auto max-w-3xl space-y-6 p-4">
		<h1 class="text-2xl font-bold">App Settings</h1>

		{#if errorMessages.length}
			<ul class="space-y-1 rounded bg-red-100 p-3 text-red-700">
				{#each errorMessages as msg}
					<li>• {msg}</li>
				{/each}
			</ul>
		{/if}

		{#if successMessage}
			<div class="rounded bg-green-100 p-3 text-green-700">{successMessage}</div>
		{/if}

		<!-- Defaults -->
		{#if config?.defaults}
			<section>
				<h2 class="mb-2 text-xl font-semibold">Defaults</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					{#each Object.entries(config.defaults) as [key, value]}
						<div>
							<label for={key} class="block font-medium capitalize">{key}</label>
							{#if defaultsFieldTypes[key] === 'boolean'}
								<select
									id={key}
									class="w-full rounded border p-2"
									bind:value={config.defaults[key]}
								>
									<option value={true}>true</option>
									<option value={false}>false</option>
								</select>
							{:else if defaultsFieldTypes[key] === 'number'}
								<input
									id={key}
									type="number"
									class="w-full rounded border p-2"
									bind:value={config.defaults[key]}
									on:input={(e) => (config.defaults[key] = +e.target.value)}
								/>
							{:else if defaultsFieldTypes[key] === 'preset'}
								<select
									id={key}
									class="w-full rounded border p-2"
									bind:value={config.defaults[key]}
								>
									<option value="">-- select --</option>
									{#each key === 'regex_preset' ? regexPresetOptions : jqPresetOptions as option}
										<option value={option}>{option}</option>
									{/each}
								</select>
							{:else}
								<input
									id={key}
									class="w-full rounded border p-2"
									bind:value={config.defaults[key]}
								/>
							{/if}
						</div>
					{/each}
				</div>
			</section>
		{/if}

		<!-- Regex Presets -->
		{#if config?.regex_presets}
			<PresetEditor
				title="Regex Presets"
				presets={config.regex_presets}
				validateValue={validateRegex}
			/>
		{/if}

		<!-- JQ Presets -->
		{#if config?.jq_presets}
			<PresetEditor title="JQ Presets" presets={config.jq_presets} validateValue={validateJQ} />
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
