<script lang="ts">

  import Button from "./lib/Button.svelte";
  import QuizCard from "./lib/QuizCard.svelte";
  import { NetService, PacketTypes, type ChangeGameState, type PlayerJoinPacket } from "./service/net";
  import type { Player, Quiz, QuizQuestion } from "./model/quiz";

  let quizzes: { _id: string; name: string }[] = [];

  let netService = new NetService();
  let currentQuestion: QuizQuestion;
  let state=-1
  let host=false
  let players:Player[]=[]
  netService.connect();
  netService.onPacket((packet: any) => {

    switch (packet.id) {
      case 2: {
        currentQuestion = packet.question;
        break;
      }

      // change a packet to am aliased struct type with the custom interfaces 

      case PacketTypes.ChnageState:
        {
            let data=packet as ChangeGameState
            console.log(data.state);
            state=data.state

            break;
        }
      
      case PacketTypes.PlayerJoin:
        {
          console.log(packet)
          let data=packet as PlayerJoinPacket;
          players=[...players,data.player]
          break
        }
    }
  });
  let code = "";
  let name=""
  async function getQuizes() {
    try {
      const response = await fetch("http://localhost:5001/api/quizzes");

      let json = await response.json();
      console.log(json);
      quizzes = json;
    } catch (error) {
      console.log("Error",error);
    }
  }

  function connect() {
   netService.sendPacket({
    id:0,
    code:code,
    name:name
   })


  }
  // default mode for entering

  function hostQuiz(quiz: Quiz) {
    host=true
    netService.sendPacket({
      id: 1,
      quiz_id: quiz.id,
    });
  }

  function startGame(){

    netService.sendPacket({
      id:PacketTypes.StartGame
    })

  }
</script>
{#if state==-1}
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
{:else if state==0}
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
{/if}