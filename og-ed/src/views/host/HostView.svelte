<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { push } from "svelte-spa-router";
  import { PlayerGame, state } from "../../service/player/player";
  import { GameState } from "../../service/net";

  let game = new PlayerGame();

  function joinGame() {
    game.join("Player1");
  }

  console.log("Current state", $state);

  $: if ($state == GameState.Lobby) {
    console.log("Player Added");
  } else if ($state == GameState.Play) {
    push("/play");
  }
</script>

<div class="w-fill h-full flex flex-col px-10 py-5">
  <div class="justify-end w-full flex">
    <Button>Host</Button>
  </div>
  <div class="w-fill h-full flex flex-col justify-center items-center">
    <Card.Root class="w-[350px]">
      <Card.Header>
        <Card.Title>Join Game</Card.Title>
        <Card.Description>Join an already happening game</Card.Description>
      </Card.Header>
      <Card.Content>
        <form>
          <div class="grid w-full items-center gap-4">
            <div class="flex flex-col space-y-1.5">
              <Label for="name">Name</Label>
              <Input id="name" placeholder="Enter your nick name" />
            </div>
          </div>
        </form>
      </Card.Content>
      <Card.Footer class="flex justify-between">
        <Button variant="outline">Cancel</Button>
        <Button on:click={joinGame}>Join</Button>
      </Card.Footer>
    </Card.Root>
  </div>
</div>
