<script>
    import { onMount } from "svelte";

    let scheduleArray = [];

    export async function getScheduledTimes() {
        console.log("data component mounted");
        await fetch(`http://0.0.0.0:8888/all_times`)
            .then((r) => r.json())
            .catch((e) => {
                console.log(e);
            })
            .then((data) => {
                if (data && data.length) {
                    scheduleArray = [...data];
                } else {
                    scheduleArray = [];
                }
            });
    }

    onMount(async () => {
        await getScheduledTimes();
    });

    const deleteScheduledRun = async (e) => {
        const idx = e.target.value;
        const rt = { runtime: scheduleArray[idx].runtime };
        let r = await fetch("http://localhost:8888/delete_schedule", {
            method: "post",
            body: JSON.stringify(rt),
        });
        let resJson = await r.json();
        if (resJson && resJson.length) {
            scheduleArray = [...resJson];
        } else {
            scheduleArray = [];
        }
    };

    const epochToDateString = (epoch) => {
        var d = new Date(epoch * 1000 + 5 * 60 * 60 * 1000);
        return `${d.toDateString()} ${d.toLocaleTimeString()}`;
    };
</script>

<style>
    #container {
        display: flex;
        flex-direction: column;
        padding: 30px;
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
    }
    .subtitle {
        font-size: 14px;
        font-weight: 300;
        text-transform: capitalize;
    }
    .tbl-outln {
        padding: 25px;
    }
</style>

<div id="container" class="w3-card-4">
    <div class="w3-card-2 tbl-outln">
        <h2>Scheduled Runs</h2>
        <span class="subtitle">
            List of all scheduled times ready to be run. Click the ❌ to remove
            from running (remove from database).
        </span>
        <br />
        <br />
        <div id="tbl-container">
            <table id="data-tbl">
                <thead>
                    <tr class="tbl-header">
                        <th id="tbl-header-col" scope="col" />
                        <th id="tbl-header-col" scope="col">username</th>
                        <th id="tbl-header-col" scope="col">
                            scheduled run-time
                        </th>
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
</div>
