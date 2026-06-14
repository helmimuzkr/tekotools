<script lang="ts">
  import Button from "$lib/components/ui/button/button.svelte";
  import { onMount, type Component } from "svelte";
  import Tekojar from "./components/Tekojar.svelte";
  import { settingState } from "./core/settingState.svelte";

  interface SettingPage {
    name: string;
    label: string;
    settingComponent: Component;
  }

  const settings: SettingPage[] = [
    {
      name: "tekojar_setting",
      label: "Tekojar",
      settingComponent: Tekojar,
    },
  ];

  onMount(() => {
    settingState.selectedSettingName = settings[0].name;
  });
</script>

<div class="flex h-full">
  <nav class="flex flex-col w-52 h-full border-r p-2 gap-1">
    <p class="text-s text-muted-foreground px-2 py-1">Settings</p>
    {#each settings as setting (setting.name)}
      <Button
        class="w-full justify-start {settingState.selectedSettingName === setting.name ? 'bg-accent' : ''}"
        size={"sm"}
        variant="ghost"
        onclick={() => (settingState.selectedSettingName = setting.name)}
      >
        {setting.label}
      </Button>
    {/each}
  </nav>

  <div class="flex-1 p-4">
    {#each settings as setting (setting.name)}
      <setting.settingComponent />
    {/each}
  </div>
</div>
