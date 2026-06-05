import { derived, writable, get } from "svelte/store";
import type { Service, Log } from "./type";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetAll, Start, Stop, Get as GetById } from "../../wailsjs/go/app/TekojarApp";

// --- State ---
export const services = writable<Service[]>([]);
export const selectedServiceId = writable<string | null>(null);
export const serviceLogs = writable<Record<string, Log[]>>({});
export const error = writable<string | null>(null);
export const stoppingProcess = writable<boolean>(false);

// --- Derived ---
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

export function initLogListener() {
  EventsOn("service:log", (data: { id: string; logView: Log }) => {
    serviceLogs.update(curr => {
      let existing = curr[data.id] ?? [];

      // timer logs are replaced instead of appended. so only the latest countdown is shown
      if (data.logView.log_type === 'TIMER') {
        existing = existing.filter(l => l.log_type !== 'TIMER')
      }

      const updated = [...existing, data.logView]

      return { ...curr, [data.id]: updated };
    });
  });
}

// --- Internal helper ---
async function refreshService(id: string) {
  const result = await GetById(id);
  services.update(curr => curr.map(s => s.id === id ? result : s));
}

// --- Public API ---
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




