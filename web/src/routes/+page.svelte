<script lang="ts">
	import LiveJsonView from '../lib/components/LiveJsonView.svelte';
	import JsonResultsView from '../lib/components/JsonResultsView.svelte';
</script>

<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
	<div>
		<JsonResultsView />
	</div>

	<div>
		<LiveJsonView></LiveJsonView>
	</div>
</div>
