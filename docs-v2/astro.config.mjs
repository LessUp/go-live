import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightImageZoom from 'starlight-image-zoom';
import tailwind from '@tailwindcss/vite';
import sitemap from '@astrojs/sitemap';
import partytown from '@astrojs/partytown';

// https://astro.build/config
export default defineConfig({
  site: 'https://go-live.dev',
  base: '/',
  
  // 构建优化
  build: {
    format: 'directory',
    inlineStylesheets: 'auto',
  },
  
  // 性能优化
  prefetch: {
    prefetchAll: true,
    defaultStrategy: 'hover',
  },
  
  // 图片优化
  image: {
    service: {
      entrypoint: 'astro/assets/services/sharp',
    },
  },
  
  vite: {
    plugins: [tailwind()],
    build: {
      cssMinify: 'lightningcss',
    },
    ssr: {
      noExternal: ['@astrojs/starlight'],
    },
  },
  
  integrations: [
    // Starlight 文档主题
    starlight({
      title: 'Go-Live',
      tagline: 'Lightweight WebRTC SFU Server',
      
      // Logo 配置
      logo: {
        light: './src/assets/logo-light.svg',
        dark: './src/assets/logo-dark.svg',
        replacesTitle: true,
      },
      
      // 编辑链接
      editLink: {
        baseUrl: 'https://github.com/LessUp/go-live/edit/master/docs-v2/',
      },
      
      // 最后更新时间
      lastUpdated: true,
      
      // 分页
      pagination: true,
      
      // 站点地图
      sitemap: {
        xmlns: {
          xhtml: true,
        },
        xslUrl: '/sitemap.xsl',
      },
      
      // 自定义 CSS
      customCss: [
        './src/styles/custom.css',
        './src/styles/code.css',
      ],
      
      // 语言配置
      defaultLocale: 'root',
      locales: {
        root: {
          label: '简体中文',
          lang: 'zh-CN',
          dir: 'ltr',
        },
        en: {
          label: 'English',
          lang: 'en',
        },
      },
      
      // 导航侧边栏配置
      sidebar: [
        {
          label: 'Getting Started',
          translations: { 'zh-CN': '快速开始' },
          items: [
            { label: 'Introduction', link: '/', translations: { 'zh-CN': '介绍' } },
            { label: 'Quick Start', link: '/getting-started/quickstart/', translations: { 'zh-CN': '快速开始' } },
            { label: 'Installation', link: '/getting-started/installation/', translations: { 'zh-CN': '安装指南' } },
          ],
        },
        {
          label: 'Deployment',
          translations: { 'zh-CN': '部署' },
          items: [
            { label: 'Docker', link: '/deployment/docker/', translations: { 'zh-CN': 'Docker' } },
            { label: 'Kubernetes', link: '/deployment/kubernetes/', translations: { 'zh-CN': 'Kubernetes' } },
            { label: 'Binary', link: '/deployment/binary/', translations: { 'zh-CN': '二进制' } },
          ],
        },
        {
          label: 'Configuration',
          translations: { 'zh-CN': '配置' },
          items: [
            { label: 'Overview', link: '/configuration/overview/', translations: { 'zh-CN': '配置概览' } },
            { label: 'Authentication', link: '/configuration/authentication/', translations: { 'zh-CN': '认证配置' } },
            { label: 'WebRTC', link: '/configuration/webrtc/', translations: { 'zh-CN': 'WebRTC' } },
            { label: 'Recording', link: '/configuration/recording/', translations: { 'zh-CN': '录制配置' } },
            { label: 'Monitoring', link: '/configuration/monitoring/', translations: { 'zh-CN': '监控配置' } },
          ],
        },
        {
          label: 'API Reference',
          translations: { 'zh-CN': 'API 参考' },
          collapsed: false,
          items: [
            { label: 'Overview', link: '/api/', translations: { 'zh-CN': '概览' } },
            { label: 'WHIP Publish', link: '/api/whip/', translations: { 'zh-CN': 'WHIP 推流' } },
            { label: 'WHEP Play', link: '/api/whep/', translations: { 'zh-CN': 'WHEP 播放' } },
            { label: 'Room Management', link: '/api/rooms/', translations: { 'zh-CN': '房间管理' } },
            { label: 'Admin API', link: '/api/admin/', translations: { 'zh-CN': '管理 API' } },
            { label: 'Health & Metrics', link: '/api/metrics/', translations: { 'zh-CN': '健康与指标' } },
          ],
        },
        {
          label: 'Architecture',
          translations: { 'zh-CN': '架构' },
          items: [
            { label: 'Overview', link: '/architecture/', translations: { 'zh-CN': '架构概览' } },
            { label: 'SFU Manager', link: '/architecture/sfu/', translations: { 'zh-CN': 'SFU 管理器' } },
            { label: 'Room & Fanout', link: '/architecture/room/', translations: { 'zh-CN': '房间与转发' } },
            { label: 'Authentication', link: '/architecture/auth/', translations: { 'zh-CN': '认证系统' } },
            { label: 'Recording & Upload', link: '/architecture/recording/', translations: { 'zh-CN': '录制与上传' } },
          ],
        },
        {
          label: 'Troubleshooting',
          translations: { 'zh-CN': '故障排除' },
          items: [
            { label: 'Common Issues', link: '/troubleshooting/common/', translations: { 'zh-CN': '常见问题' } },
            { label: 'WebRTC Connection', link: '/troubleshooting/webrtc/', translations: { 'zh-CN': 'WebRTC 连接' } },
            { label: 'Performance', link: '/troubleshooting/performance/', translations: { 'zh-CN': '性能问题' } },
          ],
        },
      ],
      
      // 社交链接
      social: {
        github: 'https://github.com/LessUp/go-live',
      },
      
      // 页脚
      pagefind: true, // 启用搜索
      expressiveCode: true, // 代码块增强
      
      // 头部组件
      components: {
        // 自定义组件覆盖
        Header: './src/components/CustomHeader.astro',
        ThemeSelect: './src/components/ThemeSelect.astro',
      },
      
      // 插件
      plugins: [
        // 图片点击放大
        starlightImageZoom(),
      ],
      
      // 头部元数据
      head: [
        // PWA
        {
          tag: 'link',
          attrs: {
            rel: 'manifest',
            href: '/manifest.json',
          },
        },
        {
          tag: 'meta',
          attrs: {
            name: 'theme-color',
            content: '#00ADD8',
          },
        },
        {
          tag: 'meta',
          attrs: {
            name: 'apple-mobile-web-app-capable',
            content: 'yes',
          },
        },
        // 预连接优化
        {
          tag: 'link',
          attrs: {
            rel: 'preconnect',
            href: 'https://fonts.googleapis.com',
          },
        },
        {
          tag: 'link',
          attrs: {
            rel: 'preconnect',
            href: 'https://fonts.gstatic.com',
            crossorigin: true,
          },
        },
        // Open Graph 图片
        {
          tag: 'meta',
          attrs: {
            property: 'og:image',
            content: 'https://go-live.dev/og-image.png',
          },
        },
        {
          tag: 'meta',
          attrs: {
            property: 'twitter:image',
            content: 'https://go-live.dev/og-image.png',
          },
        },
      ],
    }),
    
    // 站点地图
    sitemap({
      filter: (page) => !page.includes('/api-playground/'),
      changefreq: 'weekly',
      priority: 0.7,
      lastmod: new Date(),
      i18n: {
        defaultLocale: 'zh-CN',
        locales: {
          'zh-CN': 'zh-CN',
          'en': 'en-US',
        },
      },
    }),
    
    // Partytown - 将第三方脚本移至 Web Worker
    partytown({
      config: {
        forward: ['dataLayer.push'],
      },
    }),
  ],
  
  // Markdown 配置
  markdown: {
    shikiConfig: {
      theme: 'github-dark',
      themes: {
        light: 'github-light',
        dark: 'github-dark',
      },
      wrap: true,
    },
  },
});
