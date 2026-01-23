<script setup lang="ts">
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import ManagerIcon from "../icons/ManagerIcon.vue";

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
  <nav class="space-y-1">
    <button
      v-for="item in props.items"
      :key="item.key"
      type="button"
      class="group flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm transition"
      :class="
        props.active === item.key
          ? 'bg-muted/80 text-fg shadow-sm shadow-black/5'
          : 'text-fg/70 hover:bg-muted/60 hover:text-fg'
      "
      :aria-current="props.active === item.key ? 'page' : undefined"
      @click="emit('select', item.key)"
    >
      <ManagerIcon :name="item.key" class="text-fg/60 group-hover:text-fg" />
      <span class="flex-1 truncate">{{ t(item.labelKey, props.lang) }}</span>
      <span
        class="text-[10px] uppercase tracking-widest text-fg/30 group-hover:text-fg/60"
        aria-hidden="true"
      >
        {{ item.key }}
      </span>
    </button>
  </nav>
</template>
