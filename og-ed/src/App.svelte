<script lang="ts">


  import Router from "svelte-spa-router";
  import NotFound from "./views/NotFound.svelte";
  import Dashboard from "./views/Dashboard/Dashboard.svelte";


  let routes={
    '/':Dashboard,
    // '/play':PlayerView,
    // '/host':HostView,
    "*" :NotFound
  }
</script>
<div class="w-screen h-screen flex flex-col bg-gray-50 m-0 pb-2">
  <div class="w-full h-[70px] z-50 sticky top-0 bg-white shadow-md items-center flex px-2 justify-between">
    <div class="w-[50px] h-[50px] border-[0.5px] rounded-full bg-blue-200 justify-center items-center flex">Logo</div>
    <div class="flex-1 items-center justify-center flex">
    <h1 class="text-3xl font-bold">
      Scribble
    </h1>
    
    </div>
  </div>
  <div class="flex  flex-col w-full h-full">
    <!-- <div class="w-[10%] bg-white h-full z-40 shadow-sm">sidebar</div> -->
    <div class="flex flex-col flex-1 bg-gray-200">
      <Router {routes}/>
    </div>
    <div class="w-full h-[70px] z-50 sticky top-0 bg-white shadow-md items-center flex px-2 justify-center">....</div>

  
  </div>
  


</div>

<!-- {#if state==-1}
<Button on:click={getQuizes}>Get Quizes</Button>

<div>
  {#each quizzes as quiz}
    <QuizCard {quiz} on:host={() => hostQuiz(quiz)} />
  {/each}
</div>

<input bind:value={code} class="border" type="text" placeholder="game code" />
<input bind:value={name} class="border" type="text" placeholder="game name" />

<Button on:click={connect}>Join Game</Button>

{#if currentQuestion != null}
  <h2 class="text-3xl font-bold mt-8">{currentQuestion.name}</h2>
  <div class="flex">
    {#each currentQuestion.choices as choice}
      <div class="flex-1 bg-blue-400 text-center font-bold text-2xl">
        {choice.name}
      </div>
    {/each}
  </div>
{/if}
{:else if state==GameState.Lobby}
{#if host}
<Button on:click={startGame}>Start Game</Button>
<p>lobby state</p>
{#each players as player }
    <p>
      {player.name}
    </p>
  {/each
}

{/if}
{:else if state===GameState.Play}
{#if host}
Clock: {tick}
{#if currentQuestion != null}
  <h2 class="text-3xl font-bold mt-8">{currentQuestion.name}</h2>
  <div class="flex">
    {#each currentQuestion.choices as choice}
      <div class="flex-1 bg-blue-400 text-center font-bold text-2xl">
        {choice.name}
      </div>
    {/each}
  </div>
{/if}
{:else}
<p>Press correct answer</p>
{/if}
{/if} -->