import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import ScheduledTestsPanel from '../ScheduledTestsPanel.vue'

const { createPlan, updatePlan, listByAccount, listResults, deletePlan, showError, showSuccess } =
  vi.hoisted(() => ({
    createPlan: vi.fn(),
    updatePlan: vi.fn(),
    listByAccount: vi.fn(),
    listResults: vi.fn(),
    deletePlan: vi.fn(),
    showError: vi.fn(),
    showSuccess: vi.fn()
  }))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    scheduledTests: {
      listByAccount,
      create: createPlan,
      update: updatePlan,
      delete: deletePlan,
      listResults
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess
  })
}))

vi.mock('@/utils/format', () => ({
  formatDateTime: (value: string) => value
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

const selectStub = {
  props: ['modelValue', 'options'],
  emits: ['update:modelValue'],
  template:
    '<select class="select-stub" :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option></select>'
}

const inputStub = {
  props: ['modelValue', 'type'],
  emits: ['update:modelValue'],
  template:
    '<input class="input-stub" :type="type || \'text\'" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
}

function mountPanel() {
  return mount(ScheduledTestsPanel, {
    props: {
      show: false,
      accountId: 42,
      modelOptions: [{ value: 'gpt-test', label: 'GPT Test' }]
    },
    global: {
      stubs: {
        BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' },
        ConfirmDialog: true,
        HelpTooltip: { template: '<span><slot name="trigger" /><slot /></span>' },
        Select: selectStub,
        Input: inputStub,
        Toggle: false,
        Icon: true
      }
    }
  })
}

describe('ScheduledTestsPanel', () => {
  beforeEach(() => {
    listResults.mockResolvedValue([])
    createPlan.mockResolvedValue({})
    updatePlan.mockResolvedValue({})
    deletePlan.mockResolvedValue(undefined)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('includes auto_disable_on_unauth in create payload when enabled', async () => {
    listByAccount.mockResolvedValue([])
    const wrapper = mountPanel()
    await wrapper.setProps({ show: true })
    await flushPromises()

    await wrapper.find('button.btn-primary').trigger('click')
    await flushPromises()

    await wrapper.find('select.select-stub').setValue('gpt-test')
    const inputs = wrapper.findAll('input.input-stub')
    await inputs[0].setValue('*/15 * * * *')

    const toggles = wrapper.findAll('button[role="switch"]')
    await toggles[2].trigger('click')

    const saveButton = wrapper.findAll('button').find((button) => button.text().includes('common.save'))
    expect(saveButton).toBeTruthy()
    await saveButton!.trigger('click')
    await flushPromises()

    expect(createPlan).toHaveBeenCalledWith(
      expect.objectContaining({
        account_id: 42,
        model_id: 'gpt-test',
        cron_expression: '*/15 * * * *',
        auto_disable_on_unauth: true
      })
    )
  })

  it('renders the 401 auto-pause badge and includes the field in update payload', async () => {
    listByAccount.mockResolvedValue([
      {
        id: 7,
        account_id: 42,
        model_id: 'gpt-test',
        cron_expression: '*/30 * * * *',
        enabled: true,
        max_results: 100,
        auto_recover: false,
        auto_disable_on_unauth: true,
        last_run_at: null,
        next_run_at: null,
        created_at: '2026-05-10T00:00:00Z',
        updated_at: '2026-05-10T00:00:00Z'
      }
    ])
    updatePlan.mockResolvedValue({
      id: 7,
      account_id: 42,
      model_id: 'gpt-test',
      cron_expression: '*/30 * * * *',
      enabled: true,
      max_results: 100,
      auto_recover: false,
      auto_disable_on_unauth: false,
      last_run_at: null,
      next_run_at: null,
      created_at: '2026-05-10T00:00:00Z',
      updated_at: '2026-05-10T00:00:00Z'
    })

    const wrapper = mountPanel()
    await wrapper.setProps({ show: true })
    await flushPromises()

    expect(wrapper.text()).toContain('admin.scheduledTests.autoDisableOnUnauthBadge')

    const editButton = wrapper.find('button[title="admin.scheduledTests.editPlan"]')
    await editButton.trigger('click')
    await flushPromises()

    const toggles = wrapper.findAll('button[role="switch"]')
    await toggles[toggles.length - 1].trigger('click')

    const saveButton = wrapper.findAll('button').find((button) => button.text().includes('common.save'))
    expect(saveButton).toBeTruthy()
    await saveButton!.trigger('click')
    await flushPromises()

    expect(updatePlan).toHaveBeenCalledWith(
      7,
      expect.objectContaining({
        auto_disable_on_unauth: false
      })
    )
  })
})
