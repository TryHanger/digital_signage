import type { FormEvent } from 'react';
import { useEffect, useState } from 'react';
import { client } from '../api/client';
import type {Monitor, Location } from '../api/client'

export default function MonitorsPage() {
  const [monitors, setMonitors] = useState<Monitor[]>([])
  const [locations, setLocations] = useState<Location[]>([])
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Omit<Monitor, 'id' | 'token' | 'createdAt'>>({
    name: '',
    locationID: undefined,
  })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const loadData = async () => {
    setLoading(true)
    try {
      const [monitorsData, locationsData] = await Promise.all([
        client.monitors.getAll(),
        client.locations.getAll(),
      ])
      setMonitors(monitorsData)
      setLocations(locationsData)
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
    
    try {
      await client.monitors.create(formData)
      setSuccess('Монитор успешно создан')
      setFormData({ name: '', locationID: undefined })
      loadData()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка: ' + (err.response?.data?.error || err.message || 'Неизвестная ошибка'))
    }
  }

  const getStatusBadge = (status?: string) => {
    if (!status) return <span className="badge badge-warning">Неизвестно</span>
    if (status.toLowerCase().includes('online') || status.toLowerCase().includes('active')) {
      return <span className="badge badge-success">{status}</span>
    }
    return <span className="badge badge-danger">{status}</span>
  }

  const getLocationName = (locationID?: number) => {
    if (!locationID) return '-'
    const loc = locations.find(l => l.id === locationID)
    return loc ? loc.name : `ID: ${locationID}`
  }

  return (
    <div>
      <h1 className="page-title">Мониторы</h1>

      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}

      <div className="card">
        <h2 style={{ marginBottom: '16px' }}>Создать новый монитор</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Название монитора *</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Например: Монитор в зале ожидания"
              required
            />
          </div>
          <div className="form-group">
            <label>Локация</label>
            <select
              value={formData.locationID || ''}
              onChange={(e) => setFormData({ ...formData, locationID: e.target.value ? Number(e.target.value) : undefined })}
            >
              <option value="">Не выбрана</option>
              {locations.map((loc, i) => (
                <option key={loc.id ?? `loc-${i}`} value={loc.id}>
                  {loc.name}
                </option>
              ))}
            </select>
          </div>
          <div className="toolbar">
            <button type="submit" className="btn btn-primary">
              Создать
            </button>
          </div>
        </form>
      </div>

      <div className="card">
        <div className="toolbar">
          <h2 style={{ margin: 0, flex: 1 }}>Список мониторов</h2>
          <button className="btn btn-secondary" onClick={loadData} disabled={loading}>
            {loading ? 'Загрузка...' : 'Обновить'}
          </button>
        </div>

        {monitors.length === 0 ? (
          <div className="empty-state">
            <p>Мониторов пока нет. Создайте первый монитор выше.</p>
          </div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Название</th>
                  <th>Токен</th>
                  <th>Статус</th>
                  <th>Локация</th>
                  <th>Дата создания</th>
                </tr>
              </thead>
              <tbody>
                {monitors.map((monitor, i) => (
                  <tr key={monitor.id ?? `monitor-${i}`}>
                    <td>{monitor.id}</td>
                    <td>{monitor.name}</td>
                    <td>
                      <code style={{ background: '#f0f0f0', padding: '2px 6px', borderRadius: '4px', fontSize: '12px' }}>
                        {monitor.token || '-'}
                      </code>
                    </td>
                    <td>{getStatusBadge(monitor.status)}</td>
                    <td>{getLocationName(monitor.locationID)}</td>
                    <td>{monitor.createdAt ? new Date(monitor.createdAt).toLocaleDateString('ru-RU') : '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
