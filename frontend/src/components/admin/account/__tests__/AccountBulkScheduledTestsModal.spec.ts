import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, describe, expect, it, vi } from 'vitest'
import AccountBulkScheduledTestsModal from '../AccountBulkScheduledTestsModal.vue'

const { createPlan, showError, showSuccess } = vi.hoisted(() => ({
  createPlan: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    scheduledTests: {
      create: createPlan
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess
  })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

describe('AccountBulkScheduledTestsModal', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('creates the same 401 auto-pause scheduled plan for each selected account', async () => {
    createPlan.mockResolvedValue({})

    const wrapper = mount(AccountBulkScheduledTestsModal, {
      props: {
        show: true,
        accountIds: [11, 12]
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' },
          Input: {
            props: ['modelValue'],
            emits: ['update:modelValue'],
            template:
              '<input class="input-stub" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
          },
          Toggle: false
        }
      }
    })

    const inputs = wrapper.findAll('input.input-stub')
    await inputs[0].setValue('gpt-5.5')
    await inputs[1].setValue('*/15 * * * *')

    const createButton = wrapper.findAll('button').find((button) => button.text().includes('admin.accounts.bulkActions.scheduleCreate'))
    expect(createButton).toBeTruthy()
    await createButton!.trigger('click')
    await flushPromises()

    expect(createPlan).toHaveBeenCalledTimes(2)
    expect(createPlan).toHaveBeenNthCalledWith(
      1,
      expect.objectContaining({
        account_id: 11,
        model_id: 'gpt-5.5',
        cron_expression: '*/15 * * * *',
        auto_disable_on_unauth: true
      })
    )
    expect(createPlan).toHaveBeenNthCalledWith(
      2,
      expect.objectContaining({
        account_id: 12,
        model_id: 'gpt-5.5',
        cron_expression: '*/15 * * * *',
        auto_disable_on_unauth: true
      })
    )
    expect(wrapper.emitted('created')).toBeTruthy()
  })
})
