import starlightPlugin from '@astrojs/starlight-tailwind';

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      // 自定义品牌颜色
      colors: {
        // Go 品牌色
        go: {
          light: '#00ADD8',
          DEFAULT: '#00ADD8',
          dark: '#00758F',
          50: '#e6f8fc',
          100: '#ccf1f9',
          200: '#99e3f3',
          300: '#66d5ed',
          400: '#33c7e7',
          500: '#00ADD8',
          600: '#008ab0',
          700: '#006888',
          800: '#004560',
          900: '#002338',
        },
        // 直播强调色
        live: {
          DEFAULT: '#FF6B6B',
          light: '#FF8585',
          dark: '#CC5555',
        },
        // Starlight 语义颜色
        sl: {
          // 明亮模式
          'color-white': '#ffffff',
          'color-gray-1': '#f6f8fa',
          'color-gray-2': '#e1e4e8',
          'color-gray-3': '#d1d5da',
          'color-gray-4': '#959da5',
          'color-gray-5': '#6a737d',
          'color-gray-6': '#24292f',
          'color-black': '#1b1f23',
          // 暗黑模式
          'color-dark-bg': '#0d1117',
          'color-dark-surface': '#161b22',
          'color-dark-border': '#30363d',
        },
      },
      
      // 字体
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'Consolas', 'Monaco', 'monospace'],
      },
      
      // 动画
      animation: {
        'fade-in': 'fadeIn 0.5s ease-out',
        'slide-up': 'slideUp 0.5s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'bounce-slow': 'bounce 2s infinite',
      },
      
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      
      // 过渡
      transitionTimingFunction: {
        'spring': 'cubic-bezier(0.4, 0, 0.2, 1)',
      },
      
      // 间距
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
        '128': '32rem',
      },
      
      // 阴影
      boxShadow: {
        'glow': '0 0 20px rgba(0, 173, 216, 0.3)',
        'glow-dark': '0 0 20px rgba(0, 173, 216, 0.5)',
        'card': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
      },
      
      // 圆角
      borderRadius: {
        'xl': '1rem',
        '2xl': '1.5rem',
        '3xl': '2rem',
      },
    },
  },
  plugins: [
    starlightPlugin(),
    // 添加自定义工具类
    function({ addUtilities }) {
      addUtilities({
        '.text-gradient': {
          'background-clip': 'text',
          '-webkit-background-clip': 'text',
          '-webkit-text-fill-color': 'transparent',
          'background-image': 'linear-gradient(135deg, #00ADD8, #00758F)',
        },
        '.text-gradient-live': {
          'background-clip': 'text',
          '-webkit-background-clip': 'text',
          '-webkit-text-fill-color': 'transparent',
          'background-image': 'linear-gradient(135deg, #FF6B6B, #FF8585)',
        },
        '.glass': {
          'background': 'rgba(255, 255, 255, 0.1)',
          'backdrop-filter': 'blur(10px)',
          'border': '1px solid rgba(255, 255, 255, 0.2)',
        },
        '.glass-dark': {
          'background': 'rgba(13, 17, 23, 0.8)',
          'backdrop-filter': 'blur(10px)',
          'border': '1px solid rgba(48, 54, 61, 0.5)',
        },
      });
    },
  ],
};
