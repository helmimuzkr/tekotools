import { derived, writable, get } from "svelte/store";
import type { Service, Log } from "./type";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetAll, Start, Stop, Get as GetById } from "../../wailsjs/go/app/TekojarApp";

// --- State ---
export const services = writable<Service[]>([]);
export const selectedServiceId = writable<string | null>(null);
export const serviceLogs = writable<Record<string, Log[]>>({});
export const error = writable<string | null>(null);

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


export function initLogListener() {
  EventsOn("service:log", (data: { id: string; logView: Log }) => {
    serviceLogs.update(curr => ({
      ...curr,
      [data.id]: [...(curr[data.id] ?? []), data.logView]
    }));
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
    await Stop(id);
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




