<script lang="ts">
  import { push } from "svelte-spa-router";
  import { state } from "../../service/player/player";

  import PlayerPlay from "./PlayerPlay.svelte";

  if (!$state) {
    push("/host");
  }
</script>

<div>
  <PlayerPlay />
</div>
