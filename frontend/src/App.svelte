<script>
	import { router, state } from "./store.js"
	import Home from "./pages/Home.svelte"
	// import Auth from "./pages/Authenticate.svelte"

	import Navaid from "navaid"
	import { onDestroy } from "svelte"

	let Route
	let uri = location.pathname

	$router.on("/", params => {
		Route = Home
	})

	$router.on("/rooms/:room_id", params => {
		state.room_id = params.room_id
		Route = Room
	})

	$router.listen()

	onDestroy($router.unlisten)
</script>
