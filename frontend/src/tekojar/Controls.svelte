<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { restartService, selectedService, startService, stopService } from "./store";
</script>

<div class="flex justify-between items-center">
  <div class="flex items-center gap-2">
    <span class="font-bold truncate">{$selectedService.name}</span>
    <Badge variant={$selectedService.status === "ACTIVE" ? "default" : "secondary"}>
      {$selectedService.status}
    </Badge>
  </div>
  <div class="flex gap-2">
    <Button
      variant="default"
      size="sm"
      disabled={$selectedService.status === "ACTIVE"}
      onclick={() => startService($selectedService.id)}>Start</Button
    >
    <Button
      variant="destructive"
      size="sm"
      disabled={$selectedService.status !== "ACTIVE"}
      onclick={() => stopService($selectedService.id)}>Stop</Button
    >
    <Button variant="outline" size="sm" onclick={() => restartService($selectedService.id)}>Restart</Button>
  </div>
</div>
