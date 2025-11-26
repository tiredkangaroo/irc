<script lang="ts">
  import { connection } from "$lib/stores/connection";

  let connectionState = $state(0); // 0 = unknown, 1 = connecting, 2 = connected showing channels list, 3 = in channel
  let lastMessages = $state<Array<string>>([]);
  let channelList = $state<Array<string>>([]);

  connection.subscribe((value) => {
    if (value === null) {
      window.location.href = "/";
    }
  });

  let ws: WebSocket | null = $derived.by(() => {
    if (!$connection) return null;
    // where VITE_BACKEND_HOST is something like "localhost:8000"
    // however in prod, it would be empty string since the frontend is served from the same host
    console.log(
      "attempting websocket connection to backend",
      "ws://" +
        import.meta.env.VITE_BACKEND_HOST +
        "/proxy?irc_host=" +
        $connection.hostport
    );
    return new WebSocket(
      "ws://" +
        import.meta.env.VITE_BACKEND_HOST +
        "/proxy?irc_host=" +
        $connection.hostport
    );
  });

  $effect(() => {
    if (!ws || !$connection) return;
    ws.onopen = () => {
      console.log("websocket connected");
      connectionState = 1;
      ws.send("NICK " + $connection.nickname + "\r\n");
      ws.send(
        "USER " +
          $connection.nickname +
          " 0 * :" +
          $connection.realname +
          "\r\n"
      );
    };
    function handleMessage(event: MessageEvent) {
      lastMessages = [...lastMessages, event.data].slice(-50);
      const code = event.data.split(" ")[1];

      // ping event
      if (event.data.startsWith("PING")) {
        // PING :<token>
        ws!.send(event.data.replace("PING", "PONG")); // PONG :<token>
      }

      // error
      if (event.data.startsWith("ERROR")) {
        console.error("IRC ERROR:", event.data);
        ws!.close();
        connectionState = 0;
      }

      // welcome
      if (code === "001" && connectionState === 1) {
        // 001 is the RPL_WELCOME message (yall logged me in 🥹)
        connectionState = 2;
        // let's list channels
        ws!.send("LIST\r\n");
        // :irc.hackclub.com 322 beep #irs-ysws 3 :[+nt]
      }

      // channel list item
      if (code === "322") {
        const parts = event.data.split(" ");
        const channelName = parts[3];
        channelList = [...channelList, channelName];
      }
    }
    ws.onmessage = (event) => {
      console.log("received:", event.data);
      for (const line of event.data.split("\n")) {
        if (line.trim() !== "") {
          handleMessage(event);
        }
      }
    };
    ws.onerror = (event) => {
      console.error("websocket error:", event);
    };
    ws.onclose = () => {
      console.log("websocket closed");
    };
  });
</script>

{#if connectionState === 1}
  {@render connectingView(lastMessages)}
{:else if connectionState === 2}
  <div class="w-full h-full flex justify-center items-center">
    <div class="w-2/3 h-2/3 bg-gray-300 px-4 flex flex-col gap-4 rounded-xl">
      <h2 class="text-2xl font-bold text-center mt-4">Channels</h2>
      <div class="overflow-auto w-full flex flex-col gap-0.5">
        {#each channelList as channel}
          <p class="text-center">{channel}</p>
        {/each}
      </div>
    </div>
  </div>
{/if}

{#snippet connectingView(lastMessages: Array<string>)}
  <div class="w-full h-full flex justify-center items-center">
    <div class="w-2/3 h-2/3 bg-gray-300 px-4 flex flex-col gap-4 rounded-xl">
      <div class="flex flex-row align-middle items-center justify-between">
        <h2 class="text-2xl font-bold text-center mt-4">Connecting...</h2>
        <button
          class="bg-red-600 text-white ml-4 rounded-xl px-4 h-fit py-1"
          onclick={() => {
            localStorage.setItem("connection", null);
            window.location.replace("/");
          }}>Cancel</button
        >
      </div>
      <div class="overflow-auto w-full flex flex-col gap-0.5">
        {#each lastMessages as message}
          {#each message.split(" ").slice(3).join(" ").split("\n") as line}
            <p class="text-center">{line}</p>
          {/each}
        {/each}
      </div>
    </div>
  </div>
{/snippet}
