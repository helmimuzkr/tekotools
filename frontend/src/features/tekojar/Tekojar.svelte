<script lang="ts">
  import List from "../tekojar/components/List.svelte";
  import Logs from "../tekojar/components/Logs.svelte";
  import Controls from "../tekojar/components/Controls.svelte";

  import { serviceState } from "@/features/tekojar/core/serviceState.svelte";
  import { onMount } from "svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";
  import { toast } from "svelte-sonner";

  // initialize services on mount component
  onMount(async () => {
    const err = await serviceState.getAll();
    if (err) toast.error(err);
    serviceState.selectFirstService();
  });
</script>

<div class="flex h-full min-h-0">
  <div class=" min-h-0 min-w-0 p-2">
    <List />
  </div>
  <div class="flex flex-col flex-1 min-h-0 min-w-0 p-2 gap-3">
    {#if serviceState.selectedServiceId === null}
      <Skeleton class="h-4 w-full" />
    {:else}
      <Controls />
      <Logs />
    {/if}
  </div>
</div>
