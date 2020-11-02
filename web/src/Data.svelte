<script>
    import { Frequency, FrequencyArray } from "./store";

    console.log($FrequencyArray);
    const scheduleClass = async () => {
        let reqObj = {
            username: "alexanderkarlis@gmail.com",
            name: "Alexander Karlis",
            password: "921921Zz?",
            schedule: {
                classtime: "5:45pm",
                date: "11/02/2020",
                frequency: "1",
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
    #data-tbl {
        padding: 10px;
    }
    #tbl-header-col {
        text-transform: capitalize;
    }
</style>

<div id="container">
    <table id="data-tbl">
        <tr>
            <th id="tbl-header-col" scope="col">date</th>
            <th id="tbl-header-col" scope="col">day of week</th>
            <th id="tbl-header-col" scope="col">time</th>
            <th id="tbl-header-col" scope="col">freq</th>
        </tr>
        {#each $FrequencyArray as freq}
            <tr>
                <td>1/1/2020</td>
                <td>{freq.day}</td>
                <td>{freq.time}</td>
                <td>{freq.freq}</td>
            </tr>
        {/each}
    </table>
    <button on:click={() => scheduleClass()}>run</button>
</div>
