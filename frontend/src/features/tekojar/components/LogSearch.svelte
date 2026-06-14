<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { logState } from "@/features/tekojar/core/logState.svelte";
  import { debounce } from "@/shared/utils";

  const debouncedSearch = debounce((value: string) => {
    logState.setSearchQuery(value);
  }, 300);
</script>

<div class="flex items-center gap-2 p-2 border-y">
  <Input
    placeholder="Search logs..."
    value={logState.searchQuery}
    oninput={(e) => debouncedSearch(e.currentTarget.value)}
    class="h-7 text-sm"
    onkeydown={(e) => {
      if (e.key === "Enter") {
        e.shiftKey ? logState.prevMatch() : logState.nextMatch();
      }
    }}
  />

  <span class="text-sm text-muted-foreground whitespace-nowrap">
    {logState.matchingLineIndices.length} matches
  </span>

  <Button variant="ghost" size="icon" onclick={() => logState.prevMatch()}>↑</Button>
  <Button variant="ghost" size="icon" onclick={() => logState.nextMatch()}>↓</Button>
</div>
