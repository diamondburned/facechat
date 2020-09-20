<script>
	import { router, state } from "./store.js"
	import Home from "./pages/Home.svelte"
	import Login from "./pages/Login.svelte"
	import Register from "./pages/Register.svelte"
	import Accounts from "./pages/Accounts.svelte"

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

	function mustNotAuth() {
		if (browserCookies.get("token")) {
			$router.route("/", true)
			return true
		}
		return false
	}

	$router.on("/", params => {
		if (mustAuth()) return

		$state.roomID = ""
		Route = Home
	})

	$router.on("/:roomID", params => {
		if (mustAuth()) return

		$state.roomID = params.roomID
		Route = Home
	})

	$router.on("/accounts", params => {
		Route = Accounts
	})

	$router.on("/login", params => {
		if (mustNotAuth()) return

		Route = Login
	})

	$router.on("/register", params => {
		if (mustNotAuth()) return

		Route = Register
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
