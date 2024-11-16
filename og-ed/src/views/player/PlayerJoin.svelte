<script lang="ts">
  import { createEventDispatcher } from "svelte";
    import Button from "../../lib/Button.svelte";
    import type { PlayerGame } from "../../service/player/player";
  let code:string=""
  let name:string=""
  export let game:PlayerGame;
    const dispatch=createEventDispatcher()
  function join(){
    dispatch("join")
      game.join(name,code) 
  }
  
  
  </script>
  
  
  <div class="bg-purple-400 min-h-screen w-full flex items-center justify-center ">
  <div>
      <h2 class="text-white font-bold text-4xl text-center first-letter:uppercase">quiz</h2>
      <div class="flex flex-col gap-2 mt-10 items-center ">
          <input type="text" bind:value={code} placeholder="Game Code" class="p-2 rounded"/>
          <input type="text" bind:value={name} placeholder="Name" class="p-2 rounded"/>
  
          <Button on:click={join}>
          Join Game
          </Button>
      </div>
  
      </div>
  </div>
  
  