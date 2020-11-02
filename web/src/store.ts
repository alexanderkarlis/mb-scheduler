import { Writable, writable } from "svelte/store";

export const user = writable({
    name: "",
    username: "",
    password: "",
});

export const activeWeekday = writable("");
export const activeWeekdayTimes = writable([""]);

interface IFrequency {
    day: string;
    time: string;
    freq: number;
}

export const Frequency: Writable<IFrequency> = writable({
    day: null,
    time: null,
    freq: null,
})
export const FrequencyArray: Writable<Array<IFrequency>> = writable([])
