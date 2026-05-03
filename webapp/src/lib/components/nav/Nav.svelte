<script module lang="ts">
	export type SidePanelView = 'adm' | 'geojson';
</script>

<script lang="ts">
	import GlobSouthAmericaIcon from '$lib/icons/GlobSouthAmericaIcon.svelte';
	import JsonIcon from '$lib/icons/JsonIcon.svelte';

	import { colors } from '$lib/utills/colors';

	type SidePanelView = 'adm' | 'geojson';

	let {
		sidePanelView,
		onSelectView
	}: { sidePanelView: SidePanelView; onSelectView: (view: SidePanelView) => void } = $props();
</script>

<nav style:background-color={colors.offWhite}>
	<ul>
		<li
			onclick={() => onSelectView('adm')}
			role="none"
			class="nav-btn"
			class:active={sidePanelView === 'adm'}
		>
			<GlobSouthAmericaIcon />
		</li>
		<li
			onclick={() => onSelectView('geojson')}
			role="none"
			class="nav-btn"
			class:active={sidePanelView === 'geojson'}
		>
			<JsonIcon />
		</li>
	</ul>
</nav>

<style>
	nav {
		display: flex;
		flex-direction: column;
		width: 50px;
		min-width: 50px;
		border-radius: 10px;
		box-shadow: 0 0 4px 0 rgba(0, 0, 0, 0.2);
	}

	nav ul {
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	nav ul li {
		height: 50px;
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		color: white;
	}

	.nav-btn {
		cursor: pointer;
		border-radius: 10px;
	}
	.nav-btn:hover,
	.nav-btn.active {
		background-color: #ffcb77;
	}
</style>
