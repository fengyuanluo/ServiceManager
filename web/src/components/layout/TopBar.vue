<script setup lang="ts">
import { t, type Lang } from "../../i18n";
import type { ThemeId } from "../../data/navigation";

const props = defineProps<{
  title: string;
  lang: Lang;
  themes: { id: ThemeId; labelKey: string; swatch: string }[];
  activeTheme: ThemeId;
}>();

const emit = defineEmits<{
  (event: "toggle-lang"): void;
  (event: "set-theme", id: ThemeId): void;
  (event: "toggle-drawer"): void;
}>();
</script>

<template>
  <header class="sticky top-0 z-30 border-b border-border/60 bg-bg/80 backdrop-blur">
    <div class="flex items-center gap-3 px-4 py-4 lg:px-8">
      <button
        data-testid="menu-toggle"
        type="button"
        class="rounded-full border border-border/70 p-2 text-fg/70 transition hover:text-fg lg:hidden"
        :aria-label="t('action.menu', props.lang)"
        @click="emit('toggle-drawer')"
      >
        <svg
          viewBox="0 0 24 24"
          class="h-4 w-4"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
          stroke-linecap="round"
        >
          <path d="M4 7h16" />
          <path d="M4 12h16" />
          <path d="M4 17h16" />
        </svg>
      </button>

      <div class="flex-1">
        <div class="font-display text-lg tracking-wide text-fg">
          {{ props.title }}
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div class="hidden items-center gap-2 sm:flex">
          <button
            v-for="theme in props.themes"
            :key="theme.id"
            :data-testid="`theme-${theme.id}`"
            type="button"
            class="h-4 w-4 rounded-full border border-border/70 transition hover:scale-110"
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
          data-testid="lang-toggle"
          type="button"
          class="rounded-full border border-border/70 px-2.5 py-1 text-[11px] uppercase tracking-[0.35em] text-fg/70 transition hover:text-fg"
          :aria-label="t('action.language', props.lang)"
          @click="emit('toggle-lang')"
        >
          {{ props.lang === "zh" ? "中" : "EN" }}
        </button>
      </div>
    </div>
  </header>
</template>
