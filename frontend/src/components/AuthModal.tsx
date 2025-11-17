import React, { useState } from 'react'
import { client, setAccessToken, auth } from '../api/client'

type Mode = 'login' | 'register'

export default function AuthModal({ visible, onClose, onLogged }: { visible: boolean; onClose: () => void; onLogged?: () => void }) {
  const [mode, setMode] = useState<Mode>('login')
  // step: 1 = choose identifier (email/phone), 2 = password
  const [step, setStep] = useState<number>(1)

  const [usePhone, setUsePhone] = useState<boolean>(false)
  const [identifier, setIdentifier] = useState('') // email or phone depending on usePhone
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [rememberMe, setRememberMe] = useState(false)
  const [error, setError] = useState('')

  if (!visible) return null

  const goNext = () => {
    setError('')
    // basic validation for step 1
    if (!identifier || identifier.trim() === '') {
      setError(usePhone ? 'Введите номер телефона' : 'Введите email')
      return
    }
    setStep(2)
  }

  const goBack = () => {
    setError('')
    setStep(1)
  }

  const submit = async () => {
    setError('')
    try {
      if (mode === 'register') {
        const payload: any = { password }
        if (usePhone) payload.phone = identifier
        else payload.email = identifier
        // include rememberMe flag optionally
        if (rememberMe) payload.remember = true
        const res = await auth.register(payload)
        if (res?.access_token) setAccessToken(res.access_token)
      } else {
        const payload: any = { password }
        if (usePhone) payload.phone = identifier
        else payload.email = identifier
        if (rememberMe) payload.remember = true
        const res = await auth.login(payload)
        if (res?.access_token) setAccessToken(res.access_token)
      }

      onClose()
      onLogged && onLogged()
    } catch (err: any) {
      setError(err?.response?.data?.message || err?.response?.data?.error || err.message || 'Ошибка')
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h3 style={{ margin: 0 }}>{mode === 'login' ? 'Войти' : 'Регистрация'}</h3>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          <div style={{ display: 'flex', gap: 8, marginBottom: 12 }}>
            <button className={`btn ${mode === 'login' ? 'btn-primary' : ''}`} onClick={() => { setMode('login'); setStep(1) }}>Войти</button>
            <button className={`btn ${mode === 'register' ? 'btn-primary' : ''}`} onClick={() => { setMode('register'); setStep(1) }}>Регистрация</button>
          </div>

          {step === 1 ? (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              <div style={{ display: 'flex', gap: 8 }}>
                <button className={`btn ${!usePhone ? 'btn-primary' : ''}`} onClick={() => setUsePhone(false)}>Email</button>
                <button className={`btn ${usePhone ? 'btn-primary' : ''}`} onClick={() => setUsePhone(true)}>Телефон</button>
              </div>
              <input
                placeholder={usePhone ? 'Телефон (e.g. +7...)' : 'Email'}
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                onKeyDown={(e) => { if (e.key === 'Enter') goNext() }}
                style={{ fontSize: 18, padding: '12px' }}
              />
              {error && <div className="alert alert-error">{error}</div>}
              <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8, alignItems: 'center' }}>
                <button className="btn btn-secondary" onClick={onClose}>Отмена</button>
                <div style={{ display: 'flex', gap: 8 }}>
                  <button className="btn" onClick={() => { setIdentifier(''); setError('') }}>Очистить</button>
                  <button className="btn btn-primary" onClick={goNext}>Далее</button>
                </div>
              </div>
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div style={{ fontSize: 13, color: '#666' }}>Вход: <b>{identifier}</b></div>
                <button className="btn" onClick={() => { setStep(1); setPassword('') }}>Изменить</button>
              </div>
              <div style={{ display: 'flex', gap: 8 }}>
                <input placeholder="Пароль" type={showPassword ? 'text' : 'password'} value={password} onChange={(e) => setPassword(e.target.value)} onKeyDown={(e) => { if (e.key === 'Enter') submit() }} style={{ flex: 1, padding: '12px', fontSize: 16 }} />
                <button className="btn" onClick={() => setShowPassword(!showPassword)}>{showPassword ? 'Скрыть' : 'Показать'}</button>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <label style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                  <input type="checkbox" checked={rememberMe} onChange={(e) => setRememberMe(e.target.checked)} />
                  <span style={{ fontSize: 13 }}>Запомнить меня</span>
                </label>
                <a href="#" onClick={(e) => { e.preventDefault(); /* implement forgot flow later */ }}>Забыли пароль?</a>
              </div>
              {error && <div className="alert alert-error">{error}</div>}
              <div style={{ display: 'flex', justifyContent: 'space-between', gap: 8 }}>
                <button className="btn btn-secondary" onClick={goBack}>Назад</button>
                <div style={{ display: 'flex', gap: 8 }}>
                  <button className="btn" onClick={() => { setPassword(''); setError('') }}>Очистить</button>
                  <button className="btn btn-primary" onClick={submit}>{mode === 'login' ? 'Войти' : 'Зарегистрироваться'}</button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
