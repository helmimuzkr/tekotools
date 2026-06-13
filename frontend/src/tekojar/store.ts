import { derived, writable, get } from "svelte/store";
import type { Service, Log } from "./type";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetAll, Start, Stop, Get as GetById } from "../../wailsjs/go/app/TekojarApp";

// CONST
const MAX_BYTES = 1 * 1024 * 1024; // 1MB
const TRIM_RATIO = 0.2; // if logs hit the max_bytes limit, drop oldest 20%
const logBufferSizeByService: Record<string, number> = {};

// STATE 
export const services = writable<Service[]>([]);
export const selectedServiceId = writable<string | null>(null);
export const serviceLogs = writable<Record<string, Log[]>>({});
export const error = writable<string | null>(null);
export const stoppingProcess = writable<boolean>(false);

// DERIVED
export const selectedService = derived(
  [services, selectedServiceId],
  ([$services, $selectedServiceId]) =>
    $services.find(s => s.id === $selectedServiceId) ?? null
);

export const selectedServiceLogs = derived(
  [serviceLogs, selectedServiceId],
  ([$serviceLogs, $selectedServiceId]) =>
    $selectedServiceId ? ($serviceLogs[$selectedServiceId] ?? []) : []
);

export const selectedServiceStatus = derived(
  [selectedService, stoppingProcess],
  ([$selectedService, $stoppingProcess]) => {
    if ($stoppingProcess) return "STOPPING";
    return $selectedService?.status ?? "INACTIVE";
  }
);

export const isDisableButton = derived(
  [selectedService, stoppingProcess],
  ([$selectedService, $stoppingProcess]) => ({
    start: $selectedService?.status === "ACTIVE" || $stoppingProcess,
    stop: $selectedService?.status !== "ACTIVE" || $stoppingProcess,
    restart: $stoppingProcess,
  })
);

// BACKEND API
async function refreshService(id: string) {
  const result = await GetById(id);
  services.update(currServices => currServices.map(s => s.id === id ? result : s));
}

export async function getAllServices() {
  try {
    const result = await GetAll();
    services.set(result);
    error.set(null);
  } catch (err) {
    error.set(err instanceof Error ? err.message : "Failed to fetch services");
  }
}

export async function startService(id: string) {
  try {
    await clearLog(id)
    await Start(id);
    await refreshService(id);
    error.set(null);
  } catch (err) {
    error.set(err instanceof Error ? err.message : "Failed to start service");
  }
}

export async function stopService(id: string) {
  try {
    stoppingProcess.set(true)
    await Stop(id);
    stoppingProcess.set(false)
    await refreshService(id);
    error.set(null);
  } catch (err) {
    error.set(err instanceof Error ? err.message : "Failed to stop service");
  }
}

export async function restartService(id: string) {
  await stopService(id);
  await startService(id);
}

export async function startAllServices() {
  for (const service of get(services)) {
    await startService(service.id);
  }
}

export async function stopAllServices() {
  for (const service of get(services)) {
    await stopService(service.id);
  }
}

export async function restartAllServices() {
  await stopAllServices();
  await startAllServices();
}

// HELPER
export function initLogListener() {
  EventsOn("service:log", (data: { id: string; logView: Log }) => {
    if (data.logView.log_type === 'TIMER') {
      replaceTimerLog(data.id, data.logView)
    } else {
      appendLog(data.id, data.logView)
    }
  });
}

export async function clearLog(id: string) {
  serviceLogs.update(currLogs => { return { ...currLogs, [id]: [] } })
}


// timer logs are replaced instead of appended. so only the latest countdown is shown
function replaceTimerLog(id: string, entry: Log) {
  serviceLogs.update(currLogs => {
    let selectedLog = currLogs[id] ?? [];
    let withoutTimer = selectedLog.filter(l => l.log_type !== 'TIMER')
    const updatedLog = [...withoutTimer, entry]
    return { ...currLogs, [id]: updatedLog };
  });
}

export function appendLog(id: string, entry: Log) {
  serviceLogs.update((currLogs) => {
    if (!currLogs[id]) {
      currLogs[id] = [];
      logBufferSizeByService[id] = 0;
    }

    const selectedLog = currLogs[id];

    selectedLog.push(entry);
    logBufferSizeByService[id] += entry.log.length;

    if (logBufferSizeByService[id] > MAX_BYTES) {
      const dropCount = Math.floor(selectedLog.length * TRIM_RATIO);
      const droppedBytes = selectedLog
        .slice(0, dropCount)
        .reduce((sum, e) => sum + e.log.length, 0);

      selectedLog.splice(0, dropCount);
      logBufferSizeByService[id] -= droppedBytes;
    }

    return currLogs;
  });
}
