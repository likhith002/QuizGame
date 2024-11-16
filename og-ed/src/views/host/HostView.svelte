<script lang="ts">
  import App from "../../App.svelte";
  import type { Quiz } from "../../model/quiz";
  import { HostGame,gameState } from "../../service/host/host";
  import { GameState } from "../../service/net";
  import HostListView from "./HostListView.svelte";
  import HostLobby from "./HostLobby.svelte";
  import HostPlay from "./HostPlay.svelte";
   
 let game= new HostGame()
 
 console.log("GAME",game)
 let active=false 
 function onHost(event:CustomEvent){
    const quiz:Quiz=event.detail
    game.hostQuiz(quiz.id)
    active=true

}

let views:Record<GameState,any>={
    [GameState.Lobby]:HostLobby,
    [GameState.Play]:HostPlay
}



</script>


<div>
    {#if active}
        <svelte:component this={views[$gameState]} {game}/>
    {:else}
    <HostListView on:host={onHost}/>
        
    {/if}
</div>