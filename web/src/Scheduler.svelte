<script lang="typescript">
    import { ScheduleArray } from "./store";
    import { onMount } from "svelte";
    import {
        activeWeekday,
        activeWeekdayTimes,
        Frequency,
        FrequencyArray,
        PostAttemptMessage,
        user,
    } from "./store";
    import Popover from "svelte-popover";
    import Alert from "./Alert.svelte";
    import { MESSAGE_TYPES } from "./constants";

    const getScheduledTimes = async () => {
        console.log("getting new scheduled times");
        await fetch(`http://0.0.0.0:8888/all_times`)
            .then((r) => r.json())
            .catch((e) => {
                console.log(e);
            })
            .then((data) => {
                console.log(data);
                if (data && data.length) {
                    $ScheduleArray = [...data];
                } else {
                    $ScheduleArray = [];
                }
            });
    };

    onMount(async () => {
        await getScheduledTimes();
    });

    setInterval(async () => {
        try {
            await getScheduledTimes();
        } catch (e) {
            console.log(e);
        }
    }, 5000);

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

    const removeFromArray = (e) => {
        const idx = e.target.value;
        $FrequencyArray.splice(idx, 1);
        $FrequencyArray = [...$FrequencyArray];
        console.log($FrequencyArray);
    };

    let postSuccessStatus = false;

    const scheduleClass = async (e) => {
        const idx = e.target.value;
        const freqObj = $FrequencyArray[idx];
        let reqObj = {
            username: $user.username,
            name: $user.name,
            password: $user.password,
            schedule: {
                classtime: freqObj.time,
                weekday: freqObj.weekday,
                frequency: freqObj.freq,
            },
        };
        let r = await fetch("http://0.0.0.0:8888/", {
            method: "post",
            body: JSON.stringify(reqObj),
        });
        postSuccessStatus = r.status === 200;
        $PostAttemptMessage = true;
        let resJson = await r.json();

        removeFromArray({ target: { value: idx } });
        await getScheduledTimes();
    };
</script>

<style>
    .group {
        display: flex;
        flex-wrap: wrap;
    }
    .content {
        width: 150px;
        padding: 10px;
        background: #fff;
    }
    .freq-btn {
        background-color: rgb(226, 160, 224);
    }
    .main {
        padding-top: 30px;
        padding-bottom: 30px;
        display: flex;
        flex-direction: column;
    }
    #tbl-header-col {
        text-transform: capitalize;
    }
    #data-tbl {
        font-family: Arial, Helvetica, sans-serif;
        border-collapse: collapse;
        width: 100%;
    }

    #data-tbl td,
    #data-tbl th {
        border: 1px solid #ddd;
        padding: 8px;
    }

    #data-tbl tr:nth-child(even) {
        background-color: #f2f2f2;
    }

    #data-tbl tr:hover {
        background-color: #ddd;
    }

    #data-tbl th {
        padding-top: 12px;
        padding-bottom: 12px;
        text-align: left;
        background-color: #eb8531;
        color: white;
    }
    .act-btn-box {
        display: flex;
        flex-direction: row;
    }
    thead tr th {
        height: 30px;
        line-height: 30px;
    }
    .subtitle {
        font-size: 14px;
        font-weight: 300;
        text-transform: capitalize;
    }
    .tbl-outln {
        padding: 25px;
        display: flex;
    }
</style>

<div class="main w3-card-2 tbl-outln">
    <div style="padding: 15px;">
        <h2>Queued Sign-up Scheduler</h2>
        <span class="subtitle">
            Select a day and time. Review the information to be saved. Click the
            ✅ to confirm, or the ❌ to remove from insertion
        </span>
    </div>
    <div
        class="w3-card-2 tbl-outln"
        role="toolbar"
        aria-label="Toolbar with button groups">
        <div class="btn-group mr-2" role="group" aria-label="First group">
            <button
                type="button"
                class="btn btn-primary"
                id="clear-day-btn"
                style="height: 35px!important;"
                on:click={() => {
                    $activeWeekday = '';
                }}>Clear</button>
        </div>
        <div class="btn-group mr-2" role="group" aria-label="Second group">
            {#each daysOfWeek as day}
                <div>
                    <button
                        class="btn btn-secondary"
                        on:click={() => {
                            $activeWeekdayTimes = classTimes[day];
                            $activeWeekday = day;
                        }}
                        style="margin-right: 5px;"
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
        <div>
            <div>
                <table id="data-tbl">
                    <thead>
                        <tr>
                            <th id="tbl-header-col" scope="col" />
                            <th id="tbl-header-col" scope="col">username</th>
                            <th id="tbl-header-col" scope="col">day of week</th>
                            <th id="tbl-header-col" scope="col">time</th>
                            <th id="tbl-header-col" scope="col">freq</th>
                            <th id="tbl-header-col" scope="col">action</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each $FrequencyArray as freq, i}
                            <tr>
                                <td>{i}</td>
                                <td>{$user.username}</td>
                                <td>{freq.weekday}</td>
                                <td>{freq.time}</td>
                                <td>{freq.freq}</td>
                                <td>
                                    <div class="act-btn-box">
                                        <button
                                            value={i}
                                            on:click={scheduleClass}>✅
                                        </button>
                                        <button
                                            value={i}
                                            on:click={removeFromArray}>❌
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </div>
        <Alert
            message={postSuccessStatus ? 'Successfully added scheduled run(s)' : 'Error adding scheduled run(s)'}
            type={postSuccessStatus ? MESSAGE_TYPES.SUCCESS : MESSAGE_TYPES.FAILED}
            open={true} />
    </div>
</div>
