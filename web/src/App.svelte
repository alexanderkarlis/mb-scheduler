<script lang="typescript">
    import Scheduler from "./Scheduler.svelte";
    import Data from "./Data.svelte";
    import Login from "./Login.svelte";

    let serverStatus = null;
</script>

<style>
    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }
    #opts {
        display: flex;
        flex-direction: row;
    }
</style>

<main>
    <h1>mindbody scheduler</h1>
    <span>server status: {serverStatus && serverStatus.status}</span>
    <div id="opts">
        <div>
            <Login bind:serverStatus />
            <Scheduler />
        </div>
        <Data />
    </div>
</main>
