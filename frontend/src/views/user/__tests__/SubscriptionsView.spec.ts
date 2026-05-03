import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { nextTick } from 'vue'

import SubscriptionsView from '../SubscriptionsView.vue'

const routerPush = vi.hoisted(() => vi.fn())
const getMySubscriptions = vi.hoisted(() => vi.fn())
const resetQuota = vi.hoisted(() => vi.fn())
const showSuccess = vi.hoisted(() => vi.fn())
const showError = vi.hoisted(() => vi.fn())

vi.mock('vue-router', async () => {
  const actual = await vi.importActual<typeof import('vue-router')>('vue-router')
  return {
    ...actual,
    useRouter: () => ({
      push: routerPush
    })
  }
})

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

vi.mock('@/api/subscriptions', () => ({
  default: {
    getMySubscriptions,
    resetQuota
  },
  getMySubscriptions,
  resetQuota
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showSuccess,
    showError
  })
}))

const AppLayoutStub = {
  template: '<div><slot /></div>'
}

const ConfirmDialogStub = {
  props: ['show', 'title', 'message', 'confirmText', 'cancelText', 'danger'],
  emits: ['confirm', 'cancel'],
  template: `
    <div v-if="show" class="confirm-dialog-stub">
      <button class="confirm" @click="$emit('confirm')">{{ confirmText }}</button>
      <button class="cancel" @click="$emit('cancel')">{{ cancelText }}</button>
      <p class="message">{{ message }}</p>
    </div>
  `
}

describe('SubscriptionsView', () => {
  beforeEach(() => {
    routerPush.mockReset()
    getMySubscriptions.mockReset()
    resetQuota.mockReset()
    showSuccess.mockReset()
    showError.mockReset()
  })

  it('lets the user spend one day to reset today\'s usage', async () => {
    getMySubscriptions
      .mockResolvedValueOnce([
        {
          id: 501,
          user_id: 1,
          group_id: 10,
          starts_at: '2099-01-01T00:00:00Z',
          expires_at: '2099-01-03T00:00:00Z',
          status: 'active',
          daily_window_start: '2099-01-02T00:00:00Z',
          weekly_window_start: null,
          monthly_window_start: null,
          daily_usage_usd: 12.34,
          weekly_usage_usd: 0,
          monthly_usage_usd: 0,
          created_at: '2099-01-01T00:00:00Z',
          updated_at: '2099-01-02T00:00:00Z',
          group: {
            id: 10,
            name: 'Group One',
            daily_limit_usd: 100,
            weekly_limit_usd: null,
            monthly_limit_usd: null,
            platform: 'openai'
          }
        }
      ])
      .mockResolvedValueOnce([
        {
          id: 501,
          user_id: 1,
          group_id: 10,
          starts_at: '2099-01-01T00:00:00Z',
          expires_at: '2099-01-02T00:00:00Z',
          status: 'active',
          daily_window_start: '2099-01-02T00:00:00Z',
          weekly_window_start: null,
          monthly_window_start: null,
          daily_usage_usd: 0,
          weekly_usage_usd: 0,
          monthly_usage_usd: 0,
          created_at: '2099-01-01T00:00:00Z',
          updated_at: '2099-01-02T00:00:00Z',
          group: {
            id: 10,
            name: 'Group One',
            daily_limit_usd: 100,
            weekly_limit_usd: null,
            monthly_limit_usd: null,
            platform: 'openai'
          }
        }
      ])
    resetQuota.mockResolvedValueOnce({
      id: 501
    })

    const wrapper = mount(SubscriptionsView, {
      global: {
        stubs: {
          AppLayout: AppLayoutStub,
          ConfirmDialog: ConfirmDialogStub,
          Icon: true
        }
      }
    })

    await flushPromises()
    await nextTick()

    const resetButton = wrapper
      .findAll('button')
      .find((button) => button.text().includes('userSubscriptions.resetQuota'))
    expect(resetButton).toBeTruthy()

    await resetButton!.trigger('click')
    await nextTick()

    const confirmButton = wrapper.get('button.confirm')
    await confirmButton.trigger('click')

    await flushPromises()
    await nextTick()

    expect(resetQuota).toHaveBeenCalledWith(501)
    expect(showSuccess).toHaveBeenCalledWith('userSubscriptions.resetQuotaSuccess')
    expect(showError).not.toHaveBeenCalled()
    expect(getMySubscriptions).toHaveBeenCalledTimes(2)
    expect(wrapper.text()).toContain('$0.00 / $100.00')
  })
})
