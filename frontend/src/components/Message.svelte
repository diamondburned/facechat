<script>
	import { state } from "../store.js"
	import marked from "marked"

	export let msg // argument

	let squashed = false // TODO
</script>

<div class="message {squashed ? 'squashed' : ''}" id={msg.id}>
	<div class="content">
		{@html marked(msg.markdown)}
	</div>
</div>
