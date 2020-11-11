<script lang="typescript">
    export let serverStatus;
    import { user } from "./store";
    import { onMount } from "svelte";

    onMount(async () => {
        console.log("login mounted");
        await fetch(`http://0.0.0.0:8888/status`)
            .then((r) => r.json())
            .catch((e) => {
                serverStatus = "NOT ok";
            })
            .then((data) => {
                serverStatus = data;
                console.log(data);
                return data;
            });
    });
</script>

<style>
    .content {
        padding-top: 65px;
        display: grid;
        grid-template-columns: 20% 40%;
        grid-column-gap: 10px;
    }
    input {
        width: 300px;
    }
</style>

<form class="content">
    Full Name:
    <input type="text" id="fname" bind:value={$user.name} />
    Username:
    <input type="text" id="uname" bind:value={$user.username} />
    Password:
    <input type="password" id="pword" bind:value={$user.password} />
</form>
