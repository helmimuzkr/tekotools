import { GetSetting, SaveSetting } from "$wails/go/app/TekojarApp";
import type { TekojarSetting } from "../types";
import type { ErrorPromise } from "@/type";

class TekojarSettingState {
  tekojarSetting = $state<TekojarSetting | null>(null);

  async getSetting(): ErrorPromise {
    try {
      const result = await GetSetting();
      this.tekojarSetting = result;
      return null
    } catch (err) {
      return err instanceof Error ? err.message : "failed to fetch tekojar setting";
    }
  }

  async save(): ErrorPromise {
    const hasEmpty = this.tekojarSetting.service_settings.some((s) => !s.name || !s.path);
    if (hasEmpty) {
      return "All services must have a name and path";
    }
    await SaveSetting(this.tekojarSetting);
    return null
  }

  addService() {
    this.tekojarSetting.service_settings = [
      ...this.tekojarSetting.service_settings,
      { id: "", name: "", path: "", skip_flag: false, delay: 0, idx: 0 },
    ];
  }

  removeService(name: string) {
    this.tekojarSetting.service_settings = this.tekojarSetting.service_settings.filter((s) => s.name !== name);
  }

}

export const tekojarSettingState = new TekojarSettingState();
