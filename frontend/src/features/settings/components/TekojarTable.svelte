<script lang="ts">
  import * as Field from "$lib/components/ui/field/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import Button from "$lib/components/ui/button/button.svelte";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import * as Table from "$lib/components/ui/table";
  import { app } from "@/../wailsjs/go/models";
  import { tekojarSettingState } from "../core/tekojarSettingState.svelte";
  import { toast } from "svelte-sonner";

  let isError = $state(false);

  function handleAdd() {
    tekojarSettingState.addService();
  }

  function handleServiceNameChange(service: app.DTOServiceSetting, name: string) {
    service.name = name;
    isError = false;
    validateDuplicateServiceName(name);
  }

  function validateDuplicateServiceName(name: string) {
    const isDuplicate = tekojarSettingState.tekojarSetting.service_settings.filter((s) => s.name === name).length > 1;
    if (isDuplicate) {
      toast.error(`"${name}" already exists`);
      isError = true;
    }
  }
</script>

<Field.Set>
  <Field.Legend>Service</Field.Legend>
  <div class="flex flex-col gap-3">
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Name</Table.Head>
          <Table.Head>Path</Table.Head>
          <Table.Head>Skip</Table.Head>
          <Table.Head>Index</Table.Head>
          <Table.Head>Delay (s)</Table.Head>
          <Table.Head></Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each tekojarSettingState.tekojarSetting.service_settings as service, i (i)}
          <Table.Row>
            <Table.Cell>
              <Input
                value={service.name}
                placeholder="service.jar"
                oninput={(e) => handleServiceNameChange(service, e.currentTarget.value)}
                required
              />
            </Table.Cell>
            <Table.Cell>
              <Input bind:value={service.path} placeholder="/home/user/service" required />
            </Table.Cell>
            <Table.Cell>
              <Checkbox bind:checked={service.skip_flag} />
            </Table.Cell>
            <Table.Cell>
              <Input type="number" bind:value={service.idx} class="w-20" />
            </Table.Cell>
            <Table.Cell>
              <Input type="number" bind:value={service.delay} class="w-20" />
            </Table.Cell>
            <Table.Cell>
              <Button variant="destructive" size="sm" onclick={() => tekojarSettingState.removeService(service.name)}
                >Remove</Button
              >
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
    <Button variant="secondary" size="sm" disabled={isError} class="w-fit" onclick={handleAdd}>Add Service</Button>
  </div>
</Field.Set>
