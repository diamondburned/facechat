<script>
	import Loading from "../components/Loading.svelte"
	import Error from "../components/Error.svelte"

	import { gateway, state, router } from "../store.js"
	import { throttle } from "throttle-debounce"
	import { fetchx } from "../fetch.js"

	let searchValue = "",
		searchChannels = [],
		searching = false

	const searchThrottler = throttle(500, async (searchValue) => {
		searchChannels = await $state.searchChannels(searchValue)
		searching = false
	})

	$: {
		searchChannels = []

		if (searchValue) {
			searching = true
			searchThrottler(searchValue)
		}
	}

	function toRoom(room) {
		$state.updateRoom(room)
		$router.route("/" + room.id, true)
	}

	let showCreateModal = false
	function toggleCreateModal() {
		showCreateModal = !showCreateModal
	}

	$: createData = {
		name: searchValue,
		level: 0,
	}

	async function submitNewChannel() {
		let resp = await fetchx("/api/room", {
			method: "POST",
			body: JSON.stringify(createData),
		})

		let room = await resp.json()

		$state.updateRoom(room)
		$router.route("/" + room.id, true)
	}
</script>

<style>
	button.create {
		display: flex;
		margin:  auto;
	}
</style>

<div class="column col-3">
	<input class="search" type="text" placeholder="Search channels..."
		   bind:value={searchValue}
	>
	<div class="channels">
		{#if searching}
			<Loading />
		{:else if searchChannels.length == 0}
			{#if searchValue != ""}
				<button type="button" class="create btn btn-sm btn-link" on:click={toggleCreateModal}>
					<span class="material-icons">add</span>
					<span>Create Room</span>
				</button>
			{/if}
		{:else}
			<div class="channel-list">
				{#each searchChannels as channel (channel.id)}
					<button type="button" class="channel btn btn-link" on:click={toRoom (channel)}>
						<div><h5>#{channel.name}</h5></div>
						<div><span>{channel.topic}</span></div>
					</button>
				{/each}
			</div>
		{/if}
	</div>

	{#if showCreateModal}
		<div class="modal active">
			<a href="#close" class="modal-overlay" aria-label="Close"
			   on:click={toggleCreateModal}></a>
			<div class="modal-container">
				<div class="modal-header">
					<a href="#close" class="btn btn-clear float-right" aria-label="Close"
					   on:click={toggleCreateModal}
					></a>
					<div class="modal-title h5">Create Room</div>
				</div>
				<form class="modal-body" on:submit|preventDefault={submitNewChannel}>
					<div class="form-group">
						<label class="form-label" for="name">Room Name</label>
						<input class="form-input" type="text" id="name" required
							   bind:value={createData.name}
						>

						<label class="form-label">Privacy Level</label>
						<label class="form-radio">
							<input type="radio" bind:group={createData.level} value={0}>
							<i class="form-icon"></i>
							<span>Anonymous</span>
						</label>
						<label class="form-radio">
							<input type="radio" bind:group={createData.level} value={1}>
							<i class="form-icon"></i>
							<span>Half Open</span>
						</label>
						<label class="form-radio">
							<input type="radio" bind:group={createData.level} value={2}>
							<i class="form-icon"></i>
							<span>Fully Open</span>
						</label>
					</div>
					<div class="form-group">
						<button type="submit" class="btn btn-primary">Register</button>
					</div>
				</form>
			</div>
		</div>
	{/if}
</div>
