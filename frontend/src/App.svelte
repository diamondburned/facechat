<script>
	import { router, state } from "./store.js"
	import Home from "./pages/Home.svelte"
	import Room from "./pages/Room.svelte"
	import Login from "./pages/Login.svelte"

	import Navaid from "navaid"
	import browserCookies from "browser-cookies"
	import { onDestroy } from "svelte"

	let Route
	let uri = location.pathname

	function mustAuth() {
		if (!browserCookies.get("token")) {
			$router.route("/login", true)
			return true
		}
		return false
	}

	$router.on("/", params => {
		if (mustAuth()) return

		Route = Home
	})

	$router.on("/rooms/:roomID", params => {
		if (mustAuth()) return

		state.roomID = params.roomID
		Route = Room
	})

	$router.on("/login", params => {
		// Redirect to homepage if the user is logged in.
		if (browserCookies.get("token")) {
			$router.route("/", true)
			return
		}

		Route = Login
	})

	$router.listen()

	onDestroy($router.unlisten)
</script>

<svelte:head>
	<link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre.min.css">
	<link rel="stylesheet" href="https://unpkg.com/spectre.css/dist/spectre-exp.min.css">
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
</svelte:head>

<style>
	.app {
		min-height: 100vh;
	}
</style>

<div class="app">
	<svelte:component this={Route} />
</div>
