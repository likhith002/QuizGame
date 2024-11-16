<script lang="ts">
  import type { Player } from "../../model/quiz";


export let player:Player
</script>

<button on:click>
    <h3 class="text-3xl text-white px-4 py-2 rounded-full hover:line-through">
        {player.name}
    </h3>
</button>