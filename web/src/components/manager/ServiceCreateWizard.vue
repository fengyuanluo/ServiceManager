<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import type { CreateRequest } from "../../types/service";

const props = defineProps<{
  open: boolean;
  lang: Lang;
  managerKey: ManagerKey;
  submitting?: boolean;
  error?: string;
}>();

const emit = defineEmits<{
  (event: "close"): void;
  (event: "submit", payload: CreateRequest): void;
}>();

const steps = [
  { id: "template", labelKey: "create.step.template" },
  { id: "form", labelKey: "create.step.form" },
  { id: "confirm", labelKey: "create.step.confirm" }
];

const mode = ref<"form" | "raw">("form");
const rawText = ref("");
const form = ref({
  name: "",
  description: "",
  exec: "",
  workdir: "",
  user: "root",
  env: ""
});
const options = ref({ start: false, autostart: false });

const templateLabel = computed(() => t(`nav.${props.managerKey}`, props.lang));
const activeStep = computed(() => "form");
const submitting = computed(() => props.submitting === true);
const canSubmit = computed(() => form.value.name.trim().length > 0 && !submitting.value);

const resetForm = () => {
  form.value = {
    name: "",
    description: "",
    exec: "",
    workdir: "",
    user: "root",
    env: ""
  };
  options.value = { start: false, autostart: false };
  rawText.value = "";
  mode.value = "form";
};

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm();
    }
  }
);

const toggleMode = () => {
  mode.value = mode.value === "form" ? "raw" : "form";
};

const buildPayload = (): CreateRequest => {
  const name = form.value.name.trim();
  const payload: CreateRequest = {
    name,
    template_id: props.managerKey,
    options: {
      start: options.value.start,
      autostart: options.value.autostart
    }
  };
  const description = form.value.description.trim();
  const exec = form.value.exec.trim();
  const workdir = form.value.workdir.trim();
  const user = form.value.user.trim();
  const env = form.value.env.trim();
  if (description) payload.description = description;
  if (exec) payload.exec = exec;
  if (workdir) payload.workdir = workdir;
  if (user) payload.user = user;
  if (env) payload.env = env;
  if (mode.value === "raw") {
    const raw = rawText.value.trim();
    if (raw) {
      payload.raw_text = raw;
    }
  }
  return payload;
};

const submit = () => {
  if (!canSubmit.value) return;
  emit("submit", buildPayload());
};
</script>

<template>
  <div
    v-if="props.open"
    data-testid="create-wizard"
    class="fixed inset-0 z-50 flex items-start justify-center px-4 py-6 sm:py-10"
  >
    <button
      type="button"
      class="absolute inset-0 bg-black/40 backdrop-blur-sm"
      aria-hidden="true"
      @click="emit('close')"
    ></button>
    <div
      role="dialog"
      aria-modal="true"
      class="relative w-full max-w-3xl max-h-[calc(100vh-4rem)] overflow-y-auto rounded-3xl border border-border/60 bg-card/95 p-5 shadow-2xl shadow-black/20 sm:p-6"
    >
      <header class="flex flex-wrap items-center justify-between gap-2">
        <div class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("create.title", props.lang) }}
        </div>
        <div class="flex items-center gap-2">
          <button
            data-testid="create-toggle-raw"
            type="button"
            class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
            @click="toggleMode"
          >
            {{ mode === "form" ? t("action.raw", props.lang) : t("action.form", props.lang) }}
          </button>
          <button
            type="button"
            class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
            :aria-label="t('action.close', props.lang)"
            @click="emit('close')"
          >
            {{ t("action.close", props.lang) }}
          </button>
        </div>
      </header>

      <div class="mt-4 flex flex-wrap gap-2">
        <div
          v-for="step in steps"
          :key="step.id"
          class="rounded-full border px-3 py-1 text-[10px] uppercase tracking-[0.35em]"
          :class="
            step.id === activeStep
              ? 'border-primary/60 bg-primary/10 text-primary'
              : 'border-border/60 text-fg/40'
          "
        >
          {{ t(step.labelKey, props.lang) }}
        </div>
      </div>

      <section class="mt-6 grid gap-4 rounded-2xl border border-border/60 bg-card/70 p-4">
        <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("create.section.template", props.lang) }}
        </div>
        <div class="flex items-center gap-3">
          <div class="rounded-full border border-primary/60 bg-primary/10 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-primary">
            {{ templateLabel }}
          </div>
          <div class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("create.section.selected", props.lang) }}
          </div>
        </div>
      </section>

      <section class="mt-6 space-y-4" v-if="mode === 'form'">
        <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("create.section.form", props.lang) }}
        </div>
        <div class="grid gap-3 md:grid-cols-2">
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.name", props.lang) }}
            </span>
            <input
              v-model="form.name"
              data-testid="create-name"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
              :disabled="submitting"
            />
          </label>
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.user", props.lang) }}
            </span>
            <select
              v-model="form.user"
              data-testid="create-user"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
              :disabled="submitting"
            >
              <option value="root">root</option>
              <option value="app">app</option>
              <option value="node">node</option>
            </select>
          </label>
        </div>
        <label class="space-y-2 text-xs text-fg/60">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("editor.description", props.lang) }}
          </span>
          <input
            v-model="form.description"
            data-testid="create-description"
            type="text"
            class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            :disabled="submitting"
          />
        </label>
        <div class="grid gap-3 md:grid-cols-2">
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.exec", props.lang) }}
            </span>
            <input
              v-model="form.exec"
              data-testid="create-exec"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
              :disabled="submitting"
            />
          </label>
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.workdir", props.lang) }}
            </span>
            <input
              v-model="form.workdir"
              data-testid="create-workdir"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
              :disabled="submitting"
            />
          </label>
        </div>
        <label class="space-y-2 text-xs text-fg/60">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("editor.env", props.lang) }}
          </span>
          <textarea
            v-model="form.env"
            data-testid="create-env"
            rows="3"
            class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            :disabled="submitting"
          ></textarea>
        </label>
      </section>

      <section class="mt-6 space-y-4" v-else>
        <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("create.section.raw", props.lang) }}
        </div>
        <textarea
          v-model="rawText"
          data-testid="create-raw"
          rows="12"
          class="w-full rounded-2xl border border-border/60 bg-card/70 p-3 font-mono text-xs"
          :disabled="submitting"
        ></textarea>
      </section>

      <section class="mt-6 space-y-3 rounded-2xl border border-border/60 bg-card/60 p-4">
        <div class="text-[10px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("create.section.options", props.lang) }}
        </div>
        <div class="flex flex-wrap items-center gap-4 text-xs text-fg/60">
          <label class="flex items-center gap-2">
            <input
              v-model="options.start"
              data-testid="create-start"
              type="checkbox"
              class="h-3 w-3 rounded border-border/60 bg-card/70 text-primary"
              :disabled="submitting"
            />
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("action.start", props.lang) }}
            </span>
          </label>
          <label class="flex items-center gap-2">
            <input
              v-model="options.autostart"
              data-testid="create-autostart"
              type="checkbox"
              class="h-3 w-3 rounded border-border/60 bg-card/70 text-primary"
              :disabled="submitting"
            />
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("action.autostart", props.lang) }}
            </span>
          </label>
        </div>
      </section>

      <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
        <div v-if="props.error" class="text-xs text-rose-500">
          {{ props.error }}
        </div>
        <div class="ml-auto flex items-center gap-2">
          <button
            data-testid="create-cancel"
            type="button"
            class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
            @click="emit('close')"
          >
            {{ t("action.close", props.lang) }}
          </button>
          <button
            data-testid="create-submit"
            type="button"
            class="rounded-full border border-primary/60 bg-primary/10 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-primary transition hover:bg-primary/20 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!canSubmit"
            @click="submit"
          >
            {{ submitting ? t("create.submitting", props.lang) : t("action.create", props.lang) }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
