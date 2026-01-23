<script setup lang="ts">
import { computed } from "vue";
import { t, type Lang } from "../../i18n";

type BatchAction = "start" | "stop" | "restart" | "delete" | "autostart";

const props = defineProps<{
  selectedCount: number;
  totalCount: number;
  allSelected: boolean;
  lang: Lang;
}>();

const emit = defineEmits<{
  (event: "toggle-all", checked: boolean): void;
  (event: "action", action: BatchAction): void;
}>();

const disabled = computed(() => props.selectedCount === 0);

const actions: { id: BatchAction; labelKey: string; tone: string }[] = [
  { id: "start", labelKey: "action.start", tone: "text-emerald-600" },
  { id: "stop", labelKey: "action.stop", tone: "text-slate-600" },
  { id: "restart", labelKey: "action.restart", tone: "text-amber-600" },
  { id: "autostart", labelKey: "action.autostart", tone: "text-sky-600" },
  { id: "delete", labelKey: "action.delete", tone: "text-rose-600" }
];

const toggleAll = (event: Event) => {
  const target = event.target as HTMLInputElement | null;
  emit("toggle-all", target?.checked ?? false);
};

const runAction = (action: BatchAction) => {
  if (disabled.value) return;
  emit("action", action);
};
</script>

<template>
  <div class="space-y-3 pb-4">
    <div class="flex items-center justify-between text-[10px] uppercase tracking-[0.35em] text-fg/40">
      <label class="flex items-center gap-2">
        <input
          data-testid="toggle-all"
          type="checkbox"
          class="h-3 w-3 rounded border-border/60 bg-card/70 text-primary"
          :checked="props.allSelected"
          @change="toggleAll"
        />
        <span>{{ t("batch.all", props.lang) }}</span>
      </label>
      <span>{{ props.selectedCount }} / {{ props.totalCount }} {{ t("batch.selected", props.lang) }}</span>
    </div>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="action in actions"
        :key="action.id"
        :data-testid="`batch-${action.id}`"
        type="button"
        class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] transition"
        :class="[
          action.tone,
          disabled ? 'cursor-not-allowed opacity-40' : 'hover:border-primary/50'
        ]"
        :disabled="disabled"
        @click="runAction(action.id)"
      >
        {{ t(action.labelKey, props.lang) }}
      </button>
    </div>
  </div>
</template>
