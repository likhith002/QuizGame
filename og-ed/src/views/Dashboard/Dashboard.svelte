<script lang="ts">
  import { onMount } from "svelte";
  import { NetService, state } from "../../service/net";
  import HostView from "../host/HostView.svelte";
  import PlayerPlay from "../player/PlayerPlay.svelte";

  let netService: NetService = NetService.getInstance();

  onMount(() => {
    console.log("App started");
  });
</script>

{#if $state === undefined || $state === 0}
  <HostView {netService} />
{:else}
  <PlayerPlay {netService} />
{/if}
