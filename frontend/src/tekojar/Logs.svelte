<script lang="ts">
  import { selectedService } from "./store";
</script>

<div class="flex-1 min-h-0 min-w-0 overflow-auto rounded border p-3 font-mono text-sm">
  {#each $selectedService?.logs ?? [] as serviceLog}
    <p class="whitespace-pre {serviceLog.is_error ? 'text-destructive' : ''}">{serviceLog.log}</p>
  {/each}
</div>
