import Navaid from "navaid"
import { writable } from "svelte/store"
import { State } from "./state.js"

export const router = writable(Navaid("/"))
export const state = writable(new State())
