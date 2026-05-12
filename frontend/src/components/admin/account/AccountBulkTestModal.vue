<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.bulkActions.testTitle')"
    width="wide"
    :close-on-click-outside="false"
    @close="handleClose"
  >
    <div class="space-y-4">
      <div class="text-sm text-gray-600 dark:text-gray-300">
        {{ t('admin.accounts.bulkActions.testDescription', { count: accounts.length }) }}
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <label class="text-sm text-gray-700 dark:text-gray-300">
          {{ t('admin.accounts.bulkActions.testConcurrency') }}
        </label>
        <select
          v-model.number="concurrency"
          :disabled="running"
          class="rounded-md border border-gray-300 bg-white px-2 py-1 text-sm dark:border-dark-500 dark:bg-dark-700 dark:text-gray-100"
        >
          <option v-for="n in [1, 2, 3, 5, 8]" :key="n" :value="n">{{ n }}</option>
        </select>

        <span class="ml-auto text-xs text-gray-500 dark:text-gray-400">
          {{ t('admin.accounts.bulkActions.testProgress', { done: doneCount, total: results.length }) }}
        </span>

        <button
          v-if="!running"
          class="btn btn-primary btn-sm"
          :disabled="results.length === 0"
          @click="startAll"
        >
          {{ t('admin.accounts.bulkActions.testStart') }}
        </button>
        <button v-else class="btn btn-warning btn-sm" @click="cancelAll">
          {{ t('admin.accounts.bulkActions.testCancel') }}
        </button>
      </div>

      <div
        v-if="finished"
        class="rounded-md border border-gray-200 bg-gray-50 p-2 text-xs text-gray-700 dark:border-dark-500 dark:bg-dark-700 dark:text-gray-200"
      >
        {{
          t('admin.accounts.bulkActions.testSummary', {
            pass: passCount,
            fail: failCount,
            skipped: skippedCount
          })
        }}
      </div>

      <div class="max-h-[420px] overflow-auto rounded-lg border border-gray-200 dark:border-dark-500">
        <table class="min-w-full text-sm">
          <thead class="bg-gray-50 dark:bg-dark-700">
            <tr class="text-left text-xs uppercase text-gray-500 dark:text-gray-400">
              <th class="px-3 py-2 font-medium">{{ t('admin.accounts.bulkActions.testColAccount') }}</th>
              <th class="px-3 py-2 font-medium">{{ t('admin.accounts.bulkActions.testColPlatform') }}</th>
              <th class="px-3 py-2 font-medium">{{ t('admin.accounts.bulkActions.testColStatus') }}</th>
              <th class="px-3 py-2 font-medium">{{ t('admin.accounts.bulkActions.testColMessage') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in results"
              :key="row.id"
              class="border-t border-gray-100 dark:border-dark-500"
            >
              <td class="px-3 py-2">
                <div class="font-medium text-gray-900 dark:text-gray-100">{{ row.name }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">#{{ row.id }}</div>
              </td>
              <td class="px-3 py-2 text-gray-700 dark:text-gray-200">
                <span class="rounded bg-gray-100 px-1.5 py-0.5 text-xs uppercase dark:bg-dark-600">
                  {{ row.platform }}
                </span>
              </td>
              <td class="px-3 py-2">
                <span :class="['rounded px-2 py-0.5 text-xs font-semibold', statusBadgeClass(row.status)]">
                  {{ statusLabel(row.status) }}
                </span>
              </td>
              <td class="px-3 py-2 max-w-[420px] break-words text-xs text-gray-600 dark:text-gray-300">
                {{ row.message || '—' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'

type Account = {
  id: number
  name: string
  platform: string
}

type RowStatus =
  | 'pending'
  | 'running'
  | 'pass'
  | 'fail'
  | 'unauthorized'
  | 'rate_limited'
  | 'timeout'
  | 'skipped'

interface Row {
  id: number
  name: string
  platform: string
  status: RowStatus
  message: string
}

const props = defineProps<{ show: boolean; accounts: Account[] }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'finished'): void }>()

const { t } = useI18n()

const concurrency = ref(3)
const results = ref<Row[]>([])
const running = ref(false)
const cancelled = ref(false)
const controllers = new Set<AbortController>()

const doneCount = computed(
  () => results.value.filter((r) => r.status !== 'pending' && r.status !== 'running').length
)
const passCount = computed(() => results.value.filter((r) => r.status === 'pass').length)
const failCount = computed(
  () =>
    results.value.filter((r) =>
      ['fail', 'unauthorized', 'rate_limited', 'timeout'].includes(r.status)
    ).length
)
const skippedCount = computed(() => results.value.filter((r) => r.status === 'skipped').length)
const finished = computed(() => !running.value && doneCount.value === results.value.length && results.value.length > 0)

watch(
  () => props.show,
  (val) => {
    if (val) {
      // Reset state every time the modal opens
      cancelAll()
      results.value = props.accounts.map((a) => ({
        id: a.id,
        name: a.name,
        platform: a.platform,
        status: 'pending' as RowStatus,
        message: ''
      }))
    }
  },
  { immediate: true }
)

const handleClose = () => {
  if (running.value) {
    cancelAll()
  }
  emit('close')
  if (finished.value) {
    emit('finished')
  }
}

function cancelAll() {
  cancelled.value = true
  for (const c of controllers) {
    try {
      c.abort()
    } catch {
      /* ignore */
    }
  }
  controllers.clear()
  running.value = false
  // mark anything still running as skipped
  for (const r of results.value) {
    if (r.status === 'running' || r.status === 'pending') {
      r.status = 'skipped'
      r.message = ''
    }
  }
}

const startAll = async () => {
  if (running.value || results.value.length === 0) return
  cancelled.value = false
  running.value = true
  // Reset previous run results but keep order
  for (const r of results.value) {
    r.status = 'pending'
    r.message = ''
  }

  const queue = [...results.value]
  const workers = Math.max(1, Math.min(concurrency.value, queue.length))

  const runWorker = async () => {
    while (!cancelled.value) {
      const row = queue.shift()
      if (!row) return
      await runOne(row)
    }
  }

  await Promise.all(Array.from({ length: workers }, () => runWorker()))
  running.value = false
  if (!cancelled.value) {
    emit('finished')
  }
}

const runOne = async (row: Row) => {
  row.status = 'running'
  row.message = ''
  const ctrl = new AbortController()
  controllers.add(ctrl)
  try {
    const resp = await fetch(`/api/v1/admin/accounts/${row.id}/test`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('auth_token') || ''}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({}),
      signal: ctrl.signal
    })

    if (!resp.ok) {
      if (resp.status === 401) {
        row.status = 'unauthorized'
      } else if (resp.status === 429) {
        row.status = 'rate_limited'
      } else {
        row.status = 'fail'
      }
      row.message = `HTTP ${resp.status}`
      return
    }

    const reader = resp.body?.getReader()
    if (!reader) {
      row.status = 'fail'
      row.message = 'No response body'
      return
    }

    const decoder = new TextDecoder()
    let buffer = ''
    let finalStatus: RowStatus | null = null
    let finalMessage = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''
      for (const line of lines) {
        if (!line.startsWith('data: ')) continue
        const jsonStr = line.slice(6).trim()
        if (!jsonStr) continue
        try {
          const ev = JSON.parse(jsonStr) as {
            type?: string
            success?: boolean
            error?: string
          }
          if (ev.type === 'test_complete') {
            if (ev.success) {
              finalStatus = 'pass'
              finalMessage = ''
            } else {
              finalStatus = classifyError(ev.error)
              finalMessage = ev.error || ''
            }
          } else if (ev.type === 'error') {
            finalStatus = classifyError(ev.error)
            finalMessage = ev.error || ''
          }
        } catch {
          // ignore parse errors
        }
      }
    }

    if (!finalStatus) {
      row.status = 'fail'
      row.message = 'No completion event received'
    } else {
      row.status = finalStatus
      row.message = finalMessage
    }
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      row.status = 'skipped'
      row.message = ''
    } else {
      row.status = 'fail'
      row.message = err instanceof Error ? err.message : String(err)
    }
  } finally {
    controllers.delete(ctrl)
  }
}

const classifyError = (msg?: string): RowStatus => {
  const m = (msg || '').toLowerCase()
  if (/(^|[^0-9])401([^0-9]|$)|unauthor/.test(m)) return 'unauthorized'
  if (/(^|[^0-9])429([^0-9]|$)|rate.?limit|too many request/.test(m)) return 'rate_limited'
  if (/timeout|timed out/.test(m)) return 'timeout'
  return 'fail'
}

const statusLabel = (s: RowStatus): string => {
  const map: Record<RowStatus, string> = {
    pending: t('admin.accounts.bulkActions.testStatusPending'),
    running: t('admin.accounts.bulkActions.testStatusRunning'),
    pass: t('admin.accounts.bulkActions.testStatusPass'),
    fail: t('admin.accounts.bulkActions.testStatusFail'),
    unauthorized: t('admin.accounts.bulkActions.testStatusUnauthorized'),
    rate_limited: t('admin.accounts.bulkActions.testStatusRateLimited'),
    timeout: t('admin.accounts.bulkActions.testStatusTimeout'),
    skipped: t('admin.accounts.bulkActions.testStatusSkipped')
  }
  return map[s]
}

const statusBadgeClass = (s: RowStatus): string => {
  switch (s) {
    case 'pass':
      return 'bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-400'
    case 'fail':
      return 'bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-400'
    case 'unauthorized':
      return 'bg-orange-100 text-orange-700 dark:bg-orange-500/20 dark:text-orange-400'
    case 'rate_limited':
      return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/20 dark:text-yellow-300'
    case 'timeout':
      return 'bg-purple-100 text-purple-700 dark:bg-purple-500/20 dark:text-purple-300'
    case 'running':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-300'
    case 'skipped':
      return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
    default:
      return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
  }
}
</script>
