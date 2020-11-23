<script lang="typescript">
    import Scheduler from "./Scheduler.svelte";
    import Data from "./ScheduledData.svelte";
    import Login from "./Login.svelte";
    import RunHistory from "./RunHistory.svelte";

    let serverStatus = null;
    let n = new Date();

    const DOW = {
        1: "Monday",
        2: "Tuesday",
        3: "Wednesday",
        4: "Thursday",
        5: "Friday",
        6: "Saturday",
        0: "Sunday",
    };
    let timeNow = `${n.toLocaleDateString()} ${n.toLocaleTimeString()}`;
    var myfunc = setInterval(function () {
        var now = new Date();
        timeNow = `${
            DOW[now.getDay()]
        } - ${now.toLocaleDateString()} ${now.toLocaleTimeString()}`;
    }, 1000);
</script>

<style>
    #app-container {
        margin: 50px;
    }
    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }
    #opts {
        display: flex;
        flex-direction: column;
    }
    span {
        color: purple;
        font-family: "Comic Sans MS", cursive;
        font-size: 1.5em;
    }
</style>

<main>
    <div id="app-container">
        <div>
            <h1>mindbody scheduler</h1>
            <div
                style="padding: 10px; display: flex; flex-direction: row; justify-content: space-between;">
                <span>server status:
                    {serverStatus && serverStatus.status}</span>
                <span>{timeNow}</span>
            </div>
        </div>
        <div id="opts">
            <div>
                <Login bind:serverStatus />
                <Scheduler />
            </div>
            <div>
                <Data />
            </div>
            <div>
                <RunHistory />
            </div>
        </div>
    </div>
</main>
