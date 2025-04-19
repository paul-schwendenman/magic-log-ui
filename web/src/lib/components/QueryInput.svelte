<script lang="ts">
	export let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10';
	export let onQuery = (query: string) => {};
</script>

<textarea
	bind:value={query}
	on:keydown={(e) => {
		if (e.key === 'Enter' && e.shiftKey) {
			e.preventDefault();
			onQuery(query);
		} else if (e.key === 'Tab') {
			e.preventDefault();

			const target = e.target as HTMLTextAreaElement;
			const start = target.selectionStart;
			const end = target.selectionEnd;
			const SPACES = '   ';

			query = query.slice(0, start) + SPACES + query.slice(end);

			requestAnimationFrame(() => {
				target.selectionStart = target.selectionEnd = start + SPACES.length;
			});
		}
	}}
	class="w-full border border-gray-600 bg-gray-800 p-2 font-mono text-sm"
	rows={4}
></textarea>
