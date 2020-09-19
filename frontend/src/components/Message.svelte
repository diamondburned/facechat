<script>
	import { state } from "../store.js"

	export let msg // argument

	let squashed = false // TODO
</script>

<div class="message {squashed ? 'squashed' : ''}" id={msg.id}>
</div>
