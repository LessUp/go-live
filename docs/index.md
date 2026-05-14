---
layout: false
---

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  const lang = navigator.language.toLowerCase()
  const target = lang.startsWith('zh') ? '/go-live/zh/' : '/go-live/en/'
  window.location.replace(target)
})
</script>

<template>
  <div class="redirect">
    <p>Redirecting...</p>
    <p>
      <a href="/go-live/en/">English</a> | <a href="/go-live/zh/">中文</a>
    </p>
  </div>
</template>

<style>
.redirect {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
  text-align: center;
}
.redirect a {
  margin: 0 0.5rem;
}
</style>
