<script setup lang="ts">
import { ref } from "vue";
import { t, type Lang } from "../../i18n";
import type { ThemeId } from "../../data/navigation";

const props = defineProps<{
  lang: Lang;
  themes: { id: ThemeId; labelKey: string; swatch: string }[];
  activeTheme: ThemeId;
  loading?: boolean;
  error?: string;
}>();

const emit = defineEmits<{
  (event: "toggle-lang"): void;
  (event: "set-theme", id: ThemeId): void;
  (event: "submit", payload: { username: string; password: string; expire: number }): void;
}>();

const username = ref("");
const password = ref("");
</script>

<template>
  <div class="flex min-h-screen items-center justify-center px-6 py-10">
    <div
      data-testid="login-panel"
      class="relative w-full max-w-md overflow-hidden rounded-3xl border border-border/60 bg-card/80 p-8 shadow-[0_30px_60px_-40px_rgba(0,0,0,0.45)] backdrop-blur"
    >
      <div
        class="absolute inset-x-8 top-0 h-px bg-gradient-to-r from-transparent via-border/60 to-transparent"
      ></div>
      <header class="flex items-center justify-between">
        <div class="text-[10px] uppercase tracking-[0.45em] text-fg/50">SM</div>
        <div class="flex items-center gap-2">
          <div class="hidden items-center gap-2 sm:flex">
            <button
              v-for="theme in props.themes"
              :key="theme.id"
              type="button"
              class="h-3.5 w-3.5 rounded-full border border-border/70 transition hover:scale-110"
              :class="
                props.activeTheme === theme.id
                  ? 'ring-2 ring-primary ring-offset-2 ring-offset-bg'
                  : ''
              "
              :style="{ background: theme.swatch }"
              :aria-label="t(theme.labelKey, props.lang)"
              @click="emit('set-theme', theme.id)"
            ></button>
          </div>
          <button
            type="button"
            class="rounded-full border border-border/70 px-2.5 py-1 text-[11px] uppercase tracking-[0.35em] text-fg/70 transition hover:text-fg"
            :aria-label="t('action.language', props.lang)"
            @click="emit('toggle-lang')"
          >
            {{ props.lang === "zh" ? "中" : "EN" }}
          </button>
        </div>
      </header>

      <div class="mt-6 space-y-1">
        <div class="font-display text-2xl tracking-wide text-fg">
          {{ t("auth.title", props.lang) }}
        </div>
      </div>

      <form
        class="mt-6 space-y-4"
        @submit.prevent="
          emit('submit', {
            username: username.trim(),
            password,
            expire: 0
          })
        "
      >
        <div class="space-y-2">
          <label class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
            {{ t("auth.username", props.lang) }}
          </label>
          <input
            v-model="username"
            type="text"
            autocomplete="username"
            class="w-full rounded-2xl border border-border/70 bg-bg/80 px-4 py-3 text-sm text-fg shadow-sm outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/20"
          />
        </div>
        <div class="space-y-2">
          <label class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
            {{ t("auth.password", props.lang) }}
          </label>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            class="w-full rounded-2xl border border-border/70 bg-bg/80 px-4 py-3 text-sm text-fg shadow-sm outline-none transition focus:border-primary/60 focus:ring-2 focus:ring-primary/20"
          />
        </div>

        <div
          v-if="props.error"
          class="rounded-2xl border border-border/70 bg-muted/70 px-4 py-2 text-xs text-fg/70"
        >
          {{ props.error }}
        </div>

        <button
          type="submit"
          class="w-full rounded-2xl bg-primary px-4 py-3 text-sm tracking-[0.3em] text-primary-fg transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="props.loading || !username.trim() || !password"
        >
          {{ props.loading ? "..." : t("auth.submit", props.lang) }}
        </button>
      </form>
    </div>
  </div>
</template>
