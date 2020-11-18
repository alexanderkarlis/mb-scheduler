<script>
    import { ScheduleHistoryArray } from "./store";
    import { onMount } from "svelte";

    const getScheduledRunHistory = async (e) => {
        try {
            let r = await fetch("http://0.0.0.0:8888/get_run_history", {
                method: "get",
            });
            let resJson = await r.json().catch((r) => console.log(r));
            if (resJson && resJson.length) {
                $ScheduleHistoryArray = [...resJson];
            } else {
                $ScheduleHistoryArray = [];
            }
        } catch (e) {
            console.log(e);
        }
    };

    const epochToDateString = (epoch) => {
        var d = new Date(epoch * 1000);
        return `${d.toDateString()} ${d.toLocaleTimeString()}`;
    };

    setInterval(async () => {
        try {
            await getScheduledRunHistory();
        } catch (e) {
            console.log(e);
        }
    }, 5000);

    onMount(async () => {
        console.log("in mount");
        try {
            await getScheduledRunHistory();
        } catch (e) {
            console.log(e);
        }
    });
</script>

<style>
    #container {
        padding: 30px;
        display: flex;
        flex-direction: column;
    }
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
    .success {
        background-color: lightgreen!important;
    }
    .failed {
        background-color: #ff9d9b!important;
    }
</style>

<div id="container" class="w3-card-4">
    <div class="w3-card-2 tbl-outln">
        <h2>Run History</h2>
        <span class="subtitle"> List of all previous runs. </span>
        <br />
        <br />
        <div id="tbl-container">
            <table id="data-tbl">
                <thead>
                    <tr class="tbl-header">
                        <th id="tbl-header-col" scope="col" />
                        <th id="tbl-header-col" scope="col">username</th>
                        <th id="tbl-header-col" scope="col">run-time</th>
                        <th id="tbl-header-col" scope="col">class date</th>
                        <th id="tbl-header-col" scope="col">class day</th>
                        <th id="tbl-header-col" scope="col">class time</th>
                        <th id="tbl-header-col" scope="col">run status</th>
                    </tr>
                </thead>
                <tbody>
                    {#each $ScheduleHistoryArray as sched, i}
                        <tr
                            class={sched.status === 'success' ? 'success' : 'failed'}>
                            <td>{i}</td>
                            <td>{sched.username}</td>
                            <td>{epochToDateString(sched.runtime)}</td>
                            <td>{sched.date}</td>
                            <td>{sched.weekday}</td>
                            <td>{sched.classtime}</td>
                            <td>{sched.status}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</div>
