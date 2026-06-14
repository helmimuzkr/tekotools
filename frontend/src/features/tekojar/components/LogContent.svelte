<script lang="ts">
  import { logState } from "@/features/tekojar/core/logState.svelte";
  import { escapeHtml } from "@/shared/utils";

  let lineRefs: HTMLDivElement[] = $state([]);
  const serviceLogClass = "whitespace-pre";

  $effect(() => {
    const currentLineIndex = logState.matchingLineIndices[logState.currentMatchIndex];
    if (currentLineIndex === undefined) return;
    lineRefs[currentLineIndex]?.scrollIntoView({ block: "center" });
  });

  function highlightLine(log: string, lineIndex: number): string {
    let escapedLog = `<p class=${serviceLogClass}>${escapeHtml(log)}</p>`;

    if (logState.searchQuery === "") return escapedLog;
    const escapedQuery = escapeHtml(logState.searchQuery);

    const matchPosition = logState.matchingLineIndices.indexOf(lineIndex);
    if (matchPosition === -1) return escapedLog;

    const isCurrent = matchPosition === logState.currentMatchIndex;
    const hightlightClass = isCurrent ? "bg-orange-400 text-black" : "bg-yellow-300 text-black";

    escapedLog = escapedLog.replace(
      new RegExp(escapedQuery, "gi"),
      (wordMatch) => `<mark class="${hightlightClass}">${wordMatch}</mark>`,
    );

    return escapedLog;
  }

  function isMatch(lineIndex: number): boolean {
    return logState.matchingLineIndices.includes(lineIndex);
  }
</script>

<div class="flex-1 min-h-0 min-w-0 overflow-auto rounded border p-3 font-mono text-sm">
  {#each logState.selectedLogs as serviceLog, i}
    <div bind:this={lineRefs[i]} class="px-2 py-0.5">
      {#if isMatch(i)}
        {@html highlightLine(serviceLog.log, i)}
      {:else}
        <p class={serviceLogClass}>{serviceLog.log}</p>
      {/if}
    </div>
  {/each}
</div>
