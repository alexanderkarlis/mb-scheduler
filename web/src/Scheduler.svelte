<script lang="typescript">
    import {
        activeWeekday,
        activeWeekdayTimes,
        Frequency,
        FrequencyArray,
    } from "./store";
    import Popover from "svelte-popover";

    let daysOfWeek = [
        "Sunday",
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
    ];
    let weekdayTimes: Array<string> = [
        "5:45am",
        "6:45am",
        "7:45am",
        "8:45am",
        "4:45pm",
        "5:45pm",
        "6:45pm",
        "7:45pm",
    ];
    let saturdayTimes: Array<string> = ["9:00am", "10:00am"];
    let sundayTimes: Array<string> = ["9:45am"];

    type GenericObject = { [key: string]: string[] };
    let classTimes: GenericObject = {
        Sunday: sundayTimes,
        Monday: weekdayTimes,
        Tuesday: weekdayTimes,
        Wednesday: weekdayTimes,
        Thursday: weekdayTimes,
        Friday: weekdayTimes,
        Saturday: saturdayTimes,
    };
    const handleBtn = (e, day, time) => {
        console.log(e.target.value);
        $Frequency.weekday = day;
        $Frequency.freq = e.target.value;
        $Frequency.time = time;
        console.log($Frequency);
        $FrequencyArray = [
            ...$FrequencyArray,
            { weekday: day, freq: e.target.value, time: time },
        ];
        console.log([...$FrequencyArray]);
    };
</script>

<style>
    .group {
        display: flex;
        flex-wrap: wrap;
    }
    #day {
        border: "1px solid black";
        color: darkgreen;
        cursor: pointer;
    }
    #weekdayContainer {
        width: 800px;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }
    #clear-day-btn {
        height: 30px;
        background-color: "#efefef";
    }
    .content {
        width: 150px;
        padding: 10px;
        background: #fff;
    }
    .freq-btn {
        background-color: rgb(226, 160, 224);
    }
</style>

<main>
    <button
        id="clear-day-btn"
        on:click={() => {
            $activeWeekday = '';
        }}>clear
    </button>
    <div id="weekdayContainer">
        {#each daysOfWeek as day}
            <div>
                <button
                    on:click={() => {
                        $activeWeekdayTimes = classTimes[day];
                        $activeWeekday = day;
                    }}
                    id="day">{day}
                </button>
                {#if $activeWeekday === day}
                    {#each classTimes[day] as time}
                        <div class="group">
                            <Popover arrowColor="#fff">
                                <button slot="target">{time}</button>
                                <div slot="content" class="content">
                                    <button
                                        value={1}
                                        on:click={(e) => {
                                            handleBtn(e, day, time);
                                        }}
                                        class="freq-btn">
                                        1 Week
                                    </button>
                                    <button
                                        value={10}
                                        on:click={(e) => {
                                            handleBtn(e, day, time);
                                        }}
                                        class="freq-btn">
                                        10 Weeks
                                    </button>
                                    <button
                                        value={25}
                                        on:click={(e) => {
                                            handleBtn(e, day, time);
                                        }}
                                        class="freq-btn">
                                        25 Weeks
                                    </button>
                                    <button
                                        value={99}
                                        on:click={(e) => {
                                            handleBtn(e, day, time);
                                        }}
                                        class="freq-btn">
                                        99 Weeks
                                    </button>
                                </div>
                            </Popover>
                        </div>
                    {/each}
                {/if}
            </div>
        {/each}
    </div>
</main>
