<script lang="ts">
  import { logState } from "@/features/tekojar/core/logState.svelte";
  import { escapeHtml, escapeRegex } from "@/shared/utils";

  let lineRefs: HTMLDivElement[] = $state([]);
  const serviceLogClass = "whitespace-pre";

  $effect(() => {
    const currentLineIndex = logState.matchingLineIndices[logState.currentMatchIndex];
    if (currentLineIndex === undefined) return;
    lineRefs[currentLineIndex]?.scrollIntoView({ block: "center" });
  });

  function highlightLine(log: string, lineIndex: number): string {
    let escapedLog = escapeHtml(log);

    if (logState.searchQuery === "") return createLogHtml(escapedLog);
    const escapedQuery = escapeRegex(escapeHtml(logState.searchQuery));

    const matchPosition = logState.matchingLineIndices.indexOf(lineIndex);
    if (matchPosition === -1) return createLogHtml(escapedLog);

    const isCurrent = matchPosition === logState.currentMatchIndex;
    const hightlightClass = isCurrent ? "bg-orange-400 text-black" : "bg-yellow-300 text-black";

    escapedLog = escapedLog.replace(
      new RegExp(escapedQuery, "gi"),
      (wordMatch) => `<mark class="${hightlightClass}">${wordMatch}</mark>`,
    );

    return createLogHtml(escapedLog);
  }

  function createLogHtml(log: string): string {
    return `<p class=${serviceLogClass}>${log}</p>`;
  }
</script>

<div class="flex-1 min-h-0 min-w-0 overflow-auto rounded border p-3 font-mono text-sm">
  {#each logState.selectedLogs as serviceLog, i}
    <div bind:this={lineRefs[i]} class="px-2 py-0.5">
      {@html highlightLine(serviceLog.log, i)}
    </div>
  {/each}
</div>
