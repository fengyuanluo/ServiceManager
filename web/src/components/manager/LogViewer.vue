<script setup lang="ts">
import { computed, ref } from "vue";
import { t, type Lang } from "../../i18n";

type LogLevel = "all" | "info" | "warn" | "error";

const props = defineProps<{
  lines: string[];
  lang: Lang;
  active: boolean;
}>();

const emit = defineEmits<{
  (event: "download", lines: string[]): void;
}>();

const limitOptions = [50, 100, 200, 500];
const limit = ref(200);
const keyword = ref("");
const level = ref<LogLevel>("all");

const filteredLines = computed(() => {
  const slice = props.lines.slice(-limit.value);
  const loweredKeyword = keyword.value.trim().toLowerCase();
  const byLevel =
    level.value === "all"
      ? slice
      : slice.filter((line) => line.toLowerCase().includes(level.value));
  if (!loweredKeyword) return byLevel;
  return byLevel.filter((line) => line.toLowerCase().includes(loweredKeyword));
});

const downloadLogs = () => {
  emit("download", filteredLines.value);
  if (typeof window === "undefined") return;
  const blob = new Blob([filteredLines.value.join("\n")], {
    type: "text/plain"
  });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = "service.log";
  anchor.click();
  URL.revokeObjectURL(url);
};
</script>

<template>
  <div class="space-y-4">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <div class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
        {{ t("layout.logs", props.lang) }}
      </div>
      <div class="flex flex-wrap items-center gap-2 text-xs text-fg/60">
        <label class="flex items-center gap-2">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("logs.limit", props.lang) }}
          </span>
          <select
            v-model.number="limit"
            class="rounded-full border border-border/70 bg-card/70 px-2 py-1 text-xs"
          >
            <option v-for="option in limitOptions" :key="option" :value="option">
              {{ option }}
            </option>
          </select>
        </label>
        <label class="flex items-center gap-2">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("logs.level", props.lang) }}
          </span>
          <select
            v-model="level"
            class="rounded-full border border-border/70 bg-card/70 px-2 py-1 text-xs"
          >
            <option value="all">{{ t("logs.levelAll", props.lang) }}</option>
            <option value="info">{{ t("logs.levelInfo", props.lang) }}</option>
            <option value="warn">{{ t("logs.levelWarn", props.lang) }}</option>
            <option value="error">{{ t("logs.levelError", props.lang) }}</option>
          </select>
        </label>
        <input
          data-testid="log-search"
          v-model="keyword"
          type="search"
          class="rounded-full border border-border/70 bg-card/70 px-3 py-1 text-xs"
          :placeholder="t('logs.search', props.lang)"
        />
        <button
          data-testid="log-download"
          type="button"
          class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
          @click="downloadLogs"
        >
          {{ t("logs.export", props.lang) }}
        </button>
      </div>
    </header>

    <div class="rounded-2xl border border-border/60 bg-card/60 p-4 text-xs text-fg/70">
      <div v-if="!props.active" class="py-10 text-center text-[10px] uppercase tracking-[0.35em] text-fg/40">
        {{ t("layout.select", props.lang) }}
      </div>
      <div v-else-if="filteredLines.length === 0" class="py-10 text-center text-[10px] uppercase tracking-[0.35em] text-fg/40">
        {{ t("logs.empty", props.lang) }}
      </div>
      <div v-else class="space-y-1 font-mono">
        <div v-for="(line, index) in filteredLines" :key="index" data-testid="log-line" class="whitespace-pre-wrap">
          {{ line }}
        </div>
      </div>
    </div>
  </div>
</template>
