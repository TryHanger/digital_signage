import axios from 'axios'

const API_BASE =
    import.meta.env.DEV
        ? '/api/v1'  // локальная разработка с vite proxy
        : import.meta.env.VITE_API_URL + '/api/v1';

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});


export type Location = {
  id?: number
  name: string
}

export type Monitor = {
  id?: number
  name: string
  token?: string
  status?: string
  locationID?: number
  location?: Location
  groupID?: number
  group?: MonitorGroup
  createdAt?: string
}

export type MonitorGroup = {
  id?: number
  name: string
  monitors?: Monitor[]
}

export type Content = {
  id?: number
  title: string
  type: string
  path: string
  description?: string
  duration?: number
  createdAt?: string
  updatedAt?: string
}

export type ScheduleDay = {
  id?: number
  scheduleID?: number
  date: string
}

export type Schedule = {
  id?: number
  contentID: number
  monitorID?: number
  locationID?: number
  startTime: string
  endTime: string
  priority?: number
  name?: string
  description?: string
  groupId?: number
  monitors?: Array<{ id?: number }>
  createdAt?: string
  content?: Content
  monitor?: Monitor
  location?: Location
  days?: ScheduleDay[]
}

export type Template = {
  id?: number
  name: string
  description?: string
  createdAt?: string
  blocks?: TemplateBlock[]
}

export type TemplateBlock = {
  id?: number
  name: string
  startTime: string // "HH:MM"
  endTime: string   // "HH:MM"
  contents: TemplateBlockContent[]
}

export type TemplateBlockContent = {
  id?: number
  contentID: number
  duration?: number // seconds
}

export const client = {
  // helper: convert local datetime string to RFC3339 with timezone offset
  _time: {
    toRFC3339WithTZ: (localDT?: string) => {
      if (!localDT) return ''
      // localDT expected like 'YYYY-MM-DDTHH:mm' or 'YYYY-MM-DDTHH:mm:ss'
      const d = new Date(localDT)
      if (isNaN(d.getTime())) return localDT

      const YYYY = d.getFullYear()
      const MM = String(d.getMonth() + 1).padStart(2, '0')
      const DD = String(d.getDate()).padStart(2, '0')
      const hh = String(d.getHours()).padStart(2, '0')
      const mm = String(d.getMinutes()).padStart(2, '0')
      const ss = String(d.getSeconds()).padStart(2, '0')

      const offsetMin = -d.getTimezoneOffset() // minutes ahead of UTC
      const sign = offsetMin >= 0 ? '+' : '-'
      const absOffset = Math.abs(offsetMin)
      const offH = String(Math.floor(absOffset / 60)).padStart(2, '0')
      const offM = String(absOffset % 60).padStart(2, '0')

      return `${YYYY}-${MM}-${DD}T${hh}:${mm}:${ss}${sign}${offH}:${offM}`
    },
    // helper to format schedule object(s)
    formatSchedule: (s: Schedule) => ({
      ...s,
      startTime: (client as any)._time.toRFC3339WithTZ(s.startTime as any),
      endTime: (client as any)._time.toRFC3339WithTZ(s.endTime as any),
    }),
    formatSchedules: (arr: Schedule[]) => arr.map((s) => (client as any)._time.formatSchedule(s)),
  },
  // Locations
  locations: {
    getAll: () => api.get<Location[]>('/locations').then((r: any) => {
      const payload = r?.data
      if (Array.isArray(payload)) return payload
      if (payload && Array.isArray(payload.data)) return payload.data
      console.warn('/locations returned unexpected shape, normalizing to []', payload)
      return [] as Location[]
    }),
    getById: (id: number) => api.get<Location>(`/locations/${id}`).then((r: any) => r.data),
    create: (data: Location) => api.post<Location>('/locations', data).then((r: any) => r.data),
    update: (id: number, data: Location) => api.put<Location>(`/locations/${id}`, data).then((r: any) => r.data),
    delete: (id: number) => api.delete(`/locations/${id}`),
  },

  // Monitors
  monitors: {
    getAll: () => api.get<Monitor[]>('/monitors').then((r: any) => {
      const payload = r?.data
      if (Array.isArray(payload)) return payload
      if (payload && Array.isArray(payload.data)) return payload.data
      console.warn('/monitors returned unexpected shape, normalizing to []', payload)
      return [] as Monitor[]
    }),
    create: (data: Omit<Monitor, 'id' | 'token' | 'createdAt'>) => api.post<Monitor>('/monitors', data).then((r: any) => r.data),
  },

  // Contents
  contents: {
    getAll: () => api.get<Content[]>('/contents').then((r: any) => {
      const payload = r?.data
      if (Array.isArray(payload)) return payload
      if (payload && Array.isArray(payload.data)) return payload.data
      console.warn('/contents returned unexpected shape, normalizing to []', payload)
      return [] as Content[]
    }),
    getById: (id: number) => api.get<Content>(`/contents/${id}`).then((r: any) => r.data),
    create: (data: Content) => api.post<Content>('/contents', data).then((r: any) => r.data),
    update: (id: number, data: Content) => api.put<Content>(`/contents/${id}`, data).then((r: any) => r.data),
    delete: (id: number) => api.delete(`/contents/${id}`),
  },

  // Schedules
  schedules: {
    getAll: () => api.get<Schedule[]>('/schedules').then((r: any) => {
      const payload = r?.data
      if (Array.isArray(payload)) return payload
      if (payload && Array.isArray(payload.data)) return payload.data
      console.warn('/schedules returned unexpected shape, normalizing to []', payload)
      return [] as Schedule[]
    }),
    getById: (id: number) => api.get<Schedule>(`/schedules/${id}`).then((r: any) => r.data),
    create: (data: Schedule) => {
      const payload = (client as any)._time.formatSchedule(data)
      return api.post<Schedule | { error: string; conflicts: Schedule[] }>('/schedules', payload).then((r: any) => r.data)
    },
    delete: (id: number) => api.delete(`/schedules/${id}`),
    resolveConflicts: (schedules: Schedule[]) => {
      const payload = { schedules: (client as any)._time.formatSchedules(schedules) }
      return api.put('/schedules/resolve', payload).then((r: any) => r.data)
    },
    updateSchedules: (schedules: Schedule[]) => {
      const payload = (client as any)._time.formatSchedules(schedules)
      return api.put('/schedules/update', payload).then((r: any) => r.data)
    },
  },
  // Templates
  templates: {
    getAll: () => api.get<Template[]>('/templates').then((r: any) => {
      const payload = r?.data
      if (Array.isArray(payload)) return payload
      if (payload && Array.isArray(payload.data)) return payload.data
      console.warn('/templates returned unexpected shape, normalizing to []', payload)
      return [] as Template[]
    }),
    create: (data: Template) => api.post<Template>('/templates', data).then((r: any) => r.data),
  update: (id: number, data: any) => api.put<Template>(`/templates/${id}`, data).then((r: any) => r.data),
  delete: (id: number) => api.delete(`/templates/${id}`),
  getById: (id: number) => api.get<Template>(`/templates/${id}`).then((r: any) => r.data),
  },
}

export default client
