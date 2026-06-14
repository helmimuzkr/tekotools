<!-- JsonataQuery.svelte -->
<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Textarea } from "$lib/components/ui/textarea";
  import { Input } from "$lib/components/ui/input";
  import { ScrollArea } from "$lib/components/ui/scroll-area";

  let json = $state("");
  let expr = $state("");
  let result = $state("");
</script>

<div class="flex h-full p-4 gap-4">
  <Textarea class="flex-1 font-mono" placeholder="JSON" bind:value={json} />
  <div class="flex flex-col flex-1 gap-3">
    <Input class="font-mono" placeholder="$.users" bind:value={expr} />
    <ScrollArea class="flex-1 rounded border font-mono text-sm p-3">
      <p>{result || "// run query to see result"}</p>
    </ScrollArea>
    <Button>Run</Button>
  </div>
</div>
