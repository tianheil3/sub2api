<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div v-else class="cyber-home">
    <div class="cyber-page">
          <header class="nav">
            <a class="brand" href="/home" aria-label="CyberRelay AI">
              <div class="mark" aria-hidden="true">
                <span></span><span></span><span></span><span></span><span></span>
                <span></span><span></span><span></span><span></span><span></span>
                <span></span><span></span><span></span><span></span><span></span>
              </div>
              <div class="brand-copy">
                <strong>CyberRelay AI</strong>
                <small>OpenAI-compatible relay</small>
              </div>
            </a>
            <nav class="nav-links" aria-label="主导航">
              <a href="#models">Models</a>
              <a href="#flow">Setup</a>
              <a href="#notes">Notes</a>
              <a class="home-button" href="/dashboard">进入中转站 ↘</a>
            </nav>
          </header>
    
          <main>
            <section class="hero" aria-label="首页首屏">
              <div class="hero-grid">
                <div>
                  <p class="eyebrow">AI API relay for teams and builders</p>
                  <h1>
                    <span>统一接入</span>
                    <span>GPT-5.6</span>
                    <span>GPT Image 2</span>
                  </h1>
                  <p class="lead">
                    <span class="lead-main">注册、兑换、创建秘钥，即可开始调用。</span>
                    <span class="lead-mobile">注册、兑换、创建秘钥。</span>
                  </p>
                  <ul class="lead-list">
                    <li>官方直转满血 GPT-5.6，支持 GPT Image 2 生图能力。</li>
                    <li>OpenAI 兼容 API，可接入 Codex、Claude Code、opencode、Droid、Zed、龙虾 OpenClaw 等开发工具。</li>
                  </ul>
                  <div class="hero-actions">
                    <a class="home-button" href="/register">邮箱注册</a>
                    <a class="home-button secondary" href="https://pay.ldxp.cn/shop/WNAXV9CI">购买卡密</a>
                    <a class="home-button secondary" href="#flow">查看购买使用流程</a>
                  </div>
                  <div class="model-strip" aria-label="模型能力">
                    <span class="pill">GPT-5.6 full</span>
                    <span class="pill">GPT Image 2</span>
                    <span class="pill">Official relay</span>
                    <span class="pill">OpenAI-compatible API</span>
                    <span class="pill">Codex / Claude Code / OpenClaw</span>
                  </div>
                </div>
    
                <div class="product" aria-label="产品界面预览">
                  <div class="product-top">
                    <div class="dots"><span></span><span></span><span></span></div>
                    <span>cyberrelay / api-keys / gpt-5.6</span>
                  </div>
                  <div class="product-body">
                    <aside class="gateway-sidebar">
                      <p class="panel-title">Gateway</p>
                      <a class="menu-item active" href="/keys">API 秘钥 ↗</a>
                      <a class="menu-item" href="/redeem">卡密兑换 ↗</a>
                      <a class="menu-item" href="/available-channels">模型列表 ↗</a>
                      <a class="menu-item" href="/usage">用量明细 ↗</a>
                    </aside>
                    <div class="workspace">
                      <div class="command">
                        <span class="accent">curl</span> https://ai.cyberrelay.net/v1/chat/completions \<br />
                        &nbsp;&nbsp;-H "Authorization: Bearer sk-..." \<br />
                        &nbsp;&nbsp;-d '{ "model": "gpt-5.6", "messages": [...] }'
                      </div>
                      <div class="log-list">
                        <div class="log-row"><span class="ok">200</span><span>gpt-5.6 · chat completion</span><span>1.28s</span></div>
                        <div class="log-row"><span class="ok">200</span><span>gpt-image-2 · image generation</span><span>4.72s</span></div>
                        <div class="log-row"><span class="ok">200</span><span>sticky session · selected group</span><span>0.86s</span></div>
                        <div class="log-row"><span class="ok">200</span><span>quota check · 100刀 / 7天</span><span>0.19s</span></div>
                      </div>
                    </div>
                    <aside class="inspector">
                      <p class="panel-title">Status</p>
                      <div class="metric"><strong>100刀</strong><span>建议创建额度</span></div>
                      <div class="metric"><strong>7天</strong><span>秘钥有效期</span></div>
                      <div class="metric"><strong>分组</strong><span>创建秘钥时选择</span></div>
                    </aside>
                  </div>
                </div>
              </div>
            </section>
    
            <section id="models">
              <div class="section-head">
                <h2>模型与接口能力</h2>
                <p class="section-copy">
                  面向开发者和团队的 OpenAI 兼容中转能力，重点覆盖文本、推理、图像生成和企业级用量场景。
                </p>
              </div>
              <div class="cards">
                <article class="card">
                  <span class="num">01 / model</span>
                  <h3>GPT-5.6 满血模型</h3>
                  <p>支持最新 GPT-5.6 满血模型，官方直转，适合高质量推理、编码、内容生成和复杂任务。</p>
                </article>
                <article class="card">
                  <span class="num">02 / image</span>
                  <h3>GPT Image 2 生图</h3>
                  <p>支持 GPT Image 2 图像生成能力，可通过 API 对接到自己的产品、工作流或自动化系统。</p>
                </article>
                <article class="card">
                  <span class="num">03 / api</span>
                  <h3>OpenAI 兼容 API</h3>
                  <p>使用熟悉的 OpenAI API 调用方式接入，降低迁移成本，适合已有应用和开发工具快速切换配置。</p>
                </article>
                <article class="card">
                  <span class="num">04 / tools</span>
                  <h3>开发工具接入</h3>
                  <p>API 可接入 Codex、Claude Code、opencode、Droid、Zed、龙虾 OpenClaw 等工具，适合编码、代理任务和自动化工作流。</p>
                </article>
              </div>
            </section>
    
            <section id="flow">
              <div class="section-head">
                <h2>购买使用流程</h2>
                <p class="section-copy">按顺序完成注册、兑换、创建秘钥和配置，即可开始调用接口。</p>
              </div>
              <div class="steps">
                <article class="step">
                  <span class="index">01</span>
                  <h3>邮箱注册</h3>
                  <p>点击邮箱注册，完成邮箱注册。</p>
                </article>
                <article class="step">
                  <span class="index">02</span>
                  <h3>卡密兑换</h3>
                  <p>购买卡密后，在中转站内进入兑换入口完成兑换。</p>
                </article>
                <article class="step">
                  <span class="index">03</span>
                  <h3>创建秘钥</h3>
                  <p>点击 API秘钥，创建秘钥，选择 100刀 7天有效期，并选择对应分组。</p>
                </article>
                <article class="step">
                  <span class="index">04</span>
                  <h3>配置使用</h3>
                  <p>点击使用秘钥。优先使用方式一一键配置；如遇异常，再使用方式二手动配置。</p>
                </article>
              </div>
            </section>
    
            <section id="purchase">
              <div class="section-head">
                <h2>购买与咨询</h2>
                <p class="section-copy">
                  扫码购买卡密或添加微信咨询。企业对接、大用量报价和稳定性需求建议先微信沟通。
                </p>
              </div>
              <div class="purchase-grid">
                <article class="qr-card">
                  <div class="qr-frame">
                    <img src="/店铺二维码.png" alt="店铺购买二维码" />
                  </div>
                  <div>
                    <h3>购买卡密</h3>
                    <p>扫码进入店铺购买卡密，购买后回到中转站内完成兑换。</p>
                    <a class="home-button" href="https://pay.ldxp.cn/shop/WNAXV9CI">打开购买地址</a>
                    <div class="mini">https://pay.ldxp.cn/shop/WNAXV9CI</div>
                  </div>
                </article>
                <article class="qr-card">
                  <div class="qr-frame">
                    <img src="/微信扫一扫.jpg" alt="微信咨询二维码" />
                  </div>
                  <div>
                    <h3>微信咨询</h3>
                    <p>企业对接、量大优惠、报价和模型调用问题，可添加微信沟通。</p>
                    <div class="home-button secondary">微信：tianheilezz</div>
                  </div>
                </article>
              </div>
            </section>
    
            <section id="notes">
              <div class="section-head">
                <h2>注意事项</h2>
                <p class="section-copy">这些细节会影响注册、兑换和调用体验，建议首次使用前先确认。</p>
              </div>
              <div class="notice">
                <div class="notice-item">
                  <strong>邮箱注册</strong>
                  <p>请使用邮箱注册。若注册邮件未收到，请检查垃圾箱或邮件拦截设置。</p>
                </div>
                <div class="notice-item">
                  <strong>秘钥分组</strong>
                  <p>创建秘钥时需要选择对应分组，避免额度、模型或路由策略不符合预期。</p>
                </div>
                <div class="notice-item">
                  <strong>模型调用</strong>
                  <p>GPT-5.6、GPT Image 2 等模型通过 OpenAI 兼容 API 对接，使用前请确认模型名和调用方式。</p>
                </div>
                <div class="notice-item">
                  <strong>配置优先级</strong>
                  <p>点击使用秘钥后，优先选择方式一一键配置；若客户端异常，再改用方式二。</p>
                </div>
              </div>
            </section>
    
            <section>
              <div class="enterprise">
                <div>
                  <h2>企业对接与大用量报价</h2>
                  <p>
                    欢迎企业、团队和高频 API 使用场景对接。量大可优惠，具体报价、并发、额度和稳定性需求可咨询微信：
                    <strong>tianheilezz</strong>。
                  </p>
                </div>
                <div class="hero-actions">
                  <a class="home-button" href="https://pay.ldxp.cn/shop/WNAXV9CI">购买卡密</a>
                  <a class="home-button secondary" href="/dashboard">进入中转站</a>
                </div>
              </div>
            </section>
          </main>
    
          <footer>
            <span>CyberRelay AI · OpenAI-compatible relay</span>
            <span>Register with email · Redeem code · Create API key · Configure</span>
          </footer>
        </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore, useAppStore } from '@/stores'

const authStore = useAuthStore()
const appStore = useAppStore()

const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

onMounted(() => {
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.cyber-home {
        color-scheme: dark;
        --bg: #0d0f10;
        --bg-top: #113a38;
        --panel: #141719;
        --panel-soft: #191d20;
        --panel-deep: #090a0b;
        --line: rgba(255, 255, 255, 0.08);
        --line-strong: rgba(255, 255, 255, 0.16);
        --text: #f1f0e8;
        --muted: #c0c2bc;
        --dim: #8d928d;
        --accent: #14b8a6;
        --accent-soft: rgba(20, 184, 166, 0.16);
        --button: #eeeee5;
        --button-text: #111315;
        --warn: #e5d06d;
      }
      .cyber-home {
        min-height: 100vh;
        background:
          linear-gradient(180deg, rgba(17, 58, 56, 0.92) 0, rgba(13, 15, 16, 0.98) 128px),
          var(--bg);
        color: var(--text);
        font-family:
          Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
          "PingFang SC", "Microsoft YaHei", sans-serif;
        letter-spacing: 0;
      }
.cyber-home a {
        color: inherit;
        text-decoration: none;
      }

      .cyber-page {
        width: min(1220px, calc(100% - 40px));
        margin: 0 auto;
      }

      .nav {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 82px;
      }

      .brand {
        display: flex;
        align-items: center;
        gap: 14px;
        min-width: 0;
      }

      .mark {
        display: grid;
        grid-template-columns: repeat(5, 6px);
        gap: 4px;
        opacity: 0.82;
      }

      .mark span {
        width: 6px;
        height: 6px;
        border-radius: 99px;
        background: var(--muted);
      }

      .brand-copy strong {
        display: block;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 14px;
        font-weight: 700;
        letter-spacing: 0;
      }

      .brand-copy small {
        display: block;
        margin-top: 2px;
        color: var(--dim);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 10px;
        text-transform: uppercase;
      }

      .nav-links {
        display: flex;
        align-items: center;
        gap: 28px;
        color: var(--muted);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 14px;
      }

      .home-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-height: 40px;
        border: 1px solid transparent;
        border-radius: 8px;
        padding: 0 18px;
        background: var(--button);
        color: var(--button-text) !important;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 14px;
        font-weight: 700;
        white-space: nowrap;
      }

      .home-button.secondary {
        background: transparent;
        color: var(--text) !important;
        border-color: var(--line-strong);
      }

      .hero {
        padding: 128px 0 72px;
      }

      .hero-grid {
        display: grid;
        grid-template-columns: minmax(0, 0.82fr) minmax(520px, 1.18fr);
        gap: 56px;
        align-items: end;
      }

      .eyebrow {
        margin: 0 0 18px;
        color: var(--dim);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
        text-transform: uppercase;
      }
.cyber-home h1 {
        max-width: 620px;
        margin: 0;
        color: #f2f1e9;
        font-size: clamp(34px, 3.6vw, 52px);
        font-weight: 750;
        line-height: 1.12;
        letter-spacing: 0;
      }
.cyber-home h1 span {
        display: block;
        white-space: nowrap;
      }

      .lead {
        max-width: 650px;
        margin: 26px 0 0;
        color: var(--muted);
        font-size: clamp(16px, 1.5vw, 20px);
        line-height: 1.65;
      }

      .lead-main {
        white-space: nowrap;
      }

      .lead-mobile {
        display: none;
      }

      .lead-list {
        display: grid;
        gap: 8px;
        max-width: 620px;
        margin: 22px 0 0;
        padding: 0;
        list-style: none;
      }

      .lead-list li {
        display: flex;
        gap: 10px;
        color: var(--muted);
        font-size: 16px;
        line-height: 1.65;
      }

      .lead-list li::before {
        content: "";
        flex: 0 0 auto;
        width: 6px;
        height: 6px;
        margin-top: 11px;
        border-radius: 999px;
        background: var(--accent);
      }

      .hero-actions {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
        margin-top: 34px;
      }

      .model-strip {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
        margin-top: 34px;
      }

      .pill {
        border: 1px solid var(--line);
        border-radius: 999px;
        padding: 8px 12px;
        background: rgba(255, 255, 255, 0.025);
        color: var(--muted);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
      }

      .product {
        border: 1px solid var(--line);
        border-radius: 10px;
        background: #050607;
        box-shadow: 0 34px 80px rgba(0, 0, 0, 0.42);
        overflow: hidden;
      }

      .product-top {
        display: flex;
        align-items: center;
        gap: 10px;
        height: 46px;
        padding: 0 14px;
        border-bottom: 1px solid var(--line);
        color: #a8aca6;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
      }

      .dots {
        display: flex;
        gap: 7px;
      }

      .dots span {
        width: 10px;
        height: 10px;
        border: 1px solid var(--line-strong);
        border-radius: 99px;
      }

      .product-body {
        display: grid;
        grid-template-columns: 180px 1fr 180px;
        min-height: 390px;
      }

      .gateway-sidebar,
      .inspector {
        padding: 18px;
        background: #08090a;
      }

      .gateway-sidebar {
        border-right: 1px solid var(--line);
      }

      .inspector {
        border-left: 1px solid var(--line);
      }

      .panel-title {
        margin: 0 0 14px;
        color: #9fa49e;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 11px;
        text-transform: uppercase;
      }

      .menu-item,
      .metric,
      .log-row {
        border: 1px solid var(--line);
        border-radius: 7px;
        background: rgba(255, 255, 255, 0.025);
      }

      .menu-item {
        display: block;
        margin-bottom: 8px;
        padding: 10px;
        color: #c1c5be;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
        transition:
          border-color 0.18s ease,
          background 0.18s ease,
          color 0.18s ease;
      }

      .menu-item.active {
        border-color: rgba(20, 184, 166, 0.38);
        background: var(--accent-soft);
        color: var(--text);
      }

      .menu-item:hover,
      .menu-item:focus-visible {
        border-color: rgba(20, 184, 166, 0.48);
        background: rgba(20, 184, 166, 0.11);
        color: var(--text);
        outline: none;
      }

      .workspace {
        padding: 18px;
      }

      .command {
        margin-bottom: 18px;
        border: 1px solid rgba(20, 184, 166, 0.28);
        border-radius: 8px;
        background: #050607;
        padding: 14px;
        color: #f0f0e8;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
        line-height: 1.7;
      }

      .accent {
        color: var(--accent);
      }

      .log-list {
        display: grid;
        gap: 10px;
      }

      .log-row {
        display: grid;
        grid-template-columns: 58px 1fr 76px;
        gap: 12px;
        align-items: center;
        padding: 12px;
        color: #c7cac4;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
      }

      .ok {
        color: var(--accent);
      }

      .metric {
        margin-bottom: 10px;
        padding: 12px;
      }

      .metric strong {
        display: block;
        color: var(--text);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 18px;
      }

      .metric span {
        display: block;
        margin-top: 4px;
        color: #a8aca6;
        font-size: 12px;
      }

      section {
        padding: 72px 0;
        border-top: 1px solid var(--line);
      }

      .section-head {
        display: flex;
        align-items: end;
        justify-content: space-between;
        gap: 28px;
        margin-bottom: 30px;
      }
.cyber-home h2 {
        margin: 0;
        color: #f2f1e9;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: clamp(28px, 3vw, 42px);
        line-height: 1.08;
      }

      .section-copy {
        max-width: 520px;
        margin: 0;
        color: var(--muted);
        line-height: 1.7;
      }

      .cards {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
        gap: 14px;
      }

      .card {
        min-height: 180px;
        border: 1px solid var(--line);
        border-radius: 10px;
        background: rgba(255, 255, 255, 0.025);
        padding: 24px;
      }

      .card .num {
        color: #9fa49e;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
      }

      .card h3 {
        margin: 24px 0 10px;
        color: var(--text);
        font-size: 20px;
      }

      .card p {
        margin: 0;
        color: var(--muted);
        line-height: 1.7;
      }

      .steps {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 12px;
      }

      .step {
        border-top: 1px solid var(--line-strong);
        padding-top: 18px;
      }

      .step .index {
        color: var(--accent);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 13px;
      }

      .step h3 {
        margin: 18px 0 8px;
        color: var(--text);
        font-size: 18px;
      }

      .step p {
        margin: 0;
        color: var(--muted);
        font-size: 14px;
        line-height: 1.7;
      }

      .notice {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 14px;
      }

      .purchase-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 14px;
      }

      .qr-card {
        display: grid;
        grid-template-columns: 148px 1fr;
        gap: 22px;
        align-items: center;
        border: 1px solid var(--line);
        border-radius: 10px;
        background: rgba(255, 255, 255, 0.025);
        padding: 18px;
      }

      .qr-frame {
        overflow: hidden;
        border: 1px solid var(--line-strong);
        border-radius: 8px;
        background: #f7f7f2;
        padding: 8px;
      }

      .qr-frame img {
        display: block;
        width: 100%;
        aspect-ratio: 1;
        object-fit: contain;
      }

      .qr-card h3 {
        margin: 0 0 10px;
        color: var(--text);
        font-size: 20px;
      }

      .qr-card p {
        margin: 0 0 16px;
        color: var(--muted);
        line-height: 1.7;
      }

      .qr-card .mini {
        color: #9fa49e;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
        word-break: break-all;
      }

      .notice-item {
        border: 1px solid var(--line);
        border-radius: 10px;
        background: #111315;
        padding: 22px;
      }

      .notice-item strong {
        display: block;
        margin-bottom: 8px;
        color: var(--warn);
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 13px;
      }

      .notice-item p {
        margin: 0;
        color: var(--muted);
        line-height: 1.7;
      }

      .enterprise {
        display: grid;
        grid-template-columns: 1fr auto;
        gap: 24px;
        align-items: center;
        border: 1px solid rgba(20, 184, 166, 0.3);
        border-radius: 12px;
        background: linear-gradient(135deg, rgba(20, 184, 166, 0.12), rgba(255, 255, 255, 0.02));
        padding: 30px;
      }

      .enterprise p {
        max-width: 760px;
        margin: 10px 0 0;
        color: var(--muted);
        line-height: 1.7;
      }

      footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 18px;
        padding: 38px 0 54px;
        color: #9fa49e;
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
        font-size: 12px;
      }

      @media (max-width: 980px) {
        .hero {
          padding-top: 70px;
        }

        .hero-grid,
        .product-body,
        .enterprise {
          grid-template-columns: 1fr;
        }

        .section-head {
          display: block;
        }

        .section-copy {
          margin-top: 14px;
        }

        .product-body {
          min-height: auto;
        }

        .gateway-sidebar,
        .inspector {
          display: none;
        }

        .cards,
        .steps,
        .notice,
        .purchase-grid {
          grid-template-columns: 1fr;
        }
      }

      @media (max-width: 680px) {
        .cyber-page {
          width: min(100% - 28px, 1220px);
        }

        .nav {
          height: auto;
          padding: 18px 0;
          align-items: flex-start;
        }

        .nav-links {
          gap: 10px;
        }

        .nav-links a:not(.home-button) {
          display: none;
        }

        .hero {
          padding: 54px 0 48px;
        }
.cyber-home h1 {
          font-size: 33px;
        }
.cyber-home h1 span {
          white-space: normal;
        }

        .lead-main {
          display: none;
        }

        .lead-mobile {
          display: inline;
        }

        .log-row {
          grid-template-columns: 48px 1fr;
        }

        .log-row span:last-child {
          display: none;
        }

        .qr-card {
          grid-template-columns: 112px 1fr;
          gap: 14px;
          padding: 14px;
        }

        footer {
          flex-direction: column;
          align-items: flex-start;
        }
      }
</style>
