<script lang="ts">
	import { countyInfos } from '$lib/utills/flags';
	let { info }: { info: Record<string, unknown> } = $props();
</script>

<div>
	<h1 class="text-2xl font-bold">
		{info.country}
		<span
			style="font-family: 'Segoe UI Emoji', 'Segoe UI Symbol', 'Segoe UI', 'Apple Color Emoji', 'Twemoji Mozilla', 'Noto Color Emoji', 'Android Emoji';"
		>
			{Object.values(countyInfos).find((c) => c.alpha3Code === info.gid_0)?.emoji}
		</span>
	</h1>
	<table class="mt-4">
		<thead>
			<tr class="border-b border-gray-200">
				<th class="pr-2 text-left font-semibold">info</th>
			</tr>
		</thead>
		<tbody>
			{#each Object.entries(info) as [key, value] (key)}
				<tr class="border-b border-gray-200">
					<td class="pr-2"
						>{key
							.toLowerCase()
							.replaceAll('_', ' ')
							.split(' ')
							.filter((w) => {
								const n = Number(w);
								return !(typeof n === 'number' && !isNaN(n));
							})
							.join(' ')}</td
					>
					<td class="p-2">{value}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
