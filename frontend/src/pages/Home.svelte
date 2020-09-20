<script>
	import Channels from "../components/Channels.svelte"
	import Room from "../components/Room.svelte"
	import Loading from "../components/Loading.svelte"
	import Error from "../components/Error.svelte"
	import User from "../components/User.svelte"

	import { gateway, state, router } from "../store.js"
	import { onMount } from "svelte"
	
	let loading = true,
		accounts,
		error

	onMount(async () => {
		accounts = await $state.fetchAccounts()

		if (accounts.length < 1) {
			$router.route("/accounts", true)
			return
		}

		$gateway.open(
			(ev) => { loading = false },
			(ev) => { error = `${ev}` },
		)
	})
</script>

<style>
	div.columns {
		margin: 10px;
	}
</style>

{#if error}
	<Error {error} />
{:else if loading}
	<Loading />
{:else}
	<div class="columns">
		<Channels />
		<Room />
	</div>
{/if}
