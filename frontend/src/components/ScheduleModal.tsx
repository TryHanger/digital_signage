import { useState, useEffect } from 'react'
import { client } from '../api/client'
import type { Monitor, Location, Template } from '../api/client'

type Props = {
  visible: boolean
  onClose: () => void
  onCreated?: () => void
  monitors: Monitor[]
  locations: Location[]
  templates: Template[]
}

export default function ScheduleModal({ visible, onClose, onCreated, monitors, locations, templates }: Props) {
  const [selectedTemplateId, setSelectedTemplateId] = useState<number | null>(null)
  const [schedName, setSchedName] = useState('')
  const [schedDesc, setSchedDesc] = useState('')
  const [dateStart, setDateStart] = useState('')
  const [dateEnd, setDateEnd] = useState('')
  const [repeatPattern, setRepeatPattern] = useState<'none'|'daily'|'weekly'|'monthly'>('none')
  const [daysOfWeek, setDaysOfWeek] = useState<Record<number, boolean>>({0:false,1:false,2:false,3:false,4:false,5:false,6:false})
  const [exceptionsDates, setExceptionsDates] = useState<string[]>([])
  const [monthlyDay, setMonthlyDay] = useState<number | null>(null)
  const [mode, setMode] = useState<'rotation'|'override'>('rotation')
  const [formData, setFormData] = useState<any>({
    monitorID: undefined,
    locationID: undefined,
    startTime: '',
    endTime: '',
    priority: 0,
  })
  const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null)
  const [selectedMonitors, setSelectedMonitors] = useState<number[]>([])
  const [selectionMode, setSelectionMode] = useState<'location' | 'group' | 'monitors'>('location')
  const [selectedAllLocations, setSelectedAllLocations] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!visible) return
    // reset when opened
    setSelectedTemplateId(null)
    setSchedName('')
  setSchedDesc('')
    setDateStart('')
    setDateEnd('')
    setRepeatPattern('none')
    setDaysOfWeek({0:false,1:false,2:false,3:false,4:false,5:false,6:false})
  setExceptionsDates([])
  setMonthlyDay(null)
    setMode('rotation')
    setFormData({ monitorID: undefined, locationID: undefined, startTime: '', endTime: '', priority: 0 })
    setSelectedGroupId(null)
    setSelectedMonitors([])
    setSelectedAllLocations(false)
    setError('')
  }, [visible])

  useEffect(() => {
    // template selection currently does not prefill content (content selection removed)
  }, [selectedTemplateId, templates])

  const handleSubmit = async () => {
    setError('')
    // basic validation: require name
    if (!schedName || schedName.trim().length === 0) {
      setError('Введите название расписания')
      return
    }
    setLoading(true)
    try {
      const base: any = {
        ...formData,
        name: schedName,
        description: schedDesc || undefined,
        templateID: selectedTemplateId || undefined,
        dateStart: dateStart || undefined,
        dateEnd: dateEnd || undefined,
        repeatPattern: repeatPattern === 'none' ? undefined : repeatPattern,
        daysOfWeek: repeatPattern === 'weekly' ? Object.keys(daysOfWeek).filter(k => daysOfWeek[Number(k)]).map(Number) : undefined,
        mode,
      }

      // Build final payload according to exclusive selection mode
      const payload: any = { ...base }
      if (selectionMode === 'location') {
        // only include locationID
        if (!formData.locationID) delete payload.locationID
        delete payload.groupId
        delete payload.monitors
      } else if (selectionMode === 'group') {
        payload.groupId = selectedGroupId || undefined
        delete payload.locationID
        delete payload.monitors
      } else if (selectionMode === 'monitors') {
        payload.monitors = selectedMonitors.map(id => ({ id }))
        delete payload.locationID
        delete payload.groupId
      }

      // Include repetition fields
      if (repeatPattern === 'none') {
        // single day: dateStart must be set (if provided)
        if (dateStart) payload.dateStart = dateStart
        delete payload.daysOfWeek
      } else if (repeatPattern === 'daily') {
        payload.repeatPattern = 'daily'
        if (dateStart) payload.dateStart = dateStart
        if (dateEnd) payload.dateEnd = dateEnd || undefined
        if (exceptionsDates && exceptionsDates.length > 0) payload.exceptions = exceptionsDates.map(d => ({ date: d }))
      } else if (repeatPattern === 'weekly') {
        payload.repeatPattern = 'weekly'
        // convert daysOfWeek map (0..6) to backend weekdays (1..7) if true
        const days = Object.keys(daysOfWeek).filter(k => daysOfWeek[Number(k)]).map(k => Number(k) + 1)
        if (days.length > 0) payload.daysOfWeek = days
        if (dateStart) payload.dateStart = dateStart
        if (dateEnd) payload.dateEnd = dateEnd || undefined
      } else if (repeatPattern === 'monthly') {
        payload.repeatPattern = 'monthly'
        if (monthlyDay) payload.monthDay = monthlyDay
        if (dateStart) payload.dateStart = dateStart
        if (dateEnd) payload.dateEnd = dateEnd || undefined
      }

      const result = await client.schedules.create(payload)
      // simple handling: if result contains error/conflicts return message
      if (result && typeof result === 'object' && 'error' in result && result.error) {
        setError(result.error || 'Ошибка сервера')
        setLoading(false)
        return
      }

      if (onCreated) onCreated()
      onClose()
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Неизвестная ошибка')
    } finally {
      setLoading(false)
    }
  }

  if (!visible) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxHeight: '80vh', overflow: 'auto' }}>
        <div className="modal-header">
          <h3 style={{ margin: 0 }}>Создать расписание</h3>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          {error && <div className="alert alert-error">{error}</div>}
          <div className="form-group">
            <label>Шаблон (опционально)</label>
            <select value={selectedTemplateId ?? ''} onChange={(e) => setSelectedTemplateId(e.target.value ? Number(e.target.value) : null)}>
              <option value="">-- без шаблона --</option>
              {templates.map((t) => (<option key={t.id} value={t.id}>{t.name}</option>))}
            </select>
          </div>

          <div className="form-group">
            <label>Название *</label>
            <input type="text" value={schedName} onChange={(e) => setSchedName(e.target.value)} placeholder="Название расписания" required />
          </div>

          <div className="form-group">
            <label>Описание</label>
            <textarea value={schedDesc} onChange={(e) => setSchedDesc(e.target.value)} placeholder="Краткое описание (необязательно)" rows={3} />
          </div>

          {/* Выбор: Локации / Группы мониторов / Определенные мониторы */}
          <div style={{ display: 'flex', gap: 12 }}>
            <div style={{ minWidth: 160, display: 'flex', flexDirection: 'column', gap: 8 }}>
              <button type="button" className={selectionMode === 'location' ? 'btn btn-secondary' : 'btn'} onClick={() => { setSelectionMode('location'); setSelectedGroupId(null); setSelectedMonitors([]); setSelectedAllLocations(false); }}>Локации</button>
              <button type="button" className={selectionMode === 'group' ? 'btn btn-secondary' : 'btn'} onClick={() => { setSelectionMode('group'); setFormData({ ...formData, locationID: undefined }); setSelectedMonitors([]); setSelectedAllLocations(false); }}>Группы мониторов</button>
              <button type="button" className={selectionMode === 'monitors' ? 'btn btn-secondary' : 'btn'} onClick={() => { setSelectionMode('monitors'); setFormData({ ...formData, locationID: undefined }); setSelectedGroupId(null); setSelectedAllLocations(false); }}>Определенные мониторы</button>
            </div>

            <div style={{ flex: 1, borderLeft: '1px solid #eee', paddingLeft: 12 }}>
              {selectionMode === 'location' && (
                <div>
                  <label>Выберите локацию</label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 6, marginTop: 6 }}>
                    <label style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <input type="radio" name="locationChoice" checked={selectedAllLocations} onChange={() => { setSelectedAllLocations(true); setFormData({ ...formData, locationID: undefined }) }} /> Все локации
                    </label>
                    {locations.map(l => (
                      <label key={l.id} style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                        <input type="radio" name="locationChoice" value={l.id} checked={formData.locationID === l.id} onChange={() => { setSelectedAllLocations(false); setFormData({ ...formData, locationID: l.id }) }} />
                        <span>{l.name}</span>
                      </label>
                    ))}
                  </div>
                </div>
              )}

              {selectionMode === 'group' && (
                <div>
                  <label>Выберите группу мониторов</label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 6, marginTop: 6 }}>
                    <label style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <input type="radio" name="groupChoice" checked={!selectedGroupId} onChange={() => { setSelectedGroupId(null); setSelectedMonitors([]); }} /> Нет группы
                    </label>
                    {/* build unique groups from monitors */}
                    {Array.from(new Map(monitors.filter(m => m.groupID).map(m => [m.groupID, m.group?.name || `Группа ${m.groupID}`]))).map(([gid, gname]) => (
                      <label key={String(gid)} style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                        <input type="radio" name="groupChoice" value={String(gid)} checked={selectedGroupId === Number(gid)} onChange={() => {
                          const idn = Number(gid)
                          setSelectedGroupId(idn)
                          const ids = monitors.filter(m => m.groupID === idn).map(m => m.id!).filter(Boolean) as number[]
                          setSelectedMonitors(ids)
                        }} />
                        <span>{String(gname)}</span>
                      </label>
                    ))}
                  </div>
                </div>
              )}

              {selectionMode === 'monitors' && (
                <div>
                  <label>Выберите мониторы (чекбоксы)</label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 6, marginTop: 6, maxHeight: 220, overflow: 'auto' }}>
                    <label style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <input type="checkbox" checked={selectedMonitors.length === monitors.length} onChange={(e) => {
                        if (e.target.checked) setSelectedMonitors(monitors.map(m => m.id!).filter(Boolean) as number[])
                        else setSelectedMonitors([])
                      }} /> Выбрать все
                    </label>
                    {monitors.map(m => (
                      <label key={m.id} style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                        <input type="checkbox" value={m.id} checked={selectedMonitors.includes(m.id!)} onChange={(e) => {
                          const id = Number(e.target.value)
                          if (e.target.checked) setSelectedMonitors(prev => Array.from(new Set([...prev, id])))
                          else setSelectedMonitors(prev => prev.filter(x => x !== id))
                        }} />
                        <span>{m.name}{m.group ? ` — ${m.group.name}` : ''}{m.locationID ? ` (${locations.find(l=>l.id===m.locationID)?.name || ''})` : ''}</span>
                      </label>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Повторение: показывается после выбора устройств/локации */}
          <div className="form-group">
            <label>Повтор</label>
            <div style={{ display: 'flex', gap: 8, marginTop: 6 }}>
              <label style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                <input type="radio" name="repeat" checked={repeatPattern === 'none'} onChange={() => { setRepeatPattern('none'); setDaysOfWeek({0:false,1:false,2:false,3:false,4:false,5:false,6:false}); setExceptionsDates([]); setMonthlyDay(null); }} /> Нет
              </label>
              <label style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                <input type="radio" name="repeat" checked={repeatPattern === 'daily'} onChange={() => { setRepeatPattern('daily'); setDaysOfWeek({0:false,1:false,2:false,3:false,4:false,5:false,6:false}); setMonthlyDay(null); }} /> Ежедневно
              </label>
              <label style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                <input type="radio" name="repeat" checked={repeatPattern === 'weekly'} onChange={() => { setRepeatPattern('weekly'); setExceptionsDates([]); setMonthlyDay(null); }} /> Еженедельно
              </label>
            </div>
          </div>

          {/* Repeat-specific UI */}
          {repeatPattern === 'none' && (
            <div className="form-group">
              <label>Дата показа</label>
              <input type="date" value={dateStart} onChange={(e) => setDateStart(e.target.value)} />
            </div>
          )}

          {repeatPattern === 'daily' && (
            <div>
              <div className="form-row">
                <div className="form-group">
                  <label>Начальная дата</label>
                  <input type="date" value={dateStart} onChange={(e) => setDateStart(e.target.value)} />
                </div>
                <div className="form-group">
                  <label>Дата окончания (опционально)</label>
                  <input type="date" value={dateEnd} onChange={(e) => setDateEnd(e.target.value)} />
                </div>
              </div>
              <div className="form-group">
                <label>Даты-исключения (не показывать)</label>
                <div style={{ display: 'flex', gap: 8, marginTop: 6, alignItems: 'center' }}>
                  <input type="date" id="exceptionInput" />
                  <button type="button" className="btn" onClick={() => {
                    const el: any = document.getElementById('exceptionInput')
                    if (!el) return
                    const v = el.value
                    if (!v) return
                    setExceptionsDates(prev => Array.from(new Set([...prev, v])))
                    el.value = ''
                  }}>Добавить</button>
                </div>
                <div style={{ marginTop: 8 }}>
                  {exceptionsDates.map((d) => (
                    <div key={d} style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                      <span>{d}</span>
                      <button type="button" className="btn btn-small" onClick={() => setExceptionsDates(prev => prev.filter(x => x !== d))}>Удалить</button>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {repeatPattern === 'weekly' && (
            <div className="form-group">
              <label>Дни недели</label>
              <div style={{ display: 'flex', gap: 8, marginTop: 6 }}>
                {[[1,'Пн'],[2,'Вт'],[3,'Ср'],[4,'Чт'],[5,'Пт'],[6,'Сб'],[7,'Вс']].map(([num, label]) => (
                  <label key={String(num)} style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                    <input type="checkbox" checked={!!daysOfWeek[(Number(num)-1)]} onChange={(e) => setDaysOfWeek(prev => ({ ...prev, [Number(num)-1]: e.target.checked }))} />
                    <span>{label}</span>
                  </label>
                ))}
              </div>
              <div style={{ marginTop: 8 }}>
                <label>Начальная дата</label>
                <input type="date" value={dateStart} onChange={(e) => setDateStart(e.target.value)} />
                <label style={{ marginLeft: 12 }}>Дата окончания (опционально)</label>
                <input type="date" value={dateEnd} onChange={(e) => setDateEnd(e.target.value)} />
              </div>
            </div>
          )}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 12 }}>
            <button className="btn btn-secondary" onClick={onClose}>Отмена</button>
            <button className="btn btn-primary" onClick={handleSubmit} disabled={loading}>{loading ? 'Сохраняем...' : 'Создать'}</button>
          </div>
        </div>
      </div>
    </div>
  )
}
