<script>
	import { state } from "../store.js"
	import { onMount } from "svelte"

	let room = $state.currentRoom()

	onMount(async () => {
		if (room) {
			return
		}

		// Fetch the room if it's not in the state.
		let resp = await fetch("/api/room/" + $state.room_id)
		room = await resp.json()
		$state.rooms[room.id] = room
	})
</script>

<div id={ $state.room_id } class="room">
</div>
