// serviceState.svelte.ts
import type { Service } from "@/features/tekojar/types"
import { GetAll, Start, Stop, Get as GetById } from "@/../wailsjs/go/app/TekojarApp"
import type { ErrorPromise, ErrorPromiseArr } from "@/type"

class ServiceState {
  services = $state<Service[]>([])
  selectedServiceId = $state<string | null>(null)
  stoppingProcess = $state<boolean>(false)

  selectedService = $derived(
    this.services.find(s => s.id === this.selectedServiceId) ?? null
  )

  selectedServiceStatus = $derived(
    this.stoppingProcess ? "STOPPING" : (this.selectedService?.status ?? "INACTIVE")
  )

  isButtonDisabled = $derived({
    start: this.selectedService?.status === "ACTIVE" || this.stoppingProcess,
    stop: this.selectedService?.status !== "ACTIVE" || this.stoppingProcess,
    restart: this.stoppingProcess,
  })

  selectFirstService() {
    if (serviceState.services.length > 0) {
      this.selectedServiceId = this.services[0].id;
    }
  }

  async getAll(): ErrorPromise {
    try {
      this.services = await GetAll()
    } catch (err) {
      return err instanceof Error ? err.message : "Failed to fetch services"
    }
  }

  async start(id: string): ErrorPromise {
    try {
      await Start(id)
      await this.refreshService(id)
    } catch (err) {
      return err instanceof Error ? `service-id: ${id} - ${err.message}` : `service-id: ${id} - Failed to start service`
    }
  }

  async stop(id: string): ErrorPromise {
    try {
      this.stoppingProcess = true
      await Stop(id)
      this.stoppingProcess = false
      await this.refreshService(id)
    } catch (err) {
      this.stoppingProcess = false
      return err instanceof Error ? `service-id: ${id} - ${err.message}` : `service-id: ${id} - Failed to stop service`
    }
  }

  async restart(id: string): ErrorPromise {
    let err = await this.stop(id)
    if (err) return err

    err = await this.start(id)
    if (err) return err

    return null
  }

  async startAll(): ErrorPromiseArr {
    let errs: string[] = []
    for (const service of this.services) {
      const err = await this.start(service.id)
      if (err) errs.push(err)
    }
    return errs.length > 1 ? errs : null
  }

  async stopAll(): ErrorPromiseArr {
    let errs: string[] = []
    for (const service of this.services) {
      const err = await this.stop(service.id)
      if (err) errs.push(err)
    }
    return errs.length > 1 ? errs : null
  }

  async restartAll(): ErrorPromiseArr {
    let err = await this.stopAll()
    if (err) return err

    err = await this.startAll()
    if (err) return err

    return null
  }

  private async refreshService(id: string) {
    const updated = await GetById(id)
    this.services = this.services.map(s => s.id === id ? updated : s)
  }
}

export const serviceState = new ServiceState()
