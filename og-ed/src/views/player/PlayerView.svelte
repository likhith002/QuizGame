<script lang="ts">
  import { PlayerGame } from "../../service/player/player";
  import PlayerJoin from "./PlayerJoin.svelte";


let game=new PlayerGame();
let active=false

function onjoin(){
    active=true
}

</script>

<div>
{#if active}
    active
{:else}
    <PlayerJoin on:join={onjoin} game={game}/>
    
{/if}
</div>


