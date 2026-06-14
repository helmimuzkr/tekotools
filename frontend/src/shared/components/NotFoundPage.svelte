<script lang="ts">
  import * as Empty from "$lib/components/ui/empty/index.js";
  import { Spinner } from "$lib/components/ui/spinner/index.js";
</script>

<div class="flex items-center justify-center h-screen">
  <Empty.Root class="w-full  md:p-6">
    <Empty.Header>
      <Spinner />
      <Empty.Description>
        The page you're looking for doesn't exist.
      </Empty.Description>
    </Empty.Header>
  </Empty.Root>
</div>
