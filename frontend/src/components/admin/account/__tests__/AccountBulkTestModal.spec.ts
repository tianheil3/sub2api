import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import AccountBulkTestModal from '../AccountBulkTestModal.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string | number>) => {
        if (key === 'admin.accounts.bulkActions.testProgress') {
          return `${params?.done}/${params?.total}`
        }
        return key
      }
    })
  }
})

function createStreamResponse(events: Array<Record<string, unknown>>) {
  const encoder = new TextEncoder()
  const chunks = events.map((event) => encoder.encode(`data: ${JSON.stringify(event)}\n`))
  let index = 0

  return {
    ok: true,
    body: {
      getReader: () => ({
        read: vi.fn().mockImplementation(async () => {
          if (index < chunks.length) {
            return { done: false, value: chunks[index++] }
          }
          return { done: true, value: undefined }
        })
      })
    }
  } as Response
}

function mountModal(accounts = [
  { id: 1, name: 'Account 1', platform: 'openai' },
  { id: 2, name: 'Account 2', platform: 'claude' },
  { id: 3, name: 'Account 3', platform: 'gemini' }
]) {
  return mount(AccountBulkTestModal, {
    props: {
      show: true,
      accounts
    },
    global: {
      stubs: {
        BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' }
      }
    }
  })
}

describe('AccountBulkTestModal', () => {
  beforeEach(() => {
    Object.defineProperty(globalThis, 'localStorage', {
      value: {
        getItem: vi.fn((key: string) => (key === 'auth_token' ? 'test-token' : null)),
        setItem: vi.fn(),
        removeItem: vi.fn(),
        clear: vi.fn()
      },
      configurable: true
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('maps SSE completion and error events to visible row statuses', async () => {
    global.fetch = vi
      .fn()
      .mockResolvedValueOnce(createStreamResponse([{ type: 'test_complete', success: true }]))
      .mockResolvedValueOnce(createStreamResponse([{ type: 'error', error: 'HTTP 401: unauthorized' }]))
      .mockResolvedValueOnce(createStreamResponse([{ type: 'test_complete', success: false, error: 'rate limited 429' }])) as any

    const wrapper = mountModal()
    await wrapper.find('button.btn-primary').trigger('click')
    await flushPromises()
    await flushPromises()

    expect(wrapper.text()).toContain('admin.accounts.bulkActions.testStatusPass')
    expect(wrapper.text()).toContain('admin.accounts.bulkActions.testStatusUnauthorized')
    expect(wrapper.text()).toContain('admin.accounts.bulkActions.testStatusRateLimited')
  })

  it('maps HTTP 401 response to unauthorized status', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 401
    }) as any

    const wrapper = mountModal([{ id: 1, name: 'Account 1', platform: 'openai' }])
    await wrapper.find('button.btn-primary').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('admin.accounts.bulkActions.testStatusUnauthorized')
    expect(wrapper.text()).toContain('HTTP 401')
  })

  it('honors the selected concurrency limit', async () => {
    let active = 0
    let maxActive = 0
    const resolvers: Array<() => void> = []

    global.fetch = vi.fn().mockImplementation(() => {
      active += 1
      maxActive = Math.max(maxActive, active)
      return new Promise((resolve) => {
        resolvers.push(() => {
          active -= 1
          resolve(createStreamResponse([{ type: 'test_complete', success: true }]))
        })
      })
    }) as any

    const wrapper = mountModal([
      { id: 1, name: 'Account 1', platform: 'openai' },
      { id: 2, name: 'Account 2', platform: 'openai' },
      { id: 3, name: 'Account 3', platform: 'openai' },
      { id: 4, name: 'Account 4', platform: 'openai' }
    ])

    await wrapper.find('select').setValue('2')
    const running = wrapper.find('button.btn-primary').trigger('click')
    await flushPromises()
    expect(maxActive).toBe(2)

    while (resolvers.length > 0) {
      resolvers.shift()?.()
      await flushPromises()
    }
    await running

    expect(global.fetch).toHaveBeenCalledTimes(4)
    expect(maxActive).toBe(2)
  })
})
