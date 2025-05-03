<script>
	import { onMount } from 'svelte';
	let config = {};

	onMount(async () => {
		const res = await fetch('/api/config');
		config = await res.json();
	});

	async function saveConfig() {
		await fetch('/api/config', {
			method: 'POST',
			body: JSON.stringify(config),
			headers: {
				'Content-Type': 'application/json'
			}
		});
		alert('Saved!');
	}
</script>

<form on:submit|preventDefault={saveConfig} class="space-y-4 p-4">
	{#each Object.entries(config) as [key, value]}
		<div>
			<label from={key} class="block font-medium">{key}</label>
			<input id={key} class="w-full border p-2" bind:value={config[key]} />
		</div>
	{/each}
	<button class="rounded bg-blue-500 px-4 py-2 text-white">Save</button>
</form>
