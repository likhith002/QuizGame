<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { GameState, NetService, state } from "../../service/net";
  import HostView from "../host/HostView.svelte";
  import PlayerView from "../player/PlayerView.svelte";

  let netService: NetService = NetService.getInstance();

  onMount(() => {
    console.log("App started");
  });

  onDestroy(() => {
    netService.close();
  });

</script>

{#if $state === undefined}
  <HostView {netService} />
{:else if $state == GameState.Lobby}
  <h2>Wait for Players to join</h2>
{:else if $state == GameState.Wait}
  <h2>Wait for the player to choose a word</h2>
{:else}
  <PlayerView {netService} />
{/if}
