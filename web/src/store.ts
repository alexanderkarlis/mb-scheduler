import { Writable, writable } from "svelte/store";

export const user = writable({
    name: "",
    username: "",
    password: "",
});

export const activeWeekday = writable("");
export const activeWeekdayTimes = writable([""]);

interface Frequency {
    day: string;
    time: string;
    freq: number;
}

export const Frequency: Writable<Frequency> = writable({
    day: null,
    time: null,
    freq: null,
})
export const FrequencyArray: Writable<Array<Frequency>> = writable([])
