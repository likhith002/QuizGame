<script lang="ts">
  import type { Quiz } from "../../model/quiz";
  import { HostGame,gameState } from "../../service/host/host";
  import HostListView from "./HostListView.svelte";
   
 let game= new HostGame()
 
 console.log("GAME",game)
 let active=false 
 function onHost(event:CustomEvent){
    const quiz:Quiz=event.detail
    console.log(event)
    game.hostQuiz(quiz.id)
    active=true

}



</script>


<div>
    {#if active}
        <p>{$gameState}</p>
    {:else}
    <HostListView on:host={onHost}/>
        
    {/if}
</div>