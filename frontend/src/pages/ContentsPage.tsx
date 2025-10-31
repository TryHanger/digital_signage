import type {FormEvent} from 'react';
import { useEffect, useState } from 'react';
import { client} from '../api/client';
import type { Content } from '../api/client'

export default function ContentsPage() {
  const [contents, setContents] = useState<Content[]>([])
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<Content>({
    title: '',
    type: 'image',
    path: '',
    description: '',
    duration: 10,
  })
  const [editingId, setEditingId] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const loadContents = async () => {
    setLoading(true)
    try {
      const data = await client.contents.getAll()

      // Normalize/validate response: expect an array of contents.
      if (Array.isArray(data)) {
        setContents(data)
        setError('')
      } else if (data && Array.isArray((data as any).data)) {
        // some APIs wrap payload in { data: [...] }
        setContents((data as any).data)
        setError('')
      } else {
        console.error('Unexpected /contents response shape:', data)
        setContents([])
        setError('Неверный ответ от сервера: ожидался массив контента')
      }
    } catch (err: any) {
      setError('Ошибка загрузки контента: ' + (err.message || 'Неизвестная ошибка'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadContents()
  }, [])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    
    try {
      if (editingId) {
        await client.contents.update(editingId, formData)
        setSuccess('Контент успешно обновлен')
      } else {
        await client.contents.create(formData)
        setSuccess('Контент успешно создан')
      }
      setFormData({ title: '', type: 'image', path: '', description: '', duration: 10 })
      setEditingId(null)
      loadContents()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка: ' + (err.response?.data?.error || err.message || 'Неизвестная ошибка'))
    }
  }

  const handleEdit = (content: Content) => {
    setFormData({
      title: content.title,
      type: content.type,
      path: content.path,
      description: content.description || '',
      duration: content.duration || 10,
    })
    setEditingId(content.id || null)
    setError('')
    setSuccess('')
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить этот контент?')) return
    
    try {
      await client.contents.delete(id)
      setSuccess('Контент удален')
      loadContents()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка удаления: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleCancel = () => {
    setFormData({ title: '', type: 'image', path: '', description: '', duration: 10 })
    setEditingId(null)
    setError('')
    setSuccess('')
  }

  return (
    <div>
      <h1 className="page-title">Контент</h1>

      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}

      <div className="card">
        <h2 style={{ marginBottom: '16px' }}>{editingId ? 'Редактировать контент' : 'Создать новый контент'}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Название *</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              placeholder="Название контента"
              required
            />
          </div>
          <div className="form-group">
            <label>Тип *</label>
            <select
              value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value })}
              required
            >
              <option value="image">Изображение</option>
              <option value="video">Видео</option>
              <option value="url">URL</option>
            </select>
          </div>
          <div className="form-group">
            <label>Путь/URL *</label>
            <input
              type="text"
              value={formData.path}
              onChange={(e) => setFormData({ ...formData, path: e.target.value })}
              placeholder="Путь к файлу или URL"
              required
            />
          </div>
          <div className="form-row">
            <div className="form-group">
              <label>Длительность (секунды)</label>
              <input
                type="number"
                value={formData.duration || 10}
                onChange={(e) => setFormData({ ...formData, duration: Number(e.target.value) })}
                min="1"
              />
            </div>
          </div>
          <div className="form-group">
            <label>Описание</label>
            <textarea
              value={formData.description || ''}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={3}
              placeholder="Описание контента (необязательно)"
            />
          </div>
          <div className="toolbar">
            <button type="submit" className="btn btn-primary">
              {editingId ? 'Сохранить изменения' : 'Создать'}
            </button>
            {editingId && (
              <button type="button" className="btn btn-secondary" onClick={handleCancel}>
                Отмена
              </button>
            )}
          </div>
        </form>
      </div>

      <div className="card">
        <div className="toolbar">
          <h2 style={{ margin: 0, flex: 1 }}>Список контента</h2>
          <button className="btn btn-secondary" onClick={loadContents} disabled={loading}>
            {loading ? 'Загрузка...' : 'Обновить'}
          </button>
        </div>

        {contents.length === 0 ? (
          <div className="empty-state">
            <p>Контента пока нет. Создайте первый контент выше.</p>
          </div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Название</th>
                  <th>Тип</th>
                  <th>Путь</th>
                  <th>Длительность</th>
                  <th>Описание</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {contents.map((content, i) => (
                  <tr key={content.id ?? `content-${i}`}>
                    <td>{content.id}</td>
                    <td>{content.title}</td>
                    <td>
                      <span className="badge badge-primary" style={{ background: '#3498db' }}>
                        {content.type}
                      </span>
                    </td>
                    <td>
                      <code style={{ background: '#f0f0f0', padding: '2px 6px', borderRadius: '4px', fontSize: '12px' }}>
                        {content.path}
                      </code>
                    </td>
                    <td>{content.duration || '-'} сек</td>
                    <td>{content.description || '-'}</td>
                    <td>
                      <button
                        className="btn btn-secondary"
                        style={{ marginRight: '8px' }}
                        onClick={() => handleEdit(content)}
                      >
                        Редактировать
                      </button>
                      <button className="btn btn-danger" onClick={() => content.id && handleDelete(content.id)}>
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
  )
}
