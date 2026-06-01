import { derived, writable, get } from "svelte/store";
import type { Service } from "./type";

import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetAll, Start, Stop, Get as GetById } from "../../wailsjs/go/app/TekojarApp"

export const services = writable<Service[]>([]);
export const selectedServiceId = writable<string | null>(null);

export const selectedService = derived([services, selectedServiceId], ([$services, $selectedServiceId]) => {
  return $services.find(s => s.id === $selectedServiceId) ?? null
})

const logListeners = new Map<string, () => void>();

export async function getAllServices() {
  try {
    const result = await GetAll();
    services.set(result);
  } catch (err) {
    console.log(err);
  }
}

export async function getServices(id: string) {
  try {
    const result = await GetById(id);
    return result
  } catch (err) {
    console.log(err);
  }
}

export async function startService(id: string) {
  try {
    await Start(id);
    subscribeServiceLogs(id)
    await getServices(id).then(res => {
      services.update(curr =>
        curr.map(s => s.id === id ? res : s)
      )
    })
  } catch (err) {
    console.log(err);
    stopService(id)
  }

}

export async function stopService(id: string) {
  await Stop(id);
  unsubscribeServiceLogs(id)
  await getServices(id).then(res => {
    services.update(curr =>
      curr.map(s => s.id === id ? res : s)
    )
  })
}

export async function restartService(id: string) {
  await stopService(id)
  await startService(id)
}

export async function startAllService() {
  for (const service of get(services)) {
    startService(service.id)
  }
}

export async function stopAllService() {
  for (const service of get(services)) {
    stopService(service.id)
  }
}

export async function restartAllService() {
  stopAllService()
  startAllService()
}

// EventsOn registers a listener for the "service:log" event from Go. Returns an unlisten function to stop listening.
export function subscribeLogs(id: string) {
  const unlisten = EventsOn("service:log", (data: any) => {
    if (data.id === id) {
      services.update(current =>
        current.map(s => s.id === id
          ? { ...s, logs: [...s.logs, data.logView] } as Service
          : s
        )
      );
    }
  });
  return unlisten; // cleanup function
}

export function subscribeServiceLogs(id: string) {
  if (logListeners.has(id)) return; // already subscribed

  const unlisten = subscribeLogs(id);
  logListeners.set(id, unlisten);
}

export function unsubscribeServiceLogs(id: string) {
  const unlisten = logListeners.get(id);
  if (unlisten) {
    unlisten();
    logListeners.delete(id);
  }
}


