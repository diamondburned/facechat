<script>
	import Loading from "../components/Loading.svelte"
	import Error from "../components/Error.svelte"

	import { state } from "../store.js"
	import { onMount } from "svelte"

	let knownServices = [
		{
			name: "GitHub",
			login: "/api/oauth/github/login",
			svgHTML: require("simple-icons/icons/github").svg,
			account: null,
		},
		{
			name: "Twitter",
			login: "/api/oauth/twitter/login",
			svgHTML: require("simple-icons/icons/twitter").svg,
			account: null,
		},
	]

	function serviceSVG(name) {
		let svc = knownServices[name]
		if (svc) {
			return svc.svgHTML != ""
		}
		return false
	}

	let accounts,
		error

	onMount(async () => {
		try {
			accounts = await $state.fetchAccounts()
			console.log(accounts)
			accounts.forEach(a => {
				let svc = knownServices.find(s => s.name == a.service)
				svc.account = a
			})
		} catch(err) {
			error = err
		}
	})
</script>

<style>
	ul.account-steps {
		margin-top: 20vh;
	}

	div.account-container {
		margin-top: 5vh;
	}

	a.service-button {
		display: flex;
		flex-direction: row;
	}

	a.service-button.btn-success {
		opacity: 0.65;
		pointer-events: none;
	}

	a.service-button span.name {
		flex: 1;
		text-align: left;
	}

	div.service-icon {
		width:  1rem;
		height: 1rem;
		margin: auto 0;
		margin-right: 8px;
	}
</style>

{#if error}
	<Error {error} />
{:else if accounts}
	<div class="accounts">
		<ul class="account-steps step">
			<li class="step-item"><a href="">Register</a></li>
			<li class="step-item active"><a href="">Link Accounts</a></li>
			<li class="step-item"><a href="">Finish Up</a></li>
		</ul>
		<div class="account-container container grid-sm">
			<h3>Link Accounts</h3>
			<p>Link at least 2 accounts to continue.</p>

			<div class="columns">
				{#each knownServices as service (service.name)}
					<div class="column col-sm-2">
						<a target="_blank"
						   class="btn btn-lg service-button {service.account ? "btn-success" : "btn-error"}"
						   href={service.account ? service.account.url : service.login}>
							<div class="service-icon">
								{@html service.svgHTML}
							</div>
							<span class="name">{service.name}</span>
							<span class="material-icons">launch</span>
						</a>
					</div>
				{/each}
			</div>
		</div>
	</div>
{:else}
	<Loading />
{/if}
