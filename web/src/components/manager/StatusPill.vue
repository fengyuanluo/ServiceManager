<script setup lang="ts">
import { computed } from "vue";
import { t, type Lang } from "../../i18n";
import type { ServiceStatus } from "../../types/service";

const props = defineProps<{
  status: ServiceStatus;
  lang: Lang;
}>();

const toneClass = computed(() => {
  switch (props.status) {
    case "running":
      return "bg-emerald-500/10 text-emerald-600";
    case "stopped":
      return "bg-slate-500/10 text-slate-500";
    case "failed":
      return "bg-rose-500/10 text-rose-500";
    default:
      return "bg-amber-500/10 text-amber-600";
  }
});
</script>

<template>
  <span
    class="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] uppercase tracking-[0.25em]"
    :class="toneClass"
  >
    {{ t(`status.${props.status}`, props.lang) }}
  </span>
</template>
