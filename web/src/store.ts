import { writable } from "svelte/store";

export const user = writable({
    name: "",
    username: "",
    password: "",
});

export const activeWeekDay = writable("");
