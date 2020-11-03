<script lang="typescript">
    export let serverStatus;
    import { user } from "./store";
    import { onMount } from "svelte";

    async function onSubmit() {
        console.log(JSON.stringify($user));
        let r = await fetch("http://0.0.0.0:8888", {
            method: "post",
            body: JSON.stringify($user),
        });
        r.text().then((data) => {
            console.log(JSON.parse(data));
        });
    }

    onMount(async () => {
        console.log("login mounted");
        await fetch(`http://0.0.0.0:8888/`)
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
