<script setup lang="ts">
import type { ManagerKey } from "../../data/navigation";

const props = defineProps<{ name: ManagerKey }>();
</script>

<template>
  <svg
    viewBox="0 0 24 24"
    class="h-4 w-4"
    fill="none"
    stroke="currentColor"
    stroke-width="1.6"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    <g v-if="props.name === 'systemd'">
      <path d="M4 6h16" />
      <path d="M4 12h16" />
      <path d="M4 18h16" />
    </g>
    <g v-else-if="props.name === 'sysv'">
      <path d="M4 7h10" />
      <path d="M4 12h16" />
      <path d="M4 17h10" />
    </g>
    <g v-else-if="props.name === 'openrc'">
      <path d="M12 4v16" />
      <path d="M4 12h16" />
    </g>
    <g v-else-if="props.name === 'supervisor'">
      <rect x="5" y="5" width="14" height="14" rx="2" />
      <path d="M8 9h8" />
      <path d="M8 13h8" />
    </g>
    <g v-else>
      <path d="M5 12h14" />
      <circle cx="8" cy="12" r="2" />
      <circle cx="16" cy="12" r="2" />
    </g>
  </svg>
</template>
