import { useEffect, useState } from 'react';
import { client} from '../api/client';
import type { Schedule, Content, Monitor, Location, Template } from '../api/client'
import ScheduleModal from '../components/ScheduleModal'


export default function SchedulesPage() {
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [contents, setContents] = useState<Content[]>([])
  const [monitors, setMonitors] = useState<Monitor[]>([])
  const [locations, setLocations] = useState<Location[]>([])
  const [templates, setTemplates] = useState<Template[]>([])
  const [loading, setLoading] = useState(false)
  // form state is handled in the modal component
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
  const [showModal, setShowModal] = useState(false)

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

  // monitor/target details modal
  const [targetDetails, setTargetDetails] = useState<null | {
    scheduleId?: number
    location?: Location | null
    groupName?: string | null
    monitors?: Array<{ id?: number; name?: string }>
  }>(null)

  const openTargetDetails = (s: Schedule) => {
    const loc = (s as any).location ?? (s.locationID ? locations.find(l => l.id === s.locationID) : null)
    const groupName = (s as any).group?.name ?? (s.groupId ? `Группа ${s.groupId}` : null)
    let mons: Array<{ id?: number; name?: string }> = []
    if (Array.isArray((s as any).monitors) && (s as any).monitors.length > 0) {
      mons = (s as any).monitors.map((m: any) => ({ id: m.id, name: m.name || monitors.find(mm => mm.id === m.id)?.name }))
    } else if (s.monitorID) {
      const m = monitors.find(mm => mm.id === s.monitorID)
      if (m) mons = [{ id: m.id, name: m.name }]
    }

    setTargetDetails({ scheduleId: s.id, location: loc || null, groupName: groupName || null, monitors: mons })
  }

  const closeTargetDetails = () => setTargetDetails(null)

  const getContentName = (arg: Schedule | number) => {
    if (typeof arg === 'number') {
      const content = contents.find(c => c.id === arg)
      return content ? content.title : (typeof arg === 'number' ? `ID: ${arg}` : '-')
    }
    const s = arg as Schedule
    if (s.content && s.content.title) return s.content.title
    if ((s as any).name) return (s as any).name
    if (typeof s.contentID !== 'undefined' && s.contentID !== null) return `ID: ${s.contentID}`
    return '-'
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

      {/* Create button card (opens modal) */}
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>Создать новое расписание</h2>
          <div>
            <button className="btn btn-primary" onClick={() => setShowModal(true)}>Создать расписание</button>
          </div>
        </div>
      </div>

      <ScheduleModal
        visible={showModal}
        onClose={() => setShowModal(false)}
        onCreated={() => { setShowModal(false); loadData() }}
        monitors={monitors}
        locations={locations}
        templates={templates}
      />

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
                                    <td>{getContentName(schedule)}</td>
                                    <td>
                                      <button className="link" onClick={() => openTargetDetails(schedule)} style={{ background: 'none', border: 'none', padding: 0, color: '#0366d6', cursor: 'pointer' }}>
                                        {(() => {
                                          // compact summary
                                          if ((schedule as any).location?.name) return (schedule as any).location.name
                                          if (schedule.locationID) {
                                            const l = locations.find(ll => ll.id === schedule.locationID)
                                            if (l) return l.name
                                          }
                                          if ((schedule as any).group?.name) return (schedule as any).group.name
                                          if ((schedule as any).monitors && (schedule as any).monitors.length > 0) {
                                            const mons = (schedule as any).monitors
                                            // show first monitor name or count
                                            const first = mons[0]
                                            const firstName = first?.name || monitors.find(mm => mm.id === first?.id)?.name
                                            return firstName ? firstName + (mons.length > 1 ? ` (+${mons.length-1})` : '') : `${mons.length} мониторов`
                                          }
                                          if (schedule.monitorID) return getMonitorName(schedule.monitorID)
                                          return '-'
                                        })()
                                      }</button>
                                    </td>
                                    <td>{getLocationName(schedule.locationID)}</td>
                    <td>{formatDateTime(schedule.startTime)}</td>
                    <td>{formatDateTime(schedule.endTime)}</td>
                    <td>{typeof schedule.priority === 'number' ? schedule.priority : '-'}</td>
                    <td>
                      <button className="btn btn-danger" onClick={() => schedule.id && handleDelete(schedule.id)}>Удалить</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

          {/* Targets details modal */}
          {targetDetails && (
            <div className="modal-overlay" onClick={closeTargetDetails}>
              <div className="modal" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                  <h3 style={{ margin: 0 }}>Детали целевых устройств</h3>
                  <button className="modal-close" onClick={closeTargetDetails}>×</button>
                </div>
                <div className="modal-body">
                  {targetDetails.location ? (
                    <div style={{ marginBottom: 8 }}><strong>Локация:</strong> {targetDetails.location.name}</div>
                  ) : null}
                  {targetDetails.groupName ? (
                    <div style={{ marginBottom: 8 }}><strong>Группа:</strong> {targetDetails.groupName}</div>
                  ) : null}
                  {targetDetails.monitors && targetDetails.monitors.length > 0 ? (
                    <div>
                      <strong>Мониторы:</strong>
                      <ul>
                        {targetDetails.monitors.map(m => (<li key={m.id ?? m.name}>{m.name ?? `ID:${m.id}`}</li>))}
                      </ul>
                    </div>
                  ) : null}
                </div>
              </div>
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
