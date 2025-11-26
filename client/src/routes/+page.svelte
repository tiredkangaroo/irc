<script lang="ts">
  import { connection } from "$lib/stores/connection";

  let nicknameInput = $state("");
  let realnameInput = $state("");
  let hostInput = $state("");
  let portInput = $state("");
</script>

{#if $connection === null}
  <div class="w-full h-full flex justify-center items-center">
    <div class="min-w-fit w-1/3 min-h-fit h-1/2 py-6 mx-2 rounded-3xl">
      <h1 class="text-center">Connect to an IRC Server</h1>
      <div
        class="w-full h-full flex flex-col justify-center items-center gap-6"
      >
        <div class="w-full px-4 gap-2 flex flex-col">
          <div class="w-full">
            <label for="nickname" class="text-[#384d32]">Nickname</label>
            <input
              type="text"
              id="nickname"
              class="w-1/2"
              bind:value={nicknameInput}
              placeholder="Enter a nickname"
            />
          </div>
          <div class="w-full">
            <label for="realname" class="text-[#384d32]">Real Name</label>
            <input
              type="text"
              class="w-[80%]"
              bind:value={realnameInput}
              placeholder="Real Name"
            />
          </div>
          <div class="w-full">
            <label for="server" class="text-[#384d32]">Server</label>
            <div class="w-full flex flex-row">
              <input
                type="text"
                class="w-1/2"
                bind:value={hostInput}
                placeholder="Server Host"
              />
              <input
                class="w-1/4 ml-1"
                type="number"
                bind:value={portInput}
                placeholder="Port (e.g 6667)"
              />
            </div>
          </div>
        </div>
        <button
          class="bg-[#dbffd1] text-black font-bold py-2 px-4 rounded"
          onclick={() => {
            localStorage.setItem(
              "connection",
              JSON.stringify({
                nickname: nicknameInput,
                realname: realnameInput,
                hostport: hostInput + ":" + portInput,
              })
            );
            window.location.replace("/chat");
          }}>Connect</button
        >
      </div>
    </div>
  </div>
{:else}{/if}
