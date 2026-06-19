// logState.svelte.ts
import type { Log } from "@/features/tekojar/types"
import { EventsOn } from "@/../wailsjs/runtime/runtime"
import { serviceState } from "./serviceState.svelte"

const MAX_LINES = 1000
const TRIM_RATIO = 0.2

class LogState {
  serviceLogs = $state<Record<string, Log[]>>({})
  error = $state<string | null>(null)
  stoppingProcess = $state<boolean>(false)
  searchQuery = $state<string>("")
  currentMatchIndex = $state<number>(0)

  selectedLogs = $derived(
    this.serviceLogs[serviceState.selectedServiceId ?? ""] ?? []
  )

  matchingLineIndices = $derived(
    logState.searchQuery === "" || !serviceState.selectedServiceId
      ? []
      : this.selectedLogs
        .map((l, i) => l.log.toLowerCase().includes(logState.searchQuery.toLowerCase()) ? i : -1)
        .filter((i) => i !== -1)
  )

  currentMatchNumber = $derived(logState.matchingLineIndices.length > 0 ? logState.currentMatchIndex + 1 : 0);

  #logBufferByService: Record<string, number> = {}

  initLogListener() {
    EventsOn("service:log", (data: { id: string; logView: Log }) => {
      if (data.logView.log_type === "TIMER") {
        this.#replaceTimerLog(data.id, data.logView)
      } else {
        this.#appendLog(data.id, data.logView)
      }
    })
  }

  clearLog(id: string) {
    if (!this.serviceLogs[id] || !this.#logBufferByService[id]) return
    this.serviceLogs = { ...this.serviceLogs, [id]: [] }
    this.#logBufferByService[id] = 0
  }

  clearLogAllService() {
    if (!this.serviceLogs || !this.#logBufferByService) return
    Object.keys(this.serviceLogs).forEach(key => {
      this.serviceLogs[key] = [];
      this.#logBufferByService[key] = 0
    });
  }

  setSearchQuery(query: string) {
    this.searchQuery = query
    this.currentMatchIndex = 0
  }

  nextMatch() {
    if (this.matchingLineIndices.length === 0) return
    this.currentMatchIndex = (this.currentMatchIndex + 1) % this.matchingLineIndices.length
  }

  prevMatch() {
    if (this.matchingLineIndices.length === 0) return
    this.currentMatchIndex = (this.currentMatchIndex - 1) % this.matchingLineIndices.length
  }

  #replaceTimerLog(id: string, entry: Log) {
    const logs = this.serviceLogs[id] ?? []
    const withoutTimer = logs.filter((l) => l.log_type !== "TIMER")
    this.serviceLogs = { ...this.serviceLogs, [id]: [...withoutTimer, entry] }
  }

  #appendLog(id: string, entry: Log) {
    if (!this.serviceLogs[id]) {
      this.serviceLogs[id] = []
      this.#logBufferByService[id] = 0
    }

    const logs = this.serviceLogs[id]
    logs.push(entry)

    this.#logBufferByService[id] = (this.#logBufferByService[id] ?? 0) + 1

    if (this.#logBufferByService[id] > MAX_LINES) {
      const dropCount = Math.floor(logs.length * TRIM_RATIO)
      logs.splice(0, dropCount)
      this.#logBufferByService[id] -= dropCount
    }
  }
}

export const logState = new LogState()
