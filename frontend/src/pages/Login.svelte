<script>
	let email = "",
		password = "",
		error

	import { router } from "../store.js"
	import Error from "../components/Error.svelte"

	async function submit() {
		try {
			let resp = await fetch("/api/login", {
				method: "POST",
				body: JSON.stringify({
					email: email,
					password: password,
				}),
			})

			if (!resp.ok) {
				throw `Unexpected status code ${resp.status}`
			}

			$router.route("/", true)
		} catch(err) {
			error = err
		}
	}
</script>

<style>
	div.login {
		height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	form {
		width: 250px;
	}
</style>

{#if error}
	<Error {error} />
{:else}
	<div class="login">
		<form on:submit|preventDefault={submit}>
			<div class="form-group">
				<label class="form-label" for="email">Email</label>
				<input class="form-input" type="text" id="email" required
					   placeholder="somebody@something.com"
					   bind:value={email}
				>
			
				<label class="form-label" for="password">Password</label>
				<input class="form-input" type="password" id="password" required
					   bind:value={password}
				>
			</div>

			<div class="form-group">
				<button type="submit" class="btn btn-primary">Log in</button>
			</div>
		</form>
	</div>
{/if}
