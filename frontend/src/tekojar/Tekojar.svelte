<script lang="ts">
  import TekojarList from '../tekojar/List.svelte';
  import TekojarLogs from '../tekojar/Logs.svelte';
  import TekojarControls from '../tekojar/Controls.svelte';

  import { services, selectedServiceId, getAllServices } from './store';
  import { onMount } from 'svelte';
  import Skeleton from '$lib/components/ui/skeleton/skeleton.svelte';

  // initialize services on mount component
  onMount(async () => {
    await getAllServices();
    selectedServiceId.set($services[0].id);
  });
    $inspect($services)
</script>

<div class="flex h-full">
  <TekojarList />
  <div class="flex flex-col flex-1 p-4 gap-3">
    {#if $selectedServiceId === null}
      <Skeleton class="h-4 w-full" />
    {:else}
      <TekojarControls />
      <TekojarLogs />
    {/if}
  </div>
</div>
