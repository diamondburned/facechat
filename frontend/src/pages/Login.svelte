<script>
	let login = {
			email = "",
			password = "",
		},
		error

	import { router } from "../store.js"

	async function submit() {
		try {
			await fetch("/api/login", {
				method: "POST",
				body: JSON.stringify(login),
			})

			$router.route("/", true)
		} catch(err) {
			error = err
		}
	}
</script>

<div class="login">
	{#if error}
		<Error {error} />
	{:else}
		<form on:submit|preventDefault={submit}>
			<label class="form-label" for="email">Email</label>
			<input class="form-input" type="text" id="email" required
				   placeholder="somebody@something.com"
				   bind:value={$login.email}
			>
		
			<label class="form-label" for="password">Password</label>
			<input class="form-input" type="password" id="password" required
				   bind:value={$login.password}
			>
		</form>
	{/if}
</div>
