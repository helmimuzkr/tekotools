<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import {
    restartService,
    selectedService,
    startService,
    stopService,
    selectedServiceStatus,
    isDisableButton,
    clearLog,
  } from "./store";
</script>

<div class="flex justify-between items-center">
  <div class="flex items-center gap-2">
    <span class="font-bold truncate">{$selectedService.name}</span>
    <Badge variant={$selectedService.status === "ACTIVE" ? "default" : "secondary"}>
      {$selectedServiceStatus}
    </Badge>
  </div>
  <div class="flex gap-2">
    <Button variant="default" size="sm" onclick={() => clearLog($selectedService.id)}>Clear</Button>
    <Button
      variant="default"
      size="sm"
      disabled={$isDisableButton.start}
      onclick={() => startService($selectedService.id)}>Start</Button
    >
    <Button
      variant="destructive"
      size="sm"
      disabled={$isDisableButton.stop}
      onclick={() => stopService($selectedService.id)}>Stop</Button
    >
    <Button
      variant="outline"
      size="sm"
      disabled={$isDisableButton.restart}
      onclick={() => restartService($selectedService.id)}>Restart</Button
    >
  </div>
</div>
