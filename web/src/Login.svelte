<script lang="typescript">
    export let serverStatus;
    import { user } from "./store";
    import { onMount } from "svelte";

    onMount(async () => {
        console.log("login mounted");
        try {
            await fetch(`http://0.0.0.0:8888/status`)
                .then((r) => r.json())
                .catch((e) => {
                    serverStatus = "NOT ok";
                })
                .then((data) => {
                    serverStatus = data;
                    console.log(data);
                    return data;
                })
                .catch((e) => console.log(e));
        } catch (e) {
            console.log(e);
        }
    });
</script>

<style>
    .comp-box {
        padding: 25px;
    }
    .content {
        padding-top: 25px;
        display: grid;
        grid-template-columns: 20% 40%;
        grid-column-gap: 10px;
    }
    .subtitle {
        font-size: 14px;
        font-weight: 300;
    }
    .header {
        text-transform: capitalize;
    }
    .form-controller {
        justify-content: space-evenly;
        display: inline-grid;
    }
</style>

<!-- <div class="comp-box"> -->
<div class="w3-card-2 tbl-outln comp-box">
    <h2 class="header">mindbody login info</h2>
    <span class="subtitle">
        Enter your Mindbody username, password, and full name (as it appears on
        your account. e.g. - `Alexander Karlis`)
    </span>

    <div class="content">
        <form>
            <div class="form-group form-controller">
                <label for="fname">Full Name</label>
                <input
                    class="mui-form"
                    label="Full Name"
                    id="fname"
                    bind:value={$user.name} />
            </div>
            <div class="form-group form-controller">
                <label for="uname">User Email</label>
                <input
                    class="mui-form"
                    label="Username"
                    id="uname"
                    bind:value={$user.username} />
            </div>
            <div class="form-group form-controller">
                <label class="form-check-label" for="pword">Password</label>
                <input
                    class="mui-form"
                    label="Password"
                    type="password"
                    id="pword"
                    bind:value={$user.password} />
            </div>
        </form>
    </div>
</div>
