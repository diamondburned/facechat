import Navaid from "navaid"
import { writable } from "svelte/store"

export const router = writable(Navaid("/"))
export const state = writable(State())

class State {
	constructor() {
		this.room_id = 0
		this.rooms = {}
		this.users = {}
	}

	currentRoom() {
		return this.rooms[this.room_id]
	}

	user(id) {
		return this.users[id]
	}
}
