import { fetchx } from "./fetch.js"

export default function State() {
	let properties = {
		me: {},
		accounts: null,
		roomID: "",
		rooms: {},
		privateRooms: {},
		roomMessages: {}, // map[id][]Messages
		users: {},
	}

	properties.fetchAccounts = async function (refetch) {
		if (!properties.accounts || refetch) {
			let resp = await fetchx("/api/user/@me/accounts")
			properties.accounts = await resp.json()
		}

		return properties.accounts
	}

	properties.searchChannels = async function (query) {
		let form = new URLSearchParams({ q: query })
		let resp = await fetchx("/api/room?" + form)
		return await resp.json()
	}

	properties.user = function (id) {
		return properties.users[id]
	}

	properties.currentMessages = async function () {
		return properties.messages(properties.roomID)
	}

	properties.messages = async function (roomID) {
		let messages = properties.roomMessages[roomID]
		if (!messages) {
			let resp = await fetchx("/api/room/" + roomID + "/messages")
			messages = await resp.json()
			properties.roomMessages[roomID] = messages
		}

		return messages
	}

	properties.currentRoom = async function () {
		return properties.room(properties.roomID)
	}

	properties.room = async function (id) {
		let room = properties.rooms[id]
		if (!room) {
			// Fetch the room if it's not in the state.
			let resp = await fetchx("/api/room/" + id)
			room = await resp.json()
			properties.updateRoom(room)
		}

		return room
	}

	properties.updateRoom = function (room) {
		properties.rooms[room.id] = room
	}

	return properties
}
