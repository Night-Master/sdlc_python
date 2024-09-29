
<script lang="tsx">
import { defineComponent, computed } from 'vue'
import { Collapse } from '@/components/Collapse'
import { LocaleDropdown } from '@/components/LocaleDropdown'
import { SizeDropdown } from '@/components/SizeDropdown'
import { UserInfo } from '@/components/UserInfo'
import { Screenfull } from '@/components/Screenfull'
import { Breadcrumb } from '@/components/Breadcrumb'
import { BreadMessage } from './BreadMessage'
import { useAppStore } from '@/store/modules/app'
import { useDesign } from '@/hooks/web/useDesign'
import {Nav} from './Nav'
import {Title} from './Title'

const { getPrefixCls, variables } = useDesign()
const prefixCls = getPrefixCls('tool-header')
const appStore = useAppStore()
// 面包屑
const breadcrumb = computed(() => appStore.getBreadcrumb)

// 折叠图标
const hamburger = computed(() => appStore.getHamburger)

// 全屏图标
const screenfull = computed(() => appStore.getScreenfull)

// 尺寸图标
const size = computed(() => appStore.getSize)

// 布局
const layout = computed(() => appStore.getLayout)

// 多语言图标
// const locale = computed(() => appStore.getLocale)

export default defineComponent({
  name: 'ToolHeader',
  setup() {
    return () => (
      <div
        style='border-bottom: 2px solid #e9ecef!important;height:80px'
        id={`${variables.namespace}-tool-header`}
        class={[
          prefixCls,
          'h-[var(--top-tool-height)] relative px-[var(--top-tool-p-x)] flex items-center justify-between',
          'dark:bg-[var(--el-bg-color)]'
        ]}
      >
        {layout.value !== 'top' ? (
          <div class="h-full flex items-center">
            <div class="flex items-center">
             <Title />
            </div>
            {/* {hamburger.value && layout.value !== 'cutMenu' ? (
              <Collapse class="hover-trigger" color="var(--top-header-text-color)"></Collapse>
            ) : undefined}
            {breadcrumb.value ? <Breadcrumb class="<md:hidden"></Breadcrumb> : undefined} */}
          </div>
        ) : undefined}
        <div class="h-full flex items-center">
          {/* {true ? <BreadMessage /> : undefined} */}
          {screenfull.value ? (
            <Screenfull class="hover-trigger" color="var(--top-header-text-color)"></Screenfull>
          ) : undefined}
          {size.value ? (
            <SizeDropdown class="hover-trigger" color="var(--top-header-text-color)"></SizeDropdown>
          ) : undefined}
          {/* {locale.value ? (
            <LocaleDropdown
              class="hover-trigger"
              color="var(--top-header-text-color)"
            ></LocaleDropdown>
          ) : undefined} */}
          {/* <UserInfo class="hover-trigger"></UserInfo> */}
         <Nav></Nav>
        </div>
      </div>
    )
  }
})
</script>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-tool-header';

.@{prefix-cls} {
  transition: left var(--transition-time-02);
}
</style>
