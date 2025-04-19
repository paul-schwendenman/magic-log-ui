<script lang="ts">
	export let query = '';
	export let onQuery: () => void = () => {};
</script>

<textarea
	bind:value={query}
	class="w-full font-mono border border-gray-600 bg-gray-800 p-2 text-sm"
	rows={4}
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
			const allText = query;

			const lineStart = allText.lastIndexOf('\n', start - 1) + 1;
			const lineEnd = allText.indexOf('\n', end);
			const blockEnd = lineEnd === -1 ? allText.length : lineEnd;

			const selectedBlock = allText.slice(lineStart, blockEnd);
			const blockLines = selectedBlock.split('\n');

			if (e.shiftKey) {
				// OUTDENT
				const updated = blockLines.map((line) =>
					line.startsWith('   ') ? line.slice(3) :
					line.startsWith('  ') ? line.slice(2) :
					line.startsWith(' ') ? line.slice(1) :
					line
				);
				const newText = updated.join('\n');
				query = allText.slice(0, lineStart) + newText + allText.slice(blockEnd);

				requestAnimationFrame(() => {
					const delta = selectedBlock.length - newText.length;
					target.selectionStart = start - delta;
					target.selectionEnd = end - delta;
				});
			} else {
				const SPACES = '   ';
				const updated = blockLines.map((line) => SPACES + line);
				const newText = updated.join('\n');
				query = allText.slice(0, lineStart) + newText + allText.slice(blockEnd);

				requestAnimationFrame(() => {
					const delta = newText.length - selectedBlock.length;
					target.selectionStart = start + SPACES.length;
					target.selectionEnd = end + delta;
				});
			}
		}
	}}
></textarea>
