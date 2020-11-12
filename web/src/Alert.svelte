<script lang="typescript">
    import { MESSAGE_TYPES } from "./constants";
    import { PostAttemptMessage } from "./store";

    export let message = '';
    export let type: MESSAGE_TYPES;
    export let open;

    let show = open || false;
</script>

<style>
    #alert {
        position: fixed;
        top: 95%;
        left: 30%;
        width: 30%;
        text-align: center;
    }
</style>

<div>
    {#if $PostAttemptMessage}
        <div
            id="alert"
            class={MESSAGE_TYPES.FAILED === type ? 'alert alert-danger alert-dismissible fade show' : 'alert alert-success alert-dismissible fade show'}
            role="alert">
            <strong>{message}</strong>
            <button
                on:click={() => {
                    show = false;
                    $PostAttemptMessage = false;
                }}
                type="button"
                class="close"
                data-dismiss="alert"
                aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>
    {/if}
</div>
