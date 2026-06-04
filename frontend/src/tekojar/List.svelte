<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { services, selectedServiceId, startAllServices, stopAllServices, restartAllServices } from "./store";
</script>

<div class="flex flex-col h-full">
  <div class="flex-1 overflow-x-auto">
    <nav class="flex flex-col gap-1 min-w-max w-52 h-full border-r p-2">
      <div class="text-xs text-muted-foreground px-2 py-1 my-1">Tekojar</div>
      {#each $services as service (service.id)}
        <Button
          class="w-full justify-start {$selectedServiceId === service.id ? 'bg-accent' : ''}"
          size={"sm"}
          variant="ghost"
          onclick={() => selectedServiceId.set(service.id)}
        >
          <div class={service.status === "ACTIVE" ? "" : "text-muted-foreground"}>{service.name}</div>
        </Button>
      {/each}
    </nav>
  </div>

  <!-- start all / stop all -->
  <div class="flex flex-col gap-1 p-2 border-t">
    <Button variant="default" size="sm" class="w-full" onclick={() => startAllServices()}>Start All</Button>
    <Button variant="destructive" size="sm" class="w-full" onclick={() => stopAllServices()}>Stop All</Button>
    <Button variant="outline" size="sm" class="w-full" onclick={() => restartAllServices()}>Restart All</Button>
  </div>
</div>
