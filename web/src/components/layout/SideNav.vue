<script setup lang="ts">
import type { Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import NavList from "./NavList.vue";

const props = defineProps<{
  items: { key: ManagerKey; labelKey: string }[];
  active: ManagerKey;
  lang: Lang;
}>();

const emit = defineEmits<{
  (event: "select", key: ManagerKey): void;
}>();
</script>

<template>
  <aside
    class="hidden min-h-screen w-64 flex-col border-r border-border/60 bg-bg/80 backdrop-blur lg:flex"
  >
    <div class="px-6 py-6">
      <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">
        SM
      </div>
      <div class="mt-3 h-px w-10 bg-border/70"></div>
    </div>
    <div class="flex-1 px-4 pb-6">
      <NavList
        :items="props.items"
        :active="props.active"
        :lang="props.lang"
        @select="emit('select', $event)"
      />
    </div>
  </aside>
</template>
