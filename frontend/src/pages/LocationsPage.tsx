import type { FormEvent } from 'react'
import { useEffect, useState } from 'react'
import { client } from '../api/client'
import type { Location } from '../api/client'
import Modal from '../components/Modal'

export default function LocationsPage() {
  const [locations, setLocations] = useState<Location[]>([])
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Location>({ name: '' })
  const [editingId, setEditingId] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [showModal, setShowModal] = useState(false)

  const loadLocations = async () => {
    setLoading(true)
    try {
      const data = await client.locations.getAll()
      setLocations(data)
      setError('')
    } catch (err: any) {
      setError('Ошибка загрузки локаций: ' + (err.message || 'Неизвестная ошибка'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadLocations()
  }, [])

  const handleOpenNew = () => {
    setFormData({ name: '' })
    setEditingId(null)
    setError('')
    setSuccess('')
    setShowModal(true)
  }

  const handleEdit = (loc: Location) => {
    setFormData({ name: loc.name })
    setEditingId(loc.id || null)
    setError('')
    setSuccess('')
    setShowModal(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить эту локацию?')) return
    try {
      await client.locations.delete(id)
      setSuccess('Локация удалена')
      await loadLocations()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка удаления: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleCancel = () => {
    setFormData({ name: '' })
    setEditingId(null)
    setError('')
    setSuccess('')
    setShowModal(false)
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    try {
      if (editingId) {
        await client.locations.update(editingId, formData)
        setSuccess('Локация успешно обновлена')
      } else {
        await client.locations.create(formData)
        setSuccess('Локация успешно создана')
      }
      setFormData({ name: '' })
      setEditingId(null)
      await loadLocations()
      setShowModal(false)
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка: ' + (err.response?.data?.error || err.message || 'Неизвестная ошибка'))
    }
  }

  return (
    <div>
      <h1 className="page-title">Локации</h1>

      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}

      <div className="center-column">
        <div className="create-bar card small-card">
          <h3 style={{ margin: 0 }}>Локации</h3>
          <div style={{ marginLeft: 'auto' }}>
            <button className="btn btn-primary" onClick={handleOpenNew}>Создать локацию</button>
          </div>
        </div>

        <Modal visible={showModal} onClose={() => setShowModal(false)} title={editingId ? 'Редактировать локацию' : 'Создать локацию'}>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Название локации *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Например: Офис на Тверской"
                required
              />
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
              <button type="button" className="btn btn-secondary" onClick={handleCancel}>Отмена</button>
              <button type="submit" className="btn btn-primary">{editingId ? 'Сохранить' : 'Создать'}</button>
            </div>
          </form>
        </Modal>

        <div className="card list-card">
          <div className="toolbar">
            <h2 style={{ margin: 0, flex: 1 }}>Список локаций</h2>
            <button className="btn btn-secondary" onClick={loadLocations} disabled={loading}>
              {loading ? 'Загрузка...' : 'Обновить'}
            </button>
          </div>

          {locations.length === 0 ? (
            <div className="empty-state">
              <p>Локаций пока нет. Создайте первую локацию выше.</p>
            </div>
          ) : (
            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Название</th>
                    <th>Действия</th>
                  </tr>
                </thead>
                <tbody>
                  {locations.map((loc, i) => (
                    <tr key={loc.id ?? `loc-${i}`}>
                      <td>{loc.id}</td>
                      <td>{loc.name}</td>
                      <td>
                        <button
                          className="btn btn-secondary"
                          style={{ marginRight: '8px' }}
                          onClick={() => handleEdit(loc)}
                        >
                          Редактировать
                        </button>
                        <button className="btn btn-danger" onClick={() => loc.id && handleDelete(loc.id)}>
                          Удалить
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
