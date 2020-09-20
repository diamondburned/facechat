export class Gateway {
	constructor(state) {
		this.state = state
		this.eventHandlers = []
	}

	open(onDone) {
		this.ws = new WebSocket("/api/gateway")
		this.ws.addEventListener("message", (event) => {
			let ev = JSON.parse(event.data)

			switch (ev.d) {
				case "Ready":
					this.store.me = ev.me
					this.store.rooms = ev.public_rooms
					this.store.privateRooms = ev.private_rooms
					break
				case "Message":
					var msgs = [ev.d, ...this.store.roomMessages[this.store.roomID]]
					if (msgs.length > 35) {
						msgs = msgs.slice(0, 35)
					}
					this.store.roomMessages[this.store.roomID] = msgs
					break
			}
		})

		if (onDone) {
			this.ws.addEventListener("open", onDone)
		}
	}

	close() {
		this.ws.close()
	}
}
