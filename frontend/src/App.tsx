import { useState } from 'react'
import ContentsPage from './pages/ContentsPage'
import LocationsPage from './pages/LocationsPage'
import MonitorsPage from './pages/MonitorsPage'
import SchedulesPage from './pages/SchedulesPage'

type Tab = 'locations' | 'monitors' | 'contents' | 'schedules'

export default function App() {
    // по умолчанию открываем "Локации" как вы просили
    const [tab, setTab] = useState<Tab>('locations')

    return (
        <div className="app">
            <aside className="sidebar">
                <h1>Digital Signage</h1>
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
                </nav>
            </aside>

            <main className="content-area">
                {tab === 'contents' && <ContentsPage />}
                {tab === 'locations' && <LocationsPage />}
                {tab === 'monitors' && <MonitorsPage />}
                {tab === 'schedules' && <SchedulesPage />}
            </main>
        </div>
    )
}