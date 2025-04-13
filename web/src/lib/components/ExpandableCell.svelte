<script lang="ts">
	export let value: any;

	let expanded = false;

	$: formatted = expanded ? JSON.stringify(value, null, 2) : JSON.stringify(value);

	function toggle() {
		expanded = !expanded;
	}
</script>

<div class="cursor-pointer text-left whitespace-pre-wrap" on:click={toggle}>
	<pre class="max-w-full overflow-x-auto">{formatted}</pre>
</div>
