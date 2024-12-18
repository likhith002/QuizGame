<script lang="ts">
  import {
    type NetService,
    type SelectedWord,
    PacketTypes,
  } from "../../service/net";

  export let wordsToChoose:string[];
  export let netService: NetService;
  function sendSelectedWord(word: string) {
    let packet: SelectedWord = {
      id: PacketTypes.SelectedWord,
      word: word,
    };

    netService.sendPacket(packet);
    wordsToChoose=[]
  }
</script>

{#if wordsToChoose.length > 0}
  <div class="overlay">
    <div class="bg-white w-1/2 p-8 rounded-lg mx-auto shadow-lg">
      <h2 class="text-xl font-semibold text-center mb-4">Choose Word</h2>
      <div>
        {#each wordsToChoose as word}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
          <li on:click={() => sendSelectedWord(word)}>{word}</li>
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
