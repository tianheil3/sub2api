<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.bulkActions.scheduleTitle')"
    width="wide"
    :close-on-click-outside="false"
    @close="handleClose"
  >
    <div class="space-y-4">
      <p class="text-sm text-gray-600 dark:text-gray-300">
        {{ t('admin.accounts.bulkActions.scheduleDescription', { count: accountIds.length }) }}
      </p>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">
            {{ t('admin.scheduledTests.model') }}
          </label>
          <Input
            v-model="form.model_id"
            :placeholder="t('admin.accounts.bulkActions.scheduleModelPlaceholder')"
            :disabled="creating"
          />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">
            {{ t('admin.scheduledTests.cronExpression') }}
          </label>
          <Input
            v-model="form.cron_expression"
            :placeholder="'*/30 * * * *'"
            :hint="t('admin.scheduledTests.cronHelp')"
            :disabled="creating"
          />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">
            {{ t('admin.scheduledTests.maxResults') }}
          </label>
          <Input
            v-model="form.max_results"
            type="number"
            placeholder="100"
            :disabled="creating"
          />
        </div>
        <div class="grid gap-3 sm:grid-cols-3">
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
            <Toggle v-model="form.enabled" />
            {{ t('admin.scheduledTests.enabled') }}
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
            <Toggle v-model="form.auto_recover" />
            {{ t('admin.scheduledTests.autoRecover') }}
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
            <Toggle v-model="form.auto_disable_on_unauth" />
            {{ t('admin.scheduledTests.autoDisableOnUnauth') }}
          </label>
        </div>
      </div>

      <div
        v-if="summary"
        class="rounded-md border border-gray-200 bg-gray-50 p-2 text-xs text-gray-700 dark:border-dark-500 dark:bg-dark-700 dark:text-gray-200"
      >
        {{ t('admin.accounts.bulkActions.scheduleSummary', summary) }}
      </div>

      <div v-if="failures.length > 0" class="max-h-40 overflow-auto rounded-md border border-red-200 bg-red-50 p-2 text-xs text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-300">
        <div v-for="failure in failures" :key="failure.accountId">
          #{{ failure.accountId }}: {{ failure.message }}
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <button
          class="btn btn-secondary btn-sm"
          :disabled="creating"
          @click="handleClose"
        >
          {{ t('common.cancel') }}
        </button>
        <button
          class="btn btn-primary btn-sm"
          :disabled="creating || accountIds.length === 0 || !form.cron_expression"
          @click="handleCreate"
        >
          {{ creating ? t('common.loading') : t('admin.accounts.bulkActions.scheduleCreate') }}
        </button>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Input from '@/components/common/Input.vue'
import Toggle from '@/components/common/Toggle.vue'
import { adminAPI } from '@/api/admin'
import { useAppStore } from '@/stores/app'

const props = defineProps<{ show: boolean; accountIds: number[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'created'): void }>()

const { t } = useI18n()
const appStore = useAppStore()

const form = reactive({
  model_id: '',
  cron_expression: '*/30 * * * *',
  max_results: '100',
  enabled: true,
  auto_recover: false,
  auto_disable_on_unauth: true
})
const creating = ref(false)
const summary = ref<{ success: number; failed: number } | null>(null)
const failures = ref<Array<{ accountId: number; message: string }>>([])

watch(
  () => props.show,
  (visible) => {
    if (visible) {
      summary.value = null
      failures.value = []
    }
  }
)

const handleClose = () => {
  if (!creating.value) {
    emit('close')
  }
}

const handleCreate = async () => {
  if (creating.value || props.accountIds.length === 0 || !form.cron_expression) return
  creating.value = true
  summary.value = null
  failures.value = []

  let success = 0
  for (const accountId of props.accountIds) {
    try {
      await adminAPI.scheduledTests.create({
        account_id: accountId,
        model_id: form.model_id.trim(),
        cron_expression: form.cron_expression.trim(),
        enabled: form.enabled,
        max_results: Number(form.max_results) || 100,
        auto_recover: form.auto_recover,
        auto_disable_on_unauth: form.auto_disable_on_unauth
      })
      success += 1
    } catch (error: any) {
      failures.value.push({
        accountId,
        message: error?.message || t('admin.accounts.bulkActions.scheduleCreateFailed')
      })
    }
  }

  const failed = props.accountIds.length - success
  summary.value = { success, failed }
  if (failed === 0) {
    appStore.showSuccess(t('admin.accounts.bulkActions.scheduleCreateSuccess', { count: success }))
    emit('created')
  } else {
    appStore.showError(t('admin.accounts.bulkActions.scheduleCreatePartial', { success, failed }))
  }
  creating.value = false
}
</script>
