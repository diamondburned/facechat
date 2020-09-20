<script>
	export let account
	export let knownServices
</script>

<a class="btn btn-lg" target="_blank" href={account.url}>
	<div class="service-icon">
		{#if serviceSVG(account.service)}
			{@html knownServices[account.service].svgHTML}
		{:else}
			<div class="material-icons">public</div>
		{/if}
	</div>
	<span>{account.service}</span>
	<div class="material-icons">launch</div>
</a>
