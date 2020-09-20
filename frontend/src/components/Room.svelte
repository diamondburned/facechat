<script>
	import { state } from "../store.js"
	import { onMount } from "svelte"

	import Loading from "../components/Loading.svelte"
	import Message from "../components/Message.svelte"
	import Error from "../components/Error.svelte"

	import SimpleMDE from "simplemde"
	import marked from "../marked.js"

	let editor,
		message,
		error,
		simplemde

	$: if (editor) {
		simplemde = new SimpleMDE({
			element: editor,
			toolbar: false,
			status: false,
			placeholder: "Send a message...",
			previewRender: (text) => marked(text),
		})
	}

	$: room = $state.room[$state.roomID]
	$: messages = $state.roomMessages[$state.roomID]

	$: if ($state.roomID != "") {
		(async () => {
			// Update the state.
			try {
				await Promise.all([
					$state.currentRoom(),
					$state.currentMessages(),
				])
			} catch(err) {
				error = err
			}
		})()
	}

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

<div id={room ? room.id : null} class="room column col-9">
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
					bind:this={editor} bind:value={message}
					on:keypress={sendMessage}
				/>
			{:else}
				<Loading />
			{/if}
		</main>
	{/if}
</div>
