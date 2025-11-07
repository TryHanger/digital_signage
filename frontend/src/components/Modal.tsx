import React from 'react'

type Props = {
  title?: string
  visible: boolean
  onClose: () => void
  children?: React.ReactNode
}

export default function Modal({ title, visible, onClose, children }: Props) {
  if (!visible) return null

  return (
    <div className="modal-overlay" onClick={onClose} role="dialog" aria-modal="true">
      <div className="modal" onClick={(e) => e.stopPropagation()} style={{ maxHeight: 'calc(100vh - 160px)', width: '820px', overflow: 'hidden' }}>
        <div className="modal-header" style={{ flex: '0 0 auto' }}>
          <h3 style={{ margin: 0 }}>{title}</h3>
          <button className="modal-close" onClick={onClose} aria-label="Закрыть">×</button>
        </div>
        <div className="modal-body" style={{ overflowY: 'auto', maxHeight: 'calc(100vh - 220px)', padding: 12 }}>{children}</div>
      </div>
    </div>
  )
}
