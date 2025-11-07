import type { FormEvent } from 'react';
import { useEffect, useState } from 'react';
import { client} from '../api/client';
import type { Schedule, Content, Monitor, Location, Template } from '../api/client'


export default function SchedulesPage() {
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [contents, setContents] = useState<Content[]>([])
  const [monitors, setMonitors] = useState<Monitor[]>([])
  const [locations, setLocations] = useState<Location[]>([])
  const [templates, setTemplates] = useState<Template[]>([])
  const [selectedTemplateId, setSelectedTemplateId] = useState<number | null>(null)
  const [schedName, setSchedName] = useState('')
  const [dateStart, setDateStart] = useState('')
  const [dateEnd, setDateEnd] = useState('')
  const [repeatPattern, setRepeatPattern] = useState<'none'|'daily'|'weekly'|'monthly'>('none')
  const [daysOfWeek, setDaysOfWeek] = useState<Record<number, boolean>>({0:false,1:false,2:false,3:false,4:false,5:false,6:false})
  const [mode, setMode] = useState<'rotation'|'override'>('rotation')
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Schedule>({
    contentID: 0,
    monitorID: undefined,
    locationID: undefined,
    startTime: '',
    endTime: '',
    priority: 0,
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [resolveOpen, setResolveOpen] = useState(false)
  const [resolveItems, setResolveItems] = useState<Array<{
    id?: number
    contentID: number
    monitorID?: number
    locationID?: number
    startTime: string
    endTime: string
    priority?: number
    isNew?: boolean
    durationM?: number
  }>>([])

  const loadData = async () => {
    setLoading(true)
    try {
      const [schedulesData, contentsData, monitorsData, locationsData, templatesData] = await Promise.all([
        client.schedules.getAll(),
        client.contents.getAll(),
        client.monitors.getAll(),
        client.locations.getAll(),
        client.templates.getAll(),
      ])
      setSchedules(schedulesData)
      setContents(contentsData)
      setMonitors(monitorsData)
      setLocations(locationsData)
      setTemplates(templatesData)
      setError('')
    } catch (err: any) {
      setError('Ошибка загрузки данных: ' + (err.message || 'Неизвестная ошибка'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    // basic validation
    if (repeatPattern === 'weekly') {
      const anyDay = Object.values(daysOfWeek).some(Boolean)
      if (!anyDay) return setError('Выберите хотя бы один день недели для повторения')
    }
    try {
      // build payload: include schedule metadata and template reference
      const payload: any = {
        ...formData,
        name: schedName,
        templateID: selectedTemplateId || undefined,
        dateStart: dateStart || undefined,
        dateEnd: dateEnd || undefined,
        repeatPattern: repeatPattern === 'none' ? undefined : repeatPattern,
        daysOfWeek: repeatPattern === 'weekly' ? Object.keys(daysOfWeek).filter(k => daysOfWeek[Number(k)]).map(Number) : undefined,
        mode,
      }

      const result = await client.schedules.create(payload)

      if (result && typeof result === 'object' && 'error' in result && result.error) {
        // server returned conflicts — open resolve modal
        const conflicts: Schedule[] = (result as any).conflicts || []
        // include attempted schedule first
        const attempted = {
          ...formData,
          isNew: true,
        } as any

        const items = [attempted, ...conflicts].map((s: any) => {
          // compute duration in minutes
          let durationM = 0
          try {
            const st = new Date(s.startTime)
            const et = new Date(s.endTime)
            if (!isNaN(st.getTime()) && !isNaN(et.getTime())) {
              durationM = Math.max(1, Math.round((et.getTime() - st.getTime()) / 60000))
            }
          } catch {}
          return {
            id: s.id,
            contentID: s.contentID,
            monitorID: s.monitorID,
            locationID: s.locationID,
            startTime: s.startTime,
            endTime: s.endTime,
            priority: s.priority || 0,
            isNew: !!s.isNew,
            durationM,
          }
        })

        setResolveItems(items)
        setResolveOpen(true)
        return
      }

      setSuccess('Расписание успешно создано')
      setFormData({
        contentID: 0,
        monitorID: undefined,
        locationID: undefined,
        startTime: '',
        endTime: '',
        priority: 0,
      })
      setSchedName('')
      setDateStart('')
      setDateEnd('')
      setRepeatPattern('none')
      setDaysOfWeek({0:false,1:false,2:false,3:false,4:false,5:false,6:false})
      setMode('rotation')
      loadData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || err.message || 'Неизвестная ошибка'
      if (err.response?.status === 409) {
        setError(`Конфликт времени: ${errorMsg}`)
      } else {
        setError('Ошибка: ' + errorMsg)
      }
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить это расписание?')) return
    
    try {
      await client.schedules.delete(id)
      setSuccess('Расписание удалено')
      loadData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка удаления: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleUpdateSchedules = async () => {
    if (!confirm('Отправить все расписания на мониторы?')) return
    
    try {
      await client.schedules.updateSchedules(schedules)
      setSuccess('Расписания отправлены на мониторы')
      loadData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка: ' + (err.response?.data?.error || err.message))
    }
  }

  // helpers for resolve modal
  const updateResolveItem = (idx: number, patch: Partial<typeof resolveItems[number]>) => {
    setResolveItems((prev) => prev.map((it, i) => i === idx ? { ...it, ...patch } : it))
  }

  const computeEndTimeFromStart = (start: string, minutes: number) => {
    const d = new Date(start)
    if (isNaN(d.getTime())) return start
    d.setMinutes(d.getMinutes() + minutes)
    // return in format compatible with <input type="datetime-local">: YYYY-MM-DDTHH:mm
    const YYYY = d.getFullYear()
    const MM = String(d.getMonth() + 1).padStart(2, '0')
    const DD = String(d.getDate()).padStart(2, '0')
    const hh = String(d.getHours()).padStart(2, '0')
    const mm = String(d.getMinutes()).padStart(2, '0')
    return `${YYYY}-${MM}-${DD}T${hh}:${mm}`
  }

  // when user selects a template, prefill form (except startTime/endTime)
  useEffect(() => {
    if (!selectedTemplateId) return
    const t = templates.find(x => x.id === selectedTemplateId)
    if (!t) return
    // if template has blocks, pick first content of first block to prefill contentID
    if (Array.isArray(t.blocks) && t.blocks.length > 0) {
      const firstBlock = t.blocks[0]
      const firstContent = Array.isArray(firstBlock.contents) && firstBlock.contents.length > 0 ? firstBlock.contents[0] : null
      if (firstContent) {
        setFormData(prev => ({ ...prev, contentID: firstContent.contentID }))
      }
    }
  }, [selectedTemplateId, templates])

  // compute endTime automatically when startTime changes and a template with duration is selected
  useEffect(() => {
    if (!selectedTemplateId) return
    const t = templates.find(x => x.id === selectedTemplateId)
    if (!t || !Array.isArray(t.blocks) || t.blocks.length === 0) return
    if (!formData.startTime) return
    // compute total duration of first block (sum of content durations)
    const firstBlock = t.blocks[0]
    const totalSec = (firstBlock.contents || []).reduce((acc: number, c: any) => acc + (c.duration || 0), 0)
    const minutes = Math.max(1, Math.ceil(totalSec / 60))
    const end = computeEndTimeFromStart(formData.startTime, minutes)
    setFormData(prev => ({ ...prev, endTime: end }))
  }, [formData.startTime, selectedTemplateId, templates])

  const applyResolve = async () => {
    // build schedules to update
    const toUpdate: Schedule[] = resolveItems.map((it) => {
      const endTime = computeEndTimeFromStart(it.startTime, it.durationM || 0)
      return {
        id: it.id,
        contentID: it.contentID,
        monitorID: it.monitorID,
        locationID: it.locationID,
        startTime: it.startTime,
        endTime,
        priority: it.priority,
      } as Schedule
    })

    try {
      await client.schedules.updateSchedules(toUpdate)
      setResolveOpen(false)
      setResolveItems([])
      setSuccess('Конфликты разрешены и расписания обновлены')
      await loadData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка при разрешении конфликтов: ' + (err.response?.data?.error || err.message))
    }
  }

  const formatDateTime = (dateStr?: string) => {
    if (!dateStr) return '-'
    try {
      const date = new Date(dateStr)
      return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return dateStr
    }
  }

  const getContentName = (contentID: number) => {
    const content = contents.find(c => c.id === contentID)
    return content ? content.title : `ID: ${contentID}`
  }

  const getMonitorName = (monitorID?: number) => {
    if (!monitorID) return '-'
    const monitor = monitors.find(m => m.id === monitorID)
    return monitor ? monitor.name : `ID: ${monitorID}`
  }

  const getLocationName = (locationID?: number) => {
    if (!locationID) return '-'
    const location = locations.find(l => l.id === locationID)
    return location ? location.name : `ID: ${locationID}`
  }

  return (
    <div>
      <h1 className="page-title">Расписание</h1>

      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}

      {/* Create card */}
      <div className="card">
        <h2 style={{ marginBottom: '16px' }}>Создать новое расписание</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Шаблон (опционально)</label>
            <select value={selectedTemplateId ?? ''} onChange={(e) => setSelectedTemplateId(e.target.value ? Number(e.target.value) : null)}>
              <option value="">-- без шаблона --</option>
              {templates.map((t) => (
                <option key={t.id} value={t.id}>{t.name}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label>Контент *</label>
            <select
              value={formData.contentID}
              onChange={(e) => setFormData({ ...formData, contentID: Number(e.target.value) })}
              required
            >
              <option value={0}>Выберите контент</option>
              {contents.map((content, i) => (
                <option key={content.id ?? `content-${i}`} value={content.id}>
                  {content.title} ({content.type})
                </option>
              ))}
            </select>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Монитор</label>
              <select
                value={formData.monitorID || ''}
                onChange={(e) => setFormData({ ...formData, monitorID: e.target.value ? Number(e.target.value) : undefined })}
              >
                <option value="">Не выбран (для всех мониторов)</option>
                {monitors.map((monitor, i) => (
                  <option key={monitor.id ?? `monitor-${i}`} value={monitor.id}>
                    {monitor.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>Локация</label>
              <select
                value={formData.locationID || ''}
                onChange={(e) => setFormData({ ...formData, locationID: e.target.value ? Number(e.target.value) : undefined })}
              >
                <option value="">Не выбрана</option>
                {locations.map((location, i) => (
                  <option key={location.id ?? `loc-${i}`} value={location.id}>
                    {location.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Начало *</label>
              <input
                type="datetime-local"
                value={formData.startTime}
                onChange={(e) => setFormData({ ...formData, startTime: e.target.value })}
                required
              />
            </div>

            <div className="form-group">
              <label>Окончание *</label>
              <input
                type="datetime-local"
                value={formData.endTime}
                onChange={(e) => setFormData({ ...formData, endTime: e.target.value })}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label>Приоритет</label>
            <input
              type="number"
              value={formData.priority || 0}
              onChange={(e) => setFormData({ ...formData, priority: Number(e.target.value) })}
              min="0"
            />
          </div>

          <div style={{ marginTop: 12 }}>
            <h3>Дополнительно</h3>
            <div className="form-group">
              <label>Название расписания</label>
              <input value={schedName} onChange={(e) => setSchedName(e.target.value)} />
            </div>
            <div className="form-row">
              <div className="form-group">
                <label>Дата начала</label>
                <input type="date" value={dateStart} onChange={(e) => setDateStart(e.target.value)} />
              </div>
              <div className="form-group">
                <label>Дата конца</label>
                <input type="date" value={dateEnd} onChange={(e) => setDateEnd(e.target.value)} />
              </div>
            </div>
            <div className="form-group">
              <label>Режим</label>
              <select value={mode} onChange={(e) => setMode(e.target.value as any)}>
                <option value="rotation">rotation</option>
                <option value="override">override</option>
              </select>
            </div>
            <div className="form-group">
              <label>Повторение</label>
              <select value={repeatPattern} onChange={(e) => setRepeatPattern(e.target.value as any)}>
                <option value="none">none</option>
                <option value="daily">daily</option>
                <option value="weekly">weekly</option>
                <option value="monthly">monthly</option>
              </select>
            </div>
            {repeatPattern === 'weekly' && (
              <div className="form-group">
                <label>Дни недели</label>
                <div style={{ display: 'flex', gap: 8 }}>
                  {['Вс','Пн','Вт','Ср','Чт','Пт','Сб'].map((d, i) => (
                    <label key={i}><input type="checkbox" checked={daysOfWeek[i]} onChange={(e) => setDaysOfWeek(prev => ({ ...prev, [i]: e.target.checked }))} /> {d}</label>
                  ))}
                </div>
              </div>
            )}
            <div className="form-group">
              <label>Привязать к шаблону</label>
              <select value={selectedTemplateId ?? ''} onChange={(e) => setSelectedTemplateId(e.target.value ? Number(e.target.value) : null)}>
                <option value="">-- без шаблона --</option>
                {templates.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
              </select>
            </div>
            {selectedTemplateId && (
              <div style={{ marginTop: 8 }}>
                <h4>Блоки шаблона</h4>
                <ul>
                  {(templates.find(t => t.id === selectedTemplateId)?.blocks || []).map((b, i) => (
                    <li key={i}><strong>{b.name}</strong> {b.startTime}–{b.endTime} — содержимого: {(b.contents||[]).length}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          <div className="toolbar">
            <button type="submit" className="btn btn-primary">Создать</button>
          </div>
        </form>
      </div>

      {/* List card */}
      <div className="card">
        <div className="toolbar">
          <h2 style={{ margin: 0, flex: 1 }}>Список расписаний</h2>
          <button className="btn btn-success" onClick={handleUpdateSchedules} disabled={loading}>Отправить на мониторы</button>
          <button className="btn btn-secondary" onClick={loadData} disabled={loading}>{loading ? 'Загрузка...' : 'Обновить'}</button>
        </div>

        {schedules.length === 0 ? (
          <div className="empty-state"><p>Расписаний пока нет. Создайте первое расписание выше.</p></div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Контент</th>
                  <th>Монитор</th>
                  <th>Локация</th>
                  <th>Начало</th>
                  <th>Окончание</th>
                  <th>Приоритет</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {schedules.map((schedule, i) => (
                  <tr key={schedule.id ?? `schedule-${i}`}>
                    <td>{schedule.id}</td>
                    <td>{getContentName(schedule.contentID)}</td>
                    <td>{getMonitorName(schedule.monitorID)}</td>
                    <td>{getLocationName(schedule.locationID)}</td>
                    <td>{formatDateTime(schedule.startTime)}</td>
                    <td>{formatDateTime(schedule.endTime)}</td>
                    <td>{schedule.priority || 0}</td>
                    <td>
                      <button className="btn btn-danger" onClick={() => schedule.id && handleDelete(schedule.id)}>Удалить</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Resolve conflicts modal */}
      {resolveOpen && (
        <div className="modal-overlay" onClick={() => setResolveOpen(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3 style={{ margin: 0 }}>Разрешение конфликтов</h3>
              <button className="modal-close" onClick={() => setResolveOpen(false)}>×</button>
            </div>
            <div className="modal-body">
              <p>Найдены пересекающиеся расписания. Измените приоритет или длительность (в минутах), затем примените.</p>
              <div style={{ maxHeight: '50vh', overflow: 'auto' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                  <thead>
                    <tr>
                      <th>Контент</th>
                      <th>Монитор</th>
                      <th>Локация</th>
                      <th>Начало</th>
                      <th>Длительность (мин)</th>
                      <th>Приоритет</th>
                    </tr>
                  </thead>
                  <tbody>
                    {resolveItems.map((it, idx) => (
                      <tr key={it.id ?? `r-${idx}`}>
                        <td>{getContentName(it.contentID)}</td>
                        <td>{getMonitorName(it.monitorID)}</td>
                        <td>{getLocationName(it.locationID)}</td>
                        <td>{formatDateTime(it.startTime)}</td>
                        <td>
                          <input type="number" value={it.durationM || 0} min={1}
                            onChange={(e) => updateResolveItem(idx, { durationM: Number(e.target.value) })}
                            style={{ width: 100 }} />
                        </td>
                        <td>
                          <input type="number" value={it.priority || 0}
                            onChange={(e) => updateResolveItem(idx, { priority: Number(e.target.value) })}
                            style={{ width: 80 }} />
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 12 }}>
                <button className="btn btn-secondary" onClick={() => setResolveOpen(false)}>Отмена</button>
                <button className="btn btn-primary" onClick={applyResolve}>Применить</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
