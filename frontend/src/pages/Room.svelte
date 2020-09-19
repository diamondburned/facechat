<script>
	import { state } from "../store.js"
	import { onMount } from "svelte"

	import Message from "../components/Message.svelte"
	import Error from "../components/Error.svelte"

	import SimpleMDE from "simplemde"
	import marked from "../marked.js"

	let room     = $state.room[$state.roomID],
		messages = $state.roomMessages[state.roomID],
		editor,
		message,
		error

	onMount(async () => {
		// Update the state.
		try {
			await Promise.all([
				$state.currentRoom(),
				$state.currentMessages(),
			])

			let simplemde = new SimpleMDE({
				element: editor,
				toolbar: false,
				status: false,
				placeholder: "Send a message...",
				previewRender: (text) => marked(text),
			})
		} catch(err) {
			error = err
		}
	})

	function sendMessage(e) {
		if (e.keyCode === 13 && !e.shiftKey) {
			// Not add a new line
			e.preventDefault()

			// Clear the input, prevents the chat from blocking
			let content = message
			message = ""

			console.log("Sending ${message}")
		}
	}
</script>

<div id={room ? room.id : null} class="room">
	{#if error}
		<Error {error} />
	{:else}
		<header class="room-info">
			{#if room}
				<div class="columns left">
					{#if      room.type == 0}
						<span class="material-icons">group</span>
					{:else if room.type == 1}
						<span class="material-icons">person</span>
					{:else}
						<span>Unknown</span>
					{/if}
					<big>{room.name}</big>
					<div class="divider-vert"></div>
					<span>{room.Topic || "No topic."}</span>
				</div>
				<div class="columns right">
					<span>Secret: {room.level}</span>
				</div>
			{:else}
				<span>Loading...</span>
			{/if}
		</header>
		<main>
			{#if messages}
				<div class="messages">
					{#each messages as msg (msg.id)}
						<Message {msg} />
					{/each}
				</div>
				<textarea
					id="message-editor" class="message-editor"
					placeholder="Message to #{room.name}"
					bind:this={editor} bind:value={message}
					on:keypress={sendMessage}
				/>
			{:else}
				<div class="loading loading-lg"></div>
			{/if}
		</main>
	{/if}
</div>
