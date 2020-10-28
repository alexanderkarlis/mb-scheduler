<script lang="typescript">
    import { activeWeekDay } from "./store";

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

    let classTimes = {
        Sunday: sundayTimes,
        Monday: weekdayTimes,
        Tuesday: weekdayTimes,
        Wednesday: weekdayTimes,
        Thursday: weekdayTimes,
        Friday: weekdayTimes,
        Saturday: saturdayTimes,
    };
    let showClassTimes = null;
</script>

<style>
    #day {
        border: "1px solid black";
        color: darkgreen;
        cursor: pointer;
    }
    #weekdayContainer {
        width: 1000px;
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
    }
</style>

<main>
    <button
        on:click={() => {
            $activeWeekDay = '';
        }}>clear</button>
    <div id="weekdayContainer">
        {#each daysOfWeek as day}
            <div>
                <button
                    on:click={() => {
                        console.log(classTimes[day]);
                        $activeWeekDay = day;
                    }}
                    id="day">{day}
                </button>
                {#if $activeWeekDay === day}<span value={day}>times</span>{/if}
            </div>
        {/each}
    </div>
</main>
