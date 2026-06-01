<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { services, selectedServiceId, startAllService, stopAllService, restartAllService } from "./store";
</script>

<div class="flex flex-col h-full">
  <div class="flex-1 overflow-x-auto">
    <nav class="flex flex-col gap-1 min-w-max w-52 h-full border-r p-2">
      <p class="text-xs text-muted-foreground px-2 py-1">Tekojar</p>
      {#each $services as service, i (i)}
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
    <Button variant="default" size="sm" class="w-full" onclick={() => startAllService()}>Start All</Button>
    <Button variant="destructive" size="sm" class="w-full" onclick={() => stopAllService()}>Stop All</Button>
    <Button variant="outline" size="sm" class="w-full" onclick={() => restartAllService()}>Restart All</Button>
  </div>
</div>
