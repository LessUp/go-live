import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

const base = '/go-live/'

export default withMermaid(defineConfig({
  base,
  title: 'Go-Live',
  description: 'Lightweight WebRTC SFU Server — WHIP/WHEP streaming, room broadcast, recording & observability',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/go-live/media/favicon.svg' }],
    ['meta', { name: 'theme-color', content: '#0969da' }],
  ],

  locales: {
    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      title: 'Go-Live',
      description: 'Lightweight WebRTC SFU Server — WHIP/WHEP streaming, room broadcast, recording & observability',
      themeConfig: {
        nav: [
          { text: 'Getting Started', link: '/en/getting-started', activeMatch: '/en/getting-started' },
          { text: 'Architecture', link: '/en/architecture/overview', activeMatch: '/en/architecture/' },
          { text: 'Protocols', link: '/en/protocols/whip', activeMatch: '/en/protocols/' },
          { text: 'Features', link: '/en/features/auth', activeMatch: '/en/features/' },
          { text: 'API', link: '/en/api/endpoints', activeMatch: '/en/api/' },
        ],
        sidebar: {
          '/en/architecture/': [
            {
              text: 'Architecture',
              items: [
                { text: 'System Overview', link: '/en/architecture/overview' },
                { text: 'SFU Core', link: '/en/architecture/sfu-core' },
                { text: 'Data Flow', link: '/en/architecture/data-flow' },
                { text: 'Deployment', link: '/en/architecture/deployment' },
              ],
            },
          ],
          '/en/protocols/': [
            {
              text: 'Protocols',
              items: [
                { text: 'WHIP Publishing', link: '/en/protocols/whip' },
                { text: 'WHEP Playback', link: '/en/protocols/whep' },
              ],
            },
          ],
          '/en/features/': [
            {
              text: 'Features',
              items: [
                { text: 'Authentication', link: '/en/features/auth' },
                { text: 'Recording', link: '/en/features/recording' },
                { text: 'Observability', link: '/en/features/observability' },
              ],
            },
          ],
          '/en/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Endpoints', link: '/en/api/endpoints' },
                { text: 'Configuration', link: '/en/api/configuration' },
              ],
            },
          ],
        },
      },
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/zh/',
      title: 'Go-Live',
      description: '轻量级 WebRTC SFU 服务器 — WHIP/WHEP 流媒体、房间广播、录制与可观测性',
      themeConfig: {
        nav: [
          { text: '快速开始', link: '/zh/getting-started', activeMatch: '/zh/getting-started' },
          { text: '系统架构', link: '/zh/architecture/overview', activeMatch: '/zh/architecture/' },
          { text: '协议规范', link: '/zh/protocols/whip', activeMatch: '/zh/protocols/' },
          { text: '功能模块', link: '/zh/features/auth', activeMatch: '/zh/features/' },
          { text: 'API 参考', link: '/zh/api/endpoints', activeMatch: '/zh/api/' },
        ],
        sidebar: {
          '/zh/architecture/': [
            {
              text: '系统架构',
              items: [
                { text: '系统总览', link: '/zh/architecture/overview' },
                { text: 'SFU 核心', link: '/zh/architecture/sfu-core' },
                { text: '数据流', link: '/zh/architecture/data-flow' },
                { text: '部署方案', link: '/zh/architecture/deployment' },
              ],
            },
          ],
          '/zh/protocols/': [
            {
              text: '协议规范',
              items: [
                { text: 'WHIP 发布', link: '/zh/protocols/whip' },
                { text: 'WHEP 播放', link: '/zh/protocols/whep' },
              ],
            },
          ],
          '/zh/features/': [
            {
              text: '功能模块',
              items: [
                { text: '认证系统', link: '/zh/features/auth' },
                { text: '录制功能', link: '/zh/features/recording' },
                { text: '可观测性', link: '/zh/features/observability' },
              ],
            },
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '端点列表', link: '/zh/api/endpoints' },
                { text: '配置项', link: '/zh/api/configuration' },
              ],
            },
          ],
        },
      },
    },
  },

  themeConfig: {
    outline: [2, 3],
    search: { provider: 'local' },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/LessUp/go-live' },
    ],
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present LessUp',
    },
  },

  markdown: {
    lineNumbers: true,
  },
}))
