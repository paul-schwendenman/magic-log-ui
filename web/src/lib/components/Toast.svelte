<script lang="ts">
	import { fly } from 'svelte/transition';
	export let text: string;
	export let duration = 1500;
	let visible = true;

	setTimeout(() => (visible = false), duration);
</script>

{#if visible}
	<div
		class="rounded bg-black px-2 py-1 text-xs text-white shadow"
		transition:fly={{ y: -10, duration: 300 }}
	>
		{text}
	</div>
{/if}
