import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { App } from 'vue'
import { Layout, getParentLayout } from '@/utils/routerHelper'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n()
export const constantRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/',
    component: Layout,
    redirect: '/login',
    name: 'Root',
    meta: {
      hidden: true
    }
  },
  {
    path: '/redirect',
    component: Layout,
    name: 'Redirect',
    children: [
      {
        path: '/redirect/:path(.*)',
        name: 'Redirect',
        component: () => import('@/views/Login/Login.vue'),
        meta: {}
      }
    ],
    meta: {
      hidden: true,
      noTagsView: true
    }
  },
  {
    path: '/login',
    component: () => import('@/views/Login/Login.vue'),
    name: 'Login',
    meta: {
      hidden: true,
      title: t('router.login'),
      noTagsView: true
    }
  },
  {
    path: '/test',
    component: () => import('@/views/test/test.vue'),
    name: 'test',
    meta: {
      hidden: true,
      title: t('router.test'),
      noTagsView: true
    }
  },
  {
    path: '/404',
    component: () => import('@/views/Error/404.vue'),
    name: 'NoFind',
    meta: {
      hidden: true,
      title: '404',
      noTagsView: true
    }
  }
]

export const asyncRouterMap: AppRouteRecordRaw[] = [

  {
    path: '/sql_dashboard',
    component: Layout,
    // redirect: '/dashboard1/useWatermark',
    name: 'sql_dashboard',
    meta: {
      projectType:'bug_platform',
      title: 'sql注入',
      icon: '',
      alwaysShow: true
    },
    children: [
      {
        path: 'sql注入-sqlite3',
        component: () => import('@/views/sql_dashboard/index.vue'),
        name: 'sql注入-sqlite3',
        meta: {
          title: 'sql注入-sqlite3',
          desc: "SQL注入漏洞的原理是基于应用程序将用户输入的数据直接拼接到SQL查询语句中，而没有进行适当的验证和转义。攻击者可以通过构造特定的输入数据，使得拼接后的SQL语句执行非预期的操作。它发生在应用程序未能正确地过滤或转义用户输入的数据，导致攻击者能够将恶意的SQL代码注入到应用程序的数据库查询中。这种漏洞允许攻击者绕过应用程序的身份验证和授权机制，访问、修改或删除数据库中的敏感数据，甚至控制整个数据库服务器。"
        }
      }
      // {
      //   path: 'sql注入-mysql',
      //   component: () => import('@/views/sql_dashboard/index.vue'),
      //   name: 'sql注入-mysql',
      //   meta: {
      //     title: 'sql注入-mysql',
      //     desc: "sql注入-mysql"
      //   }
      // },
    ]
  },
  {
    path: '/xss_dashboard',
    component: Layout,
    name: 'xss_dashboard',
    meta: {
      projectType:'bug_platform',
      title: 'XSS',
      icon: '',
      alwaysShow: true
    },
    children: [
      {
        path: '反射型xss',
        component: () => import('@/views/xss_dashboard/reflect_xss/index.vue'),
        name: '反射型xss',
        meta: {
          title: '反射型xss',
          desc: "反射型跨站脚本攻击（Reflected Cross-Site Scripting，简称反射型XSS）是一种常见的网络安全漏洞，其原理是攻击者通过诱使用户点击一个包含恶意脚本的URL，使得恶意脚本在用户的浏览器中执行。这种攻击方式的特点是恶意脚本不会存储在服务器上，而是通过URL传递并在用户的浏览器中即时执行。"
        }
      },
      {
        path: '存储型xss',
        component: () => import('@/views/xss_dashboard/store_xss/index.vue'),
        name: '存储型xss',
        meta: {
          title: '存储型xss',
          desc: "存储型跨站脚本攻击（Stored Cross-Site Scripting，简称存储型XSS）是另一种常见的网络安全漏洞，其原理是攻击者将恶意脚本永久地存储在目标服务器上，使得所有访问受影响页面的用户都会在他们的浏览器中执行这些恶意脚本。与反射型XSS不同，存储型XSS的恶意脚本不会通过URL传递，而是存储在服务器的数据库或文件中。"
        }
      },
    ]
  },
  {
    path: '/file_dashboard',
    component: Layout,
    name: 'file_dashboard',
    meta: {
      projectType:'bug_platform',
      title: '任意文件操作',
      icon: '',
      alwaysShow: true
    },
    children: [
      {
        path: '任意文件上传',
        component: () => import('@/views/file_dashboard/upload/index.vue'),
        name: '任意文件上传',
        meta: {
          title: '任意文件上传',
          desc: "任意文件上传（Arbitrary File Upload）是一种常见的安全漏洞，攻击者可以利用此漏洞将恶意文件上传到目标服务器。这种漏洞可能导致服务器被入侵、数据泄露、网站被篡改等严重后果。一般发生在存在文件上传的功能的同时，没有对文件类型和路径进行校验的场景。"
        }
      },
      {
        path: '任意文件下载',
        component: () => import('@/views/file_dashboard/download/index.vue'),
        name: '任意文件下载',
        meta: {
          title: '任意文件下载',
          desc: "任意文件下载（Arbitrary File Download），也称为路径遍历（Path Traversal），是一种安全漏洞，攻击者可以利用此漏洞绕过应用程序的访问控制，下载服务器上的任意文件。这种漏洞可能导致敏感信息泄露，如配置文件、源代码、用户数据等。"
        }
      },
    ]
  },
  
  {
    path: '/rce_dashboard',
    component: Layout,
    name: 'rce_dashboard',
    redirect: '/rce_dashboard/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: 'RCE远程命令执行',
        component: () => import('@/views/rce_dashboard/index.vue'),
        name: 'RCE远程命令执行',
        meta: {
          title: 'RCE远程命令执行',
          desc: "远程命令执行（Remote Command Execution，简称RCE）是一种严重的安全漏洞，攻击者可以利用此漏洞在目标服务器上执行任意命令，从而完全控制服务器。当应用程序没有正确处理用户输入时，攻击者可以通过注入恶意代码来执行系统命令。"
        }
      }
    ]
  },
  {
    path: '/xxe_dashboard',
    component: Layout,
    name: 'xxe_dashboard',
    redirect: '/xxe_dashboard/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: 'XXE',
        component: () => import('@/views/xxe_dashboard/index.vue'),
        name: 'XXE',
        meta: {
          title: 'XXE',
          desc: "XXE（Unauthorized Access to API）是一种常见的安全漏洞，攻击者可以利用此漏洞访问或操作未经授权的API接口，从而获取敏感信息或执行恶意操作。这种漏洞通常发生在Web应用程序中，当API接口没有正确实施访问控制时，攻击者可以通过直接调用API接口来绕过前端的验证和授权机制。"
          // icon: 'clarity:document-solid'
        }
      }
    ]
  },
  {
    path: '/personinfo_dashboard',
    component: Layout,
    name: 'personinfo_dashboard',
    redirect: '/personinfo_dashboard/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: '越权漏洞',
        component: () => import('@/views/personinfo_dashboard/index.vue'),
        name: '越权漏洞',
        meta: {
          title: '越权漏洞',
          desc: "越权漏洞（Authorization Bypass）是一种常见的安全漏洞，攻击者可以利用此漏洞绕过系统的访问控制机制，访问或操作不属于的资源或者功能。例如，一个普通用户能够访问另一个普通用户的私人信息或数据。"
          // icon: 'clarity:document-solid'
        }
      }
    ]
  },
  {
    path: '/personinfo_dashboard_unauthorized',
    component: Layout,
    name: 'personinfo_dashboard_unauthorized',
    redirect: '/personinfo_dashboard_unauthorized/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: '接口未授权访问',
        component: () => import('@/views/personinfo_dashboard_unauthorized/index.vue'),
        name: '接口未授权访问',
        meta: {
          title: '接口未授权访问',
          desc: "接口未授权访问（Unauthorized Access to API）是一种常见的安全漏洞，攻击者可以利用此漏洞访问或操作未经授权的API接口，从而获取敏感信息或执行恶意操作。这种漏洞通常发生在Web应用程序中，当API接口没有正确实施访问控制时，攻击者可以通过直接调用API接口来绕过前端的验证和授权机制。"
          // icon: 'clarity:document-solid'
        }
      }
    ]
  },
  {
    path: '/buy_dashboard',
    component: Layout,
    name: 'buy_dashboard',
    redirect: '/buy_dashboard/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: '负值反冲',
        component: () => import('@/views/buy_dashboard/index.vue'),
        name: '负值反冲',
        meta: {
          title: '负值反冲',
          desc: "负值反冲漏洞（Negative Rebound Vulnerability）是一个在某些计算机系统或网络协议中可能出现的安全漏洞。这个漏洞的原理通常涉及到在处理特定类型的数据或执行某些操作时，系统未能正确处理负值或异常值，导致系统状态出现意外的反冲或回弹。具体来说，当系统接收到一个负值输入时，它可能会错误地处理这个值，导致系统状态或资源分配出现负值，进而引发一系列不正常的操作或状态变化。这种不正常的操作可能会被恶意用户利用，通过发送特定的负值数据包来触发系统的异常行为，从而实施攻击，比如电商系统购买负数商品使账户余额不减反增。为了防止这种漏洞，系统开发者需要确保在设计和实现系统时，对所有可能的输入值进行严格的验证和处理，特别是对于负值和异常值的处理要格外小心，确保系统在任何情况下都能保持稳定和安全。"
          // icon: 'clarity:document-solid'
        }
      }
    ]
  },
  {
    path: '/edit_password_dashboard',
    component: Layout,
    name: 'edit_password_dashboard',
    redirect: '/edit_password_dashboard/index',
    meta: {
      projectType:'bug_platform',
    },
    children: [
      {
        path: '敏感信息明文传输',
        component: () => import('@/views/edit_password_dashboard/index.vue'),
        name: '敏感信息明文传输',
        meta: {
          title: '敏感信息明文传输',
          desc: "敏感信息明文传输是指在数据传输过程中，敏感信息（如密码、信用卡号、个人身份信息等）没有经过加密处理，直接以原始的、可读的形式在网络上传输。这种传输方式存在极大的安全风险，因为任何能够截获数据包的攻击者都可以轻易地读取和窃取这些敏感信息。"
          // icon: 'clarity:document-solid'
        }
      }
    ]
  },


]
export const asyncRouterMapSast: AppRouteRecordRaw[] = [
  {
    path: '/scan',
    component: Layout,
    name: 'scan',
    // redirect: '/scan1',
    meta: {
      projectType:'sast',
      title: 'scan',
      icon: '',
      alwaysShow: true
    },
    children: [
      {
        path: 'createscan',
        component: () => import('@/sastscanviews/createscan/index.vue'),
        name: 'createscan',
        meta: {
          title: '创建扫描任务',
          desc: "扫码器1"
        }
      },
      {
        path: 'takslist',
        component: () => import('@/sastscanviews/takslist/index.vue'),
        name: 'takslist',
        meta: {
          title: '查看扫描结果',
          desc: "查看扫描结果"
        }
      },
      {
        path: '/result/:taskID',
        component: () => import('@/sastscanviews/TaskResult/index.vue'),
        name: 'TaskResult',
        meta: {
          title: '任务结果',
          desc: '查看具体任务的扫描结果',
          hidden: true // 当设置 true 的时候该路由不会再侧边栏出现
        },
      },
    ]
  }
   
]

export const asyncRouterMapAiScan: AppRouteRecordRaw[] = [
  {
    path: '/scan',
    component: Layout,
    name: 'scan',
    // redirect: '/scan1',
    meta: {
      projectType:'aiscan',
      title: 'scan',
      icon: '',
      alwaysShow: true
    },
    children: [
      {
        path: 'createscan',
        component: () => import('@/aiscanviews/createscan/index.vue'),
        name: 'createscan',
        meta: {
          title: '创建扫描任务',
          desc: "扫码器1"
        }
      },
      {
        path: 'takslist',
        component: () => import('@/aiscanviews/takslist/index.vue'),
        name: 'takslist',
        meta: {
          title: '查看扫描结果',
          desc: "查看扫描结果"
        }
      },
      {
        path: '/result/:taskID',
        component: () => import('@/aiscanviews/TaskResult/index.vue'),
        name: 'TaskResult',
        meta: {
          title: '任务结果',
          desc: '查看具体任务的扫描结果',
          hidden: true // 当设置 true 的时候该路由不会再侧边栏出现
        },
      },
    ]
  }
   
]
const router = createRouter({
  history: createWebHashHistory(),
  strict: true,
  routes: constantRouterMap as RouteRecordRaw[],
  scrollBehavior: () => ({ left: 0, top: 0 })
})

export const resetRouter = (): void => {
  const resetWhiteNameList = ['Redirect', 'Login', 'NoFind', 'Root']
  router.getRoutes().forEach((route) => {
    const { name } = route
    if (name && !resetWhiteNameList.includes(name as string)) {
      router.hasRoute(name) && router.removeRoute(name)
    }
  })
}

export const setupRouter = (app: App<Element>) => {
  app.use(router)
}

export default router
