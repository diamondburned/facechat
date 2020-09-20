<script>
	export let error
</script>

<style>
	div.error {
		height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	span.material-icons {
		font-size: 5em;
		color: darkred;
	}

	small.error-msg {
		margin: 1em 0;
	}
</style>

<div class="error">
	<span  class="material-icons">error</span>
	<small class="error-msg">{error}</small>
</div>
