<script lang="ts">
  import * as Field from "$lib/components/ui/field/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { tekojarSettingState } from "../core/tekojarSettingState.svelte";
</script>

<Field.Set>
  <Field.Legend>General</Field.Legend>
  <Field.Group>
    <Field.Field orientation="horizontal">
      <Field.Label class="text-sm w-24 font-normal">Command</Field.Label>
      <Input placeholder="java -jar $PATH" bind:value={tekojarSettingState.tekojarSetting.command} />
    </Field.Field>
  </Field.Group>
</Field.Set>
