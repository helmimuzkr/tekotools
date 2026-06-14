import { logState } from "@/features/tekojar/core/logState.svelte";

export async function initTekojar() {
  logState.initLogListener()
}

