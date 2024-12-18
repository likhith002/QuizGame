<script lang="ts">
  import { getContext, onMount } from "svelte";
  import {
    currentPlayer,
    gameState,
    player,
    GameState,
    NetService,
    PacketTypes,
    players,
    type ChooseWordPacket,
    type DrawPoint,
    type ChangeGameState,
    type Player,
    type GameSettingsPacket,
    state,
    type ResultPacket,
  } from "../../service/net";
  import ChooseWord from "./ChooseWord.svelte";
  import PlayerPlay from "./PlayerPlay.svelte";
  import Results from "./Results.svelte";

  export let netService: NetService;
  let pointerEvents: boolean;
  let wordsToChoose: string[] = [];
  let childMethods: any;
  let levelResults: ResultPacket;

  function handleExpose(methods: any) {
    childMethods = methods;
  }

  function callDrawLine(data: DrawPoint) {
    childMethods.drawLine(
      data.x1,
      data.y1,
      data.x2,
      data.y2,
      data.color,
      data.lineWidth
    );
  }

  netService.onPacket((packet) => {
    switch (packet.id) {
      case PacketTypes.ChangeState: {
        let data = packet as ChangeGameState;
        gameState.set(data.state);

        if (data.state == GameState.UpdatePlayer) {
          const currPlayer = data.payload.player as Player;
          currentPlayer.set(currPlayer);

          if ($player.id) pointerEvents = currPlayer.id == $player.id;
        }

        break;
      }

      case PacketTypes.Coordinates: {
        const data = packet as DrawPoint;

        callDrawLine(data);
        break;
      }

      case PacketTypes.GameSettings: {
        const data = packet as GameSettingsPacket;

        players.update((prev) => [...prev, ...data.players]);

        if (data.currentPlayer)
          pointerEvents = data.currentPlayer.id == $player.id;
        if (data?.coordinates) {
          if (data.coordinates.length == 0) {
            childMethods.resetCanvas();
          } else {
            for (const point of data.coordinates) {
              callDrawLine(point);
            }
          }
        }
        break;
      }
      case PacketTypes.ChooseWord: {
        const data = packet as ChooseWordPacket;

        console.log("Choosing words....", data.words);
        wordsToChoose = data.words;

        break;
      }

      case PacketTypes.LevelResult: {
        console.log("got results", packet);

        levelResults = packet as ResultPacket;
      }
    }
  });
</script>

<div>
  <ChooseWord {wordsToChoose} {netService} />
  <Results {levelResults} />
  <PlayerPlay exposeMethods={handleExpose} {pointerEvents} {netService} />
</div>
