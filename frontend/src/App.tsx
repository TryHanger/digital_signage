import { useState } from 'react'
import ContentsPage from './pages/ContentsPage'
import LocationsPage from './pages/LocationsPage'
import MonitorsPage from './pages/MonitorsPage'
import SchedulesPage from './pages/SchedulesPage'
import TemplatesPage from './pages/TemplatesPage'
import AuthModal from './components/AuthModal'
import { setAccessToken } from './api/client'

type Tab = 'locations' | 'monitors' | 'contents' | 'schedules' | 'templates'

export default function App() {
    // по умолчанию открываем "Локации" как вы просили
    const [tab, setTab] = useState<Tab>('locations')

    const [authOpen, setAuthOpen] = useState(false)
    const [userEmail, setUserEmail] = useState<string | null>(localStorage.getItem('access_token') ? 'Вошёл' : null)

    const onLogged = () => {
        setUserEmail('Вошёл')
    }

    return (
        <div className="app">
            <aside className="sidebar">
                <h1>Digital Signage</h1>
                <div style={{ marginBottom: 12 }}>
                    {!userEmail ? (
                        <>
                            <button className="btn btn-primary" style={{ marginRight: 8 }} onClick={() => setAuthOpen(true)}>Войти</button>
                            <button className="btn btn-secondary" onClick={() => { setAuthOpen(true); }}>Регистрация</button>
                        </>
                    ) : (
                        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                            <span>{userEmail}</span>
                            <button className="btn btn-danger" onClick={() => { setAccessToken(null); setUserEmail(null); }}>Выйти</button>
                        </div>
                    )}
                </div>
                <nav>
                    <button className={`nav-link ${tab === 'locations' ? 'active' : ''}`} onClick={() => setTab('locations')}>
                        Локации
                    </button>
                    <button className={`nav-link ${tab === 'monitors' ? 'active' : ''}`} onClick={() => setTab('monitors')}>
                        Мониторы
                    </button>
                    <button className={`nav-link ${tab === 'contents' ? 'active' : ''}`} onClick={() => setTab('contents')}>
                        Контент
                    </button>
                    <button className={`nav-link ${tab === 'schedules' ? 'active' : ''}`} onClick={() => setTab('schedules')}>
                        Расписание
                    </button>
                        <button className={`nav-link ${tab === 'templates' ? 'active' : ''}`} onClick={() => setTab('templates')}>
                            Шаблоны
                        </button>
                </nav>
            </aside>

            <AuthModal visible={authOpen} onClose={() => setAuthOpen(false)} onLogged={onLogged} />

            <main className="content-area">
                {tab === 'contents' && <ContentsPage />}
                {tab === 'locations' && <LocationsPage />}
                {tab === 'monitors' && <MonitorsPage />}
                    {tab === 'schedules' && <SchedulesPage />}
                    {tab === 'templates' && <TemplatesPage />}
            </main>
        </div>
    )
}