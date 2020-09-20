import Navaid from "navaid"
import State from "./state.js"
import Gateway from "./gateway.js"
import { writable } from "svelte/store"

const internalState = State()

export const router = writable(Navaid("/"))
export const state = writable(internalState)
export const gateway = writable(Gateway(internalState))
