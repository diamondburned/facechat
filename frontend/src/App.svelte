<script>
	import { router, state } from "./store.js"
	import Home from "./pages/Home.svelte"
	import Room from "./pages/Room.svelte"

	import Navaid from "navaid"
	import { onDestroy } from "svelte"

	let Route
	let uri = location.pathname

	$router.on("/", params => {
		Route = Home
	})

	$router.on("/rooms/:roomID", params => {
		state.roomID = params.roomID
		Route = Room
	})

	$router.listen()

	onDestroy($router.unlisten)
</script>

<svelte:head>
	<link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre.min.css">
	<link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre-exp.min.css">
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
</svelte:head>

<div class="app">
	<svelte:component this={Route} />
</div>
