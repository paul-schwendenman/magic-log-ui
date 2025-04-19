<script lang="ts">
	export let data: unknown;

	// Safe fallback string
	let fallback = JSON.stringify(data);

	// Convert JS value to a stylized fragment
	function render(value: unknown): any {
		if (typeof value === 'string') {
			return `<span class="text-green-500">"${escapeHtml(value)}"</span>`;
		} else if (typeof value === 'number') {
			return `<span class="text-blue-500">${value}</span>`;
		} else if (typeof value === 'boolean') {
			return `<span class="text-rose-500">${value}</span>`;
		} else if (value === null) {
			return `<span class="text-gray-500 italic">null</span>`;
		} else if (Array.isArray(value)) {
			return `[${value.map(render).join(', ')}]`;
		} else if (typeof value === 'object' && value !== null) {
			return `{${Object.entries(value)
				.map(
					([key, val]) => `<span class="text-teal-500">"${escapeHtml(key)}"</span>: ${render(val)}`
				)
				.join(', ')}}`;
		} else {
			return `<span class="text-gray-400">${escapeHtml(String(value))}</span>`;
		}
	}

	function escapeHtml(str: string) {
		return str
			.replace(/&/g, '&amp;')
			.replace(/</g, '&lt;')
			.replace(/>/g, '&gt;')
			.replace(/"/g, '&quot;')
			.replace(/'/g, '&#039;');
	}

	$: rendered = render(data);
</script>

<!-- ✅ Syntax-highlighted inline JSON -->
<div class="overflow-hidden font-mono text-sm text-ellipsis whitespace-nowrap">
	{@html rendered || fallback}
</div>
