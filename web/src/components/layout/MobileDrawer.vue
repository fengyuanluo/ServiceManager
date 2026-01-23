<script setup lang="ts">
import type { Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import NavList from "./NavList.vue";

const props = defineProps<{
  open: boolean;
  items: { key: ManagerKey; labelKey: string }[];
  active: ManagerKey;
  lang: Lang;
}>();

const emit = defineEmits<{
  (event: "close"): void;
  (event: "select", key: ManagerKey): void;
}>();
</script>

<template>
  <div
    data-testid="mobile-drawer"
    :data-state="props.open ? 'open' : 'closed'"
    class="fixed inset-0 z-40 lg:hidden"
    :class="props.open ? 'pointer-events-auto' : 'pointer-events-none'"
  >
    <div
      class="absolute inset-0 bg-black/30 transition-opacity duration-200"
      :class="props.open ? 'opacity-100' : 'opacity-0'"
      @click="emit('close')"
    ></div>
    <div
      class="absolute left-0 top-0 flex h-full w-72 flex-col border-r border-border/60 bg-bg/95 backdrop-blur transition-transform duration-200"
      :class="props.open ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="px-6 py-5">
        <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">SM</div>
        <div class="mt-3 h-px w-10 bg-border/70"></div>
      </div>
      <div class="flex-1 px-4 pb-6">
        <NavList
          :items="props.items"
          :active="props.active"
          :lang="props.lang"
          @select="
            (key) => {
              emit('select', key);
              emit('close');
            }
          "
        />
      </div>
    </div>
  </div>
</template>
