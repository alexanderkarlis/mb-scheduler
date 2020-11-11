<script>
    import { onMount } from "svelte";
    import { user, FrequencyArray } from "./store";

    let scheduleArray = [];

    const getScheduledTimes = async () => {
        console.log("data component mounted");
        await fetch(`http://0.0.0.0:8888/all_times`)
            .then((r) => r.json())
            .catch((e) => {
                console.log(e);
            })
            .then((data) => {
                scheduleArray = [...data];
            });
    };
    onMount(async () => {
        await getScheduledTimes();
    });

    const removeFromArray = (e) => {
        const idx = e.target.value;
        $FrequencyArray.splice(idx, 1);
        $FrequencyArray = [...$FrequencyArray];
        console.log($FrequencyArray);
    };

    const scheduleClass = async (e) => {
        const idx = e.target.value;
        console.log($user, $FrequencyArray[idx]);
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
        console.log(reqObj);
        let r = await fetch("http://localhost:8888/", {
            method: "post",
            body: JSON.stringify(reqObj),
        });
        let resJson = await r.json();
        console.log(resJson);
        removeFromArray({ target: { value: idx } });
        await getScheduledTimes();
    };

    const deleteScheduledRun = async (e) => {
        const idx = e.target.value;
        console.log(idx);
        console.log(scheduleArray[idx].runtime);
        const rt = { runtime: scheduleArray[idx].runtime };
        let r = await fetch("http://localhost:8888/delete_schedule", {
            method: "post",
            body: JSON.stringify(rt),
        });
        let resJson = await r.json();
        console.log(resJson);
    };

    const epochToDateString = (epoch) => {
        const d = new Date(epoch * 1000);
        return d.toGMTString();
    };
</script>

<style>
    #container {
        margin-left: 40px;
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
    #tbl-container {
        height: 500px;
        overflow-y: scroll;
    }
    thead tr th {
        height: 30px;
        line-height: 30px;
        /* text-align: left; */
    }
    /* table#data-tbl thead {
        display: block;
    } */
</style>

<div id="container">
    <div>
        <h2>Queue Sign-up</h2>
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
                                <button value={i} on:click={scheduleClass}>✅
                                </button>
                                <button value={i} on:click={removeFromArray}>❌
                                </button>
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
    <h2>Scheduled Runs</h2>
    <div id="tbl-container">
        <table id="data-tbl">
            <thead>
                <tr class="tbl-header">
                    <th id="tbl-header-col" scope="col" />
                    <th id="tbl-header-col" scope="col">username</th>
                    <th id="tbl-header-col" scope="col">scheduled run-time</th>
                    <th id="tbl-header-col" scope="col">class date</th>
                    <th id="tbl-header-col" scope="col">class day</th>
                    <th id="tbl-header-col" scope="col">class time</th>
                    <th id="tbl-header-col" scope="col">delete</th>
                </tr>
            </thead>
            <tbody>
                {#each scheduleArray as sched, i}
                    <tr>
                        <td>{i}</td>
                        <td>{sched.username}</td>
                        <td>{epochToDateString(sched.runtime)}</td>
                        <td>{sched.date}</td>
                        <td>{sched.weekday}</td>
                        <td>{sched.classtime}</td>
                        <td>
                            <div class="act-btn-box">
                                <button
                                    value={i}
                                    on:click={deleteScheduledRun}>❌
                                </button>
                            </div>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
