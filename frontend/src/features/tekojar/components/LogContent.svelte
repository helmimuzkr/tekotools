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
    let escapedLogHtml = `<p class=${serviceLogClass}>${escapeHtml(log)}</p>`;

    if (logState.searchQuery === "") return escapedLogHtml;
    const escapedQuery = escapeRegex(escapeHtml(logState.searchQuery));

    const matchPosition = logState.matchingLineIndices.indexOf(lineIndex);
    if (matchPosition === -1) return escapedLogHtml;

    const isCurrent = matchPosition === logState.currentMatchIndex;
    const hightlightClass = isCurrent ? "bg-orange-400 text-black" : "bg-yellow-300 text-black";

    escapedLogHtml = escapedLogHtml.replace(
      new RegExp(escapedQuery, "gi"),
      (wordMatch) => `<mark class="${hightlightClass}">${wordMatch}</mark>`,
    );

    return escapedLogHtml;
  }
</script>

<div class="flex-1 min-h-0 min-w-0 overflow-auto rounded border p-3 font-mono text-sm">
  {#each logState.selectedLogs as serviceLog, i}
    <div bind:this={lineRefs[i]} class="px-2 py-0.5">
      {@html highlightLine(serviceLog.log, i)}
    </div>
  {/each}
</div>
