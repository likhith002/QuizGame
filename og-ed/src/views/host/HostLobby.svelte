
<script lang="ts">

    import Button from "../../lib/Button.svelte";
  import { HostGame } from "../../service/host/host";
    export let game:HostGame;

    function start() {
        game.start()

    }

</script>

<div class="p-8">
    <div class="flex justify-end">
        <Button on:click={start}>Start Game</Button>
    </div>

</div>