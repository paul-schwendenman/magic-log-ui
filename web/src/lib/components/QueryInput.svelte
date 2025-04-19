<script lang="ts">
	export let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10';
	export let onQuery = () => {};
</script>

<textarea
	bind:value={query}
	on:keydown={(e) => {
		const target = e.target as HTMLTextAreaElement;

		if (e.key === 'Enter' && e.shiftKey) {
			e.preventDefault();
			onQuery();
			return;
		}

		if (e.key === 'Tab') {
			e.preventDefault();

			const start = target.selectionStart;
			const end = target.selectionEnd;

			if (e.shiftKey) {
				const before = query.slice(0, start);
				const after = query.slice(end);
				let removed = 0;

				if (before.endsWith('   ')) {
					removed = 3;
				} else if (before.endsWith('  ')) {
					removed = 2;
				} else if (before.endsWith(' ')) {
					removed = 1;
				}

				query = before.slice(0, before.length - removed) + after;

				requestAnimationFrame(() => {
					const pos = start - removed;
					target.selectionStart = target.selectionEnd = pos;
				});
			} else {
				const SPACES = '   ';
				query = query.slice(0, start) + SPACES + query.slice(end);

				requestAnimationFrame(() => {
					const pos = start + SPACES.length;
					target.selectionStart = target.selectionEnd = pos;
				});
			}
		}
	}}
	class="w-full border border-gray-600 bg-gray-800 p-2 font-mono text-sm"
	rows={4}
></textarea>
