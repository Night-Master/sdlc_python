<script setup lang="ts">
import { useTagsViewStore } from '@/store/modules/tagsView'
import { useAppStore } from '@/store/modules/app'
import { Footer } from '@/components/Footer'
import { computed } from 'vue'

const appStore = useAppStore()

const layout = computed(() => appStore.getLayout)

const fixedHeader = computed(() => appStore.getFixedHeader)

const footer = computed(() => appStore.getFooter)

const tagsViewStore = useTagsViewStore()

const getCaches = computed((): string[] => {
  return tagsViewStore.getCachedViews
})
</script>

<template>
  <section
    style="margin-left: 10px;width: calc(100% - 10px);padding-left: 10px!important;"
    :class="[
      'p-[var(--app-content-padding)]  bg-[var(--app-content-bg-color)] dark:bg-[var(--el-bg-color)]',
      {
        '!min-h-[calc(100%-var(--app-footer-height))]':
          ((fixedHeader && (layout === 'classic' || layout === 'topLeft')) || layout === 'top') &&
          footer,

        '!min-h-[calc(100%-var(--tags-view-height)-var(--top-tool-height)-var(--app-footer-height))]':
          !fixedHeader && layout === 'classic' && footer,

        '!min-h-[calc(100%-var(--tags-view-height)-var(--app-footer-height))]':
          !fixedHeader && (layout === 'topLeft' || layout === 'top') && footer,

        '!min-h-[calc(100%-var(--top-tool-height))]': fixedHeader && layout === 'cutMenu' && footer,

        '!min-h-[calc(100%-var(--top-tool-height)-var(--tags-view-height))]':
          !fixedHeader && layout === 'cutMenu' && footer
      }
    ]"
  >
    <router-view>
      <template #default="{ Component, route }">
        <keep-alive :include="getCaches">
          <component :is="Component" :key="route.fullPath" />
        </keep-alive>
      </template>
    </router-view>
  </section>
  <Footer v-if="footer" />
</template>
