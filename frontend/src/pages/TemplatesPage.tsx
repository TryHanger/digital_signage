import { useEffect, useState } from 'react'
import { client } from '../api/client'
import type { Template, Content, TemplateBlock } from '../api/client'
import Modal from '../components/Modal'

export default function TemplatesPage() {
  const [templates, setTemplates] = useState<Template[]>([])
  const [contents, setContents] = useState<Content[]>([])
  
  const [loading, setLoading] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [form, setForm] = useState<Partial<Template>>({ blocks: [] })
  const [editingTemplate, setEditingTemplate] = useState<Template | null>(null)
  const [showModal, setShowModal] = useState(false)
  const [dragInfo, setDragInfo] = useState<{ blockIdx: number, contentIdx: number } | null>(null)
  const [dragOverBlock, setDragOverBlock] = useState<number | null>(null)
  const [openContentDropdown, setOpenContentDropdown] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const load = async () => {
    setLoading(true)
    try {
      const [tpls, conts] = await Promise.all([
        client.templates.getAll(),
        client.contents.getAll(),
      ])
      setTemplates(tpls)
      setContents(conts)
      setError('')
    } catch (err: any) {
      setError('Ошибка загрузки: ' + (err.message || err))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  // close content dropdown when clicking outside
  useEffect(() => {
    const onDocClick = (e: MouseEvent) => {
      const el = e.target as HTMLElement | null
      if (!el) return
      if (!el.closest('[data-dropdown-block]')) setOpenContentDropdown(null)
    }
    document.addEventListener('click', onDocClick)
    return () => document.removeEventListener('click', onDocClick)
  }, [])

  const handleCreate = async (e: any) => {
    e.preventDefault()
    setError('')
    try {
      if (!name) return setError('Укажите имя шаблона')
      // validate blocks times
      const blocksRaw: any[] = (form.blocks || []).map((b: any) => ({ ...b }))
      // ensure blocks do not overlap (HH:MM ranges)
      const toMinutes = (hhmm: string) => {
        const parts = (hhmm || '00:00').split(':').map((x: string) => Number(x))
        return (parts[0] || 0) * 60 + (parts[1] || 0)
      }
      for (let i = 0; i < blocksRaw.length; i++) {
        const aStart = toMinutes(blocksRaw[i].startTime)
        const aEnd = toMinutes(blocksRaw[i].endTime)
        if (aStart >= aEnd) return setError(`У блока "${blocksRaw[i].name}" время начала должно быть раньше конца`)
        for (let j = i + 1; j < blocksRaw.length; j++) {
          const bStart = toMinutes(blocksRaw[j].startTime)
          const bEnd = toMinutes(blocksRaw[j].endTime)
          // overlap if ranges intersect
          if (!(aEnd <= bStart || bEnd <= aStart)) {
            return setError(`Блоки "${blocksRaw[i].name}" и "${blocksRaw[j].name}" пересекаются по времени`)
          }
        }
      }
      // normalize blocks: convert HH:MM into RFC3339 strings and normalize contents keys to backend shape
  const blocks: any[] = blocksRaw.map((b) => {
        const startHH = b.startTime || '00:00'
        const endHH = b.endTime || '00:00'
        // build local date-time using today's date
        const today = new Date()
        const YYYY = today.getFullYear()
        const MM = String(today.getMonth() + 1).padStart(2, '0')
        const DD = String(today.getDate()).padStart(2, '0')
        const startLocal = `${YYYY}-${MM}-${DD}T${startHH}`
        const endLocal = `${YYYY}-${MM}-${DD}T${endHH}`
        const startTime = (client as any)._time.toRFC3339WithTZ(startLocal)
        const endTime = (client as any)._time.toRFC3339WithTZ(endLocal)

        const contents = (b.contents || []).map((ct: any, ci: number) => ({
          id: Number(ct.id || 0),
          contentId: Number(ct.contentID || ct.contentId),
          duration: Number(ct.duration || 0),
          order: ci + 1,
          type: 'content',
        }))

        return {
          id: Number(b.id || 0),
          templateId: Number(editingTemplate?.id || 0),
          name: b.name,
          startTime,
          endTime,
          contents,
        }
      })
      for (const b of blocks) {
        if (!b.name) return setError('У блока должно быть имя')
        const st = b.startTime
        const et = b.endTime
        if (!st || !et) return setError('Укажите время начала/конца для каждого блока')
        // compare by Date parsing
        try {
          const ds = new Date(st)
          const de = new Date(et)
          if (isNaN(ds.getTime()) || isNaN(de.getTime())) return setError(`Некорректное время у блока "${b.name}"`)
          if (ds.getTime() >= de.getTime()) return setError(`У блока "${b.name}" время начала должно быть раньше конца`)
        } catch (err) {
          return setError(`Некорректное время у блока "${b.name}"`)
        }
        // contents presence
        if (!Array.isArray(b.contents) || b.contents.length === 0) return setError(`У блока "${b.name}" должен быть как минимум один контент`)
      }

      const payload: any = {
        name,
        description,
        blocks,
      }
      if (editingTemplate && editingTemplate.id) {
        await client.templates.update(editingTemplate.id as number, payload)
        setSuccess('Шаблон обновлён')
      } else {
        await client.templates.create(payload)
        setSuccess('Шаблон создан')
      }
      setName('')
      setDescription('')
      setForm({ blocks: [] })
      setEditingTemplate(null)
      setShowModal(false)
      await load()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err: any) {
      setError('Ошибка создания: ' + (err.message || err))
    }
  }

  const startEdit = async (tpl: Template) => {
    // convert backend template blocks (RFC3339 start/end) to UI-friendly HH:MM
    const uiBlocks = (tpl.blocks || []).map((b: any) => {
      const s = b.startTime ? new Date(b.startTime) : null
      const e = b.endTime ? new Date(b.endTime) : null
      const fmt = (d: Date | null) => d ? `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}` : '00:00'
      return {
        id: b.id,
        name: b.name,
        startTime: fmt(s),
        endTime: fmt(e),
  contents: (b.contents || []).map((ct: any) => ({ id: ct.id, contentID: ct.contentId, duration: ct.duration })),
      }
    })
    setEditingTemplate(tpl)
    setName(tpl.name || '')
    setDescription(tpl.description || '')
    setForm({ blocks: uiBlocks })
    setShowModal(true)
  }

  const addBlock = () => {
    setForm(prev => ({ blocks: [...(prev.blocks || []), { name: 'Новый блок', startTime: '08:00', endTime: '09:00', contents: [] }] }))
  }

  const removeBlock = (idx: number) => {
    setForm(prev => ({ blocks: (prev.blocks || []).filter((_, i) => i !== idx) }))
  }

  const updateBlock = (idx: number, patch: Partial<TemplateBlock>) => {
    setForm(prev => ({ blocks: (prev.blocks || []).map((b: any, i: number) => i === idx ? { ...b, ...patch } : b) }))
  }

  const addContentToBlock = (blockIdx: number, contentID: number) => {
    setForm(prev => ({
      blocks: (prev.blocks || []).map((b: any, i: number) => i === blockIdx ? { ...b, contents: [...(b.contents || []), { contentID, duration: 10 }] } : b)
    }))
  }

  const moveContent = (fromBlock: number, fromIdx: number, toBlock: number, toIdx: number) => {
    setForm(prev => {
      const blocks = (prev.blocks || []).map((b: any) => ({ ...b, contents: [...(b.contents || [])] }))
      if (!blocks[fromBlock] || !blocks[fromBlock].contents) return prev
      const item = blocks[fromBlock].contents.splice(fromIdx, 1)[0]
      if (!item) return { blocks }
      // if dropping into same block and index after removal needs adjustment
      const insertIdx = (toBlock === fromBlock && toIdx > fromIdx) ? toIdx - 1 : toIdx
      if (!blocks[toBlock]) {
        // append to last block if target missing
        blocks[blocks.length - 1].contents.push(item)
      } else {
        if (insertIdx === undefined || insertIdx === null) {
          blocks[toBlock].contents.push(item)
        } else {
          blocks[toBlock].contents.splice(Math.max(0, insertIdx), 0, item)
        }
      }
      return { blocks }
    })
  }

  const handleDragStart = (e: React.DragEvent, blockIdx: number, contentIdx: number) => {
    setDragInfo({ blockIdx, contentIdx })
    try { e.dataTransfer.setData('text/plain', '') } catch (err) { /* ignore */ }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDragEnterBlock = (e: React.DragEvent, blockIdx: number) => {
    e.preventDefault()
    setDragOverBlock(blockIdx)
  }

  const handleDragLeaveBlock = (blockIdx: number) => {
    // ignore if entering a child; clear only when leaving the block container
    setDragOverBlock((cur) => (cur === blockIdx ? null : cur))
  }

  const handleDropOnItem = (e: React.DragEvent, blockIdx: number, contentIdx: number) => {
    e.preventDefault()
    if (!dragInfo) return
    moveContent(dragInfo.blockIdx, dragInfo.contentIdx, blockIdx, contentIdx)
    setDragInfo(null)
  }

  const handleDropOnListEnd = (e: React.DragEvent, blockIdx: number) => {
    e.preventDefault()
    if (!dragInfo) return
    const len = ((form.blocks || [])[blockIdx]?.contents || []).length
    moveContent(dragInfo.blockIdx, dragInfo.contentIdx, blockIdx, len)
    setDragInfo(null)
  }

  const removeContentFromBlock = (blockIdx: number, contentIdx: number) => {
    setForm(prev => ({
  blocks: (prev.blocks || []).map((b: any, i: number) => i === blockIdx ? { ...b, contents: (b.contents || []).filter((_cc: any, j: number) => j !== contentIdx) } : b)
    }))
  }

  const handleDelete = async (id?: number) => {
    if (!id) return
    if (!confirm('Удалить шаблон?')) return
    try {
      await client.templates.delete(id)
      await load()
      setSuccess('Удалено')
      setTimeout(() => setSuccess(''), 2000)
    } catch (err: any) {
      setError('Ошибка удаления: ' + (err.message || err))
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
        <h1 className="page-title" style={{ margin: 0 }}>Шаблоны</h1>
        <div style={{ marginLeft: 'auto', display: 'flex', gap: 8 }}>
          <button className="btn btn-primary" onClick={() => { setEditingTemplate(null); setName(''); setDescription(''); setForm({ blocks: [] }); setShowModal(true) }}>+ Создать шаблон</button>
          <button className="btn btn-secondary" onClick={load} disabled={loading}>{loading ? '...' : 'Обновить'}</button>
        </div>
      </div>
      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}
      <Modal title={editingTemplate ? 'Редактировать шаблон' : 'Создать шаблон'} visible={showModal} onClose={() => setShowModal(false)}>
        <form onSubmit={handleCreate} onClick={(e) => { const el = e.target as HTMLElement; if (!el.closest('[data-dropdown-block]')) setOpenContentDropdown(null) }}>
          <div className="form-group">
            <label>Имя шаблона</label>
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="form-group">
            <label>Описание</label>
            <textarea value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>

          <div style={{ marginTop: 12 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h3>Блоки</h3>
              <div>
                <button type="button" className="btn btn-secondary" onClick={addBlock}>+ Добавить блок</button>
              </div>
            </div>

            {(form.blocks || []).length === 0 && <div className="empty-state"><p>Блоков пока нет. Добавьте первый блок.</p></div>}

            {(form.blocks || []).map((b: any, idx: number) => (
              <div key={idx} style={{ border: '1px solid #ddd', padding: 8, marginTop: 8 }}>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <input value={b.name} onChange={(e) => updateBlock(idx, { name: e.target.value })} style={{ flex: 1 }} />
                  <button type="button" className="btn btn-danger" onClick={() => removeBlock(idx)}>Удалить блок</button>
                </div>
                <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
                  <div className="form-group" style={{ flex: 1 }}>
                    <label>Начало (HH:MM)</label>
                    <input value={b.startTime} onChange={(e) => updateBlock(idx, { startTime: e.target.value })} placeholder="08:00" />
                  </div>
                  <div className="form-group" style={{ flex: 1 }}>
                    <label>Конец (HH:MM)</label>
                    <input value={b.endTime} onChange={(e) => updateBlock(idx, { endTime: e.target.value })} placeholder="10:00" />
                  </div>
                </div>
                <div style={{ marginTop: 8 }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8 }}>
                    <label style={{ margin: 0 }}>Контент в блоке</label>
                    <div style={{ position: 'relative', display: 'inline-block' }} data-dropdown-block>
                      <button type="button" className="btn btn-secondary" onClick={(e) => { e.stopPropagation(); setOpenContentDropdown(openContentDropdown === idx ? null : idx) }}>+ Добавить контент ▾</button>
                      {openContentDropdown === idx && (
                        <div style={{ position: 'absolute', top: 'calc(100% + 6px)', right: 0, zIndex: 30, background: '#fff', border: '1px solid #ddd', boxShadow: '0 6px 18px rgba(0,0,0,0.08)', borderRadius: 6, minWidth: 220, maxHeight: 260, overflow: 'auto' }}>
                          {contents.length === 0 && <div style={{ padding: 10, color: '#666' }}>Нет доступного контента</div>}
                          {contents.map(c => (
                            <div key={c.id} style={{ padding: 8, cursor: 'pointer', borderBottom: '1px solid #f1f1f1' }} onClick={(e) => { e.stopPropagation(); addContentToBlock(idx as number, Number(c.id)); setOpenContentDropdown(null) }}>{c.title}</div>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                  <ul onDragOver={handleDragOver} onDrop={(e) => handleDropOnListEnd(e, idx)} onDragEnter={(e) => handleDragEnterBlock(e, idx)} onDragLeave={() => handleDragLeaveBlock(idx)} style={{ padding: 8, borderRadius: 4, background: dragOverBlock === idx ? '#f7fbff' : '#fbfdff', border: dragOverBlock === idx ? '2px dashed #66a3ff' : '1px solid #eee' }}>
                    {(b.contents || []).map((ct: any, ci: number) => (
                      <li key={ci} draggable={true} onDragStart={(e) => handleDragStart(e as any, idx, ci)} onDragOver={handleDragOver as any} onDrop={(e) => handleDropOnItem(e as any, idx, ci)} style={{ display: 'flex', gap: 8, alignItems: 'center', cursor: 'grab', padding: 6, borderRadius: 4, background: (dragInfo && dragInfo.blockIdx === idx && dragInfo.contentIdx === ci) ? '#e8f0ff' : undefined, boxShadow: (dragInfo && dragInfo.blockIdx === idx && dragInfo.contentIdx === ci) ? '0 2px 8px rgba(0,0,0,0.08)' : undefined, transform: (dragInfo && dragInfo.blockIdx === idx && dragInfo.contentIdx === ci) ? 'scale(1.01)' : undefined }}>
                        <span style={{ flex: 1 }}>{contents.find(x => x.id === ct.contentID)?.title ?? `ID:${ct.contentID}`}</span>
                        <input type="number" value={ct.duration || 10} onChange={(e) => updateBlock(idx, { contents: (b.contents || []).map((cc: any, j: number) => j === ci ? { ...cc, duration: Number(e.target.value) } : cc) })} style={{ width: 80 }} />
                        <button type="button" className="btn btn-sm btn-danger" onClick={() => removeContentFromBlock(idx, ci)}>Удалить</button>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            ))}
          </div>

          <div className="toolbar" style={{ marginTop: 12 }}>
            <button className="btn btn-primary" type="submit">{editingTemplate ? 'Сохранить' : 'Сохранить шаблон'}</button>
          </div>
        </form>
      </Modal>

  <div className="card" style={{ marginTop: 16 }}>
        <div className="toolbar">
        </div>
        {templates.length === 0 ? (
          <div className="empty-state"><p>Шаблонов пока нет.</p></div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr><th>Имя</th><th>Описание</th><th>Блоки</th><th>Действия</th></tr>
              </thead>
              <tbody>
                {templates.map(t => (
                  <tr key={t.id}>
                      <td>{t.name}</td>
                      <td style={{ maxWidth: 300 }}>{t.description ?? '-'}</td>
                      <td>{(t.blocks || []).length}</td>
                      <td>
                        <button className="btn" onClick={() => startEdit(t)}>Редактировать</button>
                        <button className="btn btn-danger" onClick={() => handleDelete(t.id)}>Удалить</button>
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
