<script lang="ts">
  import { connection } from "$lib/stores/connection";

  if ($connection === null) {
    window.location.replace("/");
  }
</script>

{#if $connection}
  <p>
    connected as {$connection.nickname}, real name: {$connection.realname}, on
    server: {$connection.hostport}
  </p>
{/if}
