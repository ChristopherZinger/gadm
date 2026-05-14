<script lang="ts">
	import { colors } from '$lib/utills/colors';
	let { height = 24, color = colors.blackAsh }: { height?: number; color?: string } = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	{height}
	viewBox="0 -960 960 960"
	width={height}
	fill={color}><path d="M120-120v-720h720v720H120Zm80-80h560v-560H200v560Zm0 0v-560 560Z" /></svg
>
