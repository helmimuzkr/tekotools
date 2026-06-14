<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { serviceState } from "../core/serviceState.svelte";
  import { logState } from "../core/logState.svelte";
  import { toast } from "svelte-sonner";

  async function handleStart() {
    logState.clearLog(serviceState.selectedServiceId);
    const err = await serviceState.start(serviceState.selectedServiceId);
    if (err) toast.error(err);
  }

  async function handleStop() {
    logState.clearLog(serviceState.selectedServiceId);
    const err = await serviceState.stop(serviceState.selectedServiceId);
    if (err) toast.error(err);
  }

  async function handleRestart() {
    logState.clearLog(serviceState.selectedServiceId);
    const err = await serviceState.restart(serviceState.selectedServiceId);
    if (err) toast.error(err);
  }
</script>

<div class="flex justify-between items-center">
  <div class="flex items-center gap-2">
    <span class="font-bold truncate">{serviceState.selectedService.name}</span>
    <Badge variant={serviceState.selectedService.status === "ACTIVE" ? "default" : "secondary"}>
      {serviceState.selectedServiceStatus}
    </Badge>
  </div>
  <div class="flex gap-2">
    <Button variant="secondary" size="sm" onclick={() => logState.clearLog(serviceState.selectedServiceId)}
      >Clear</Button
    >
    <Button variant="secondary" size="sm" disabled={serviceState.isButtonDisabled.start} onclick={handleStart}
      >Start</Button
    >
    <Button variant="secondary" size="sm" disabled={serviceState.isButtonDisabled.stop} onclick={handleStop}
      >Stop</Button
    >
    <Button variant="secondary" size="sm" disabled={serviceState.isButtonDisabled.restart} onclick={handleRestart}
      >Restart</Button
    >
  </div>
</div>
