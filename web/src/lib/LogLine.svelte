<script>
	export let log;

	let expanded = false;

	function toggle() {
		expanded = !expanded;
	}
</script>

<pre
	class="cursor-pointer rounded bg-gray-800 p-2 font-mono text-sm whitespace-pre-wrap transition-all hover:bg-gray-700"
	on:click={toggle}>
    {#if expanded}
		{JSON.stringify(log, null, 2)}
	{:else}
		{JSON.stringify(log)}
	{/if}
  </pre>
