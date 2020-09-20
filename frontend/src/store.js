import Navaid from "navaid"
import { writable } from "svelte/store"
import { State } from "./state.js"
import { Gateway } from "./gateway.js"

export const router = writable(Navaid("/"))
export const state = writable(new State())
export const gateway = writable(new Gateway(state))
