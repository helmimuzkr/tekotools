<script lang="ts">
  import { selectedServiceLogs } from "./store";
</script>

<div class="flex-1 min-h-0 min-w-0 overflow-auto rounded border p-3 font-mono text-sm">
  {#each $selectedServiceLogs as serviceLog}
    <p class="whitespace-pre {serviceLog.log_type === 'ERROR' ? 'text-destructive' : ''}">{serviceLog.log}</p>
  {/each}
</div>
