<script lang="ts">
  import TekojarList from "../tekojar/List.svelte";
  import TekojarLogs from "../tekojar/Logs.svelte";
  import TekojarControls from "../tekojar/Controls.svelte";

  import { services, selectedServiceId, getAllServices } from "./store";
  import { onMount } from "svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";

  // initialize services on mount component
  onMount(async () => {
    await getAllServices();
    selectedServiceId.set($services[0].id);
  });
</script>

<div class="flex h-full min-h-0">
  <TekojarList />
  <div class="flex flex-col flex-1 min-h-0 min-w-0 p-4 gap-3">
    {#if $selectedServiceId === null}
      <Skeleton class="h-4 w-full" />
    {:else}
      <TekojarControls />
      <TekojarLogs />
    {/if}
  </div>
</div>
