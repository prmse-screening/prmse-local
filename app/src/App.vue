<template>
  <div>
    <toaster
      :theme="theme.global.name.value as 'light' | 'dark'"
      position="top-right"
      expand
      close-button
      rich-colors
    />
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useTheme } from 'vuetify'
import { Toaster } from 'vue-sonner'
// import { useColorMode } from '@vueuse/core'
// import { watch } from 'vue'

// const mode = useColorMode()
const theme = useTheme()

// watch(mode, () => (theme.global.name.value = mode.value))
</script>
