export class State {
	constructor() {
		this.roomID = 0
		this.rooms = {}
		this.roomMessages = {} // map[id][]Messages
		this.users = {}
	}

	user(id) {
		return this.users[id]
	}

	async currentMessages() {
		return this.messages(this.roomID)
	}

	async messages(roomID) {
		let messages = this.roomMessages[roomID]
		if (!messages) {
			let resp = await fetch("/api/room/" + roomID + "/messages")
			messages = await resp.json()
			this.roomMessages[roomID] = messages
		}

		return messages
	}

	async currentRoom() {
		return this.room(this.roomID)
	}

	async room(id) {
		let room = this.rooms[id]
		if (!room) {
			// Fetch the room if it's not in the state.
			let resp = await fetch("/api/room/" + id)
			room = await resp.json()
			this.rooms[id] = room
		}

		return room
	}
}
