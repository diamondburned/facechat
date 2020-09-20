<script>
	let email = "",
		username = "",
		password = "",
		error

	import { router } from "../store.js"
	import { fetchx } from "../fetch.js"
	import Error from "../components/Error.svelte"

	async function submit() {
		try {
			let resp = await fetchx("/api/register", {
				method: "POST",
				body: JSON.stringify({
					email: email,
					username: username,
					password: password,
				}),
			})

			$router.route("/", true)
		} catch(err) {
			error = err
		}
	}
</script>

<style>
	div.register {
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
	<div class="register">
		<form on:submit|preventDefault={submit}>
			<div class="form-group">
				<label class="form-label" for="username">Username</label>
				<input class="form-input" type="text" id="username" required
					   bind:value={username}
				>

				<label class="form-label" for="email">Email</label>
				<input class="form-input" type="text" id="email" required
					   bind:value={email}
				>
			
				<label class="form-label" for="password">Password</label>
				<input class="form-input" type="password" id="password" required
					   bind:value={password}
				>
			</div>

			<div class="form-group">
				<button type="submit" class="btn btn-primary">Register</button>
			</div>
		</form>
	</div>
{/if}
