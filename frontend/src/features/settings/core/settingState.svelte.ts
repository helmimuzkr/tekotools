class SettingState {
  selectedSettingName = $state<string | null>(null);
  selectedChildSetting = $state<string | null>(null)
}

export const settingState = new SettingState()
