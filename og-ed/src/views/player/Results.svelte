<script lang="ts">
  import type { ResultPacket } from "src/service/net";


export let levelResults:ResultPacket

let result=levelResults?.result || []


</script>
{#if result.length > 0}
  <div class="overlay">
    <div class="bg-white w-1/2 p-8 rounded-lg mx-auto shadow-lg">
      <h2 class="text-xl font-semibold text-center mb-4">Choose Word</h2>
      <div>
        {#each result as player}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
          
          <div class="flex gap-2 w-full">
            <p >{player.name}</p>
            <p >{player.points}</p>
          </div>


        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  /* Style for fullscreen overlay */
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 9999;
  }
</style>
