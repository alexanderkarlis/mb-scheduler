<script>
    import { Frequency, FrequencyArray } from "./store";

    console.log($FrequencyArray);
    function calculateScheduledTimes() {
        return [];
    }
</script>

<style>
    #container {
        margin-left: 40px;
        display: flex;
        flex-direction: column;
    }
</style>

<div id="container">
    {#each $FrequencyArray as freq}<span>{JSON.stringify(freq)}</span>{/each}
</div>
