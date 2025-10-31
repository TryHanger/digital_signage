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
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3 style={{ margin: 0 }}>{title}</h3>
          <button className="modal-close" onClick={onClose} aria-label="Закрыть">×</button>
        </div>
        <div className="modal-body">{children}</div>
      </div>
    </div>
  )
}
