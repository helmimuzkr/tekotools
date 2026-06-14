<script lang="ts">
  import { onMount } from "svelte";
  import { toast } from "svelte-sonner";
  import * as Field from "$lib/components/ui/field/index.js";
  import Button from "$lib/components/ui/button/button.svelte";
  import Skeleton from "$lib/components/ui/skeleton/skeleton.svelte";
  import TekojarGeneral from "./TekojarGeneral.svelte";
  import TekojarTable from "./TekojarTable.svelte";
  import { tekojarSettingState } from "../core/tekojarSettingState.svelte";

  onMount(async () => {
    const err = await tekojarSettingState.getSetting();
    if (err) toast.error(err);
  });

  async function save() {
    const err = await tekojarSettingState.save();
    if (!err) {
      toast.success("Settings saved");
    } else {
      toast.error(err);
    }
  }
</script>

<div class="w-full max-w-2xl">
  {#if tekojarSettingState.tekojarSetting}
    <form>
      <Field.Group>
        <TekojarGeneral />

        <Field.Separator />

        <TekojarTable />

        <Field.Field orientation="horizontal" class="justify-end">
          <Button variant="secondary" size="sm" class="px-5 py-2" onclick={save}>Save</Button>
        </Field.Field>
      </Field.Group>
    </form>
  {:else}
    <Skeleton />
  {/if}
</div>
