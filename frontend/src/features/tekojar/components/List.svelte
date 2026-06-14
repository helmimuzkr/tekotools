<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { toast } from "svelte-sonner";
  import { logState } from "../core/logState.svelte";
  import { serviceState } from "../core/serviceState.svelte";

  async function handleStartAll() {
    logState.clearLogAllService();
    const errs = await serviceState.startAll();
    if (errs) errs.forEach((e) => toast.error(e));
  }

  async function handleStopAll() {
    const errs = await serviceState.stopAll();
    if (errs) errs.forEach((e) => toast.error(e));
  }

  async function handleRestartAll() {
    logState.clearLogAllService();
    const errs = await serviceState.restartAll();
    if (errs) errs.forEach((e) => toast.error(e));
  }
</script>

<div class="flex flex-col h-full">
  <div class="flex-1 overflow-x-auto">
    <nav class="flex flex-col gap-1 min-w-max w-52 h-full border-r p-2">
      <div class="text-s text-muted-foreground px-2 py-1">Tekojar</div>
      {#each serviceState.services as service (service.id)}
        <Button
          class="w-full justify-start {serviceState.selectedServiceId === service.id ? 'bg-accent' : ''}"
          size={"sm"}
          variant="ghost"
          onclick={() => (serviceState.selectedServiceId = service.id)}
        >
          <div class={service.status === "ACTIVE" ? "" : "text-muted-foreground"}>{service.name}</div>
        </Button>
      {/each}
    </nav>
  </div>

  <!-- start all / stop all -->
  <div class="flex flex-col gap-2 p-2 border-t">
    <Button variant="secondary" size="sm" class="w-full" onclick={handleStartAll}>Start All</Button>
    <Button variant="secondary" size="sm" class="w-full" onclick={handleStopAll}>Stop All</Button>
    <Button variant="secondary" size="sm" class="w-full" onclick={handleRestartAll}>Restart All</Button>
  </div>
</div>
