<script>
    import { user, FrequencyArray } from "./store";

    console.log($FrequencyArray);
    const scheduleClass = async () => {
        let reqObj = {
            username: "alexanderkarlis@gmail.com",
            name: "Alexander Karlis",
            password: "921921Zz?",
            schedule: {
                classtime: "5:45pm",
                date: "11/11/2020",
                frequency: "10",
            },
        };
        let r = await fetch("http://localhost:8888", {
            method: "post",
            body: JSON.stringify(reqObj),
        });
        let resJson = await r.json();
        console.log(resJson);
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
</style>

<div id="container">
    <h2>Queue signups</h2>
    <table id="data-tbl">
        <tr>
            <th id="tbl-header-col" scope="col" />
            <th id="tbl-header-col" scope="col">username</th>
            <th id="tbl-header-col" scope="col">day of week</th>
            <th id="tbl-header-col" scope="col">time</th>
            <th id="tbl-header-col" scope="col">freq</th>
            <th id="tbl-header-col" scope="col">action</th>
        </tr>
        {#each $FrequencyArray as freq, i}
            <tr>
                <td>{i}</td>
                <td>{$user.username}</td>
                <td>{freq.weekday}</td>
                <td>{freq.time}</td>
                <td>{freq.freq}</td>
                <td>
                    <div class="act-btn-box"><button>✅</button> <button>❌</button></div>
                </td>
            </tr>
        {/each}
    </table>
    <h2>Scheduled Runs</h2>
    <table id="data-tbl">
        <tr>
            <th id="tbl-header-col" scope="col" />
            <th id="tbl-header-col" scope="col">schedule time</th>
            <th id="tbl-header-col" scope="col">day of week</th>
            <th id="tbl-header-col" scope="col">date</th>
            <th id="tbl-header-col" scope="col">class time</th>
        </tr>
        {#each $FrequencyArray as freq, i}
            <tr>
                <td>{i}</td>
                <td>{$user.username}</td>
                <td>{freq.day}</td>
                <td>{freq.time}</td>
            </tr>
        {/each}
    </table>
    <button on:click={() => scheduleClass()}>run</button>
</div>
