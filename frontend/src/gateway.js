export default function Gateway(state) {
	let properties = {
		ws: null,
		state: state,
		eventHandlers: [],
	}

	properties.open = function (onDone, onError) {
		if (properties.ws) return

		properties.ws = new WebSocket(`ws://${location.host}/api/gateway`)
		properties.ws.addEventListener("message", (event) => {
			let ev = JSON.parse(event.data)

			switch (ev.e) {
				case "Ready":
					properties.state.me = ev.d.me
					properties.state.rooms = ev.d.public_rooms
					properties.state.privateRooms = ev.d.private_rooms

					if (onDone) onDone(ev.d)

					properties.ws.send(JSON.stringify({ e: "Hello" }))
					break

				case "Message":
					var msgs = [
						ev.d,
						...properties.state.roomMessages[properties.state.roomID],
					]
					if (msgs.length > 35) {
						msgs = msgs.slice(0, 35)
					}
					properties.state.roomMessages[properties.state.roomID] = msgs
					break
			}
		})
		properties.ws.addEventListener("error", (event) => {
			if (onError) onError(event)
		})
	}

	properties.close = function () {
		properties.ws.close()
	}

	return properties
}
