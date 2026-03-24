# KanbanMaster — Proje Planı

> Kanban tabanlı görev yönetim sistemi.
> Bireysel ve kurumsal kullanıcılar için görev takibi, takım yönetimi ve performans analizi.

---

## 1. Proje Özeti

KanbanMaster, bireylerin ve kurumların görevlerini Kanban tahtası üzerinde yönetebileceği,
takım atama, görev devretme, rapor isteme ve personel performans takibi yapılabilen
modern bir görev yönetim sistemidir.

**Hedef Kullanıcılar:**

- Takım liderleri — görev atama, devretme, rapor isteme, performans izleme
- Takım üyeleri — görev takibi, durum güncelleme, raporlama
- Kurum yöneticileri — genel bakış, takım bazlı analiz

---

## 2. Teknoloji Yığını

| # | Katman | Teknoloji | Versiyon | Amaç |
|---|--------|-----------|----------|------|
| 1 | Frontend Framework | Vue.js 3 | 3.x | SPA, Composition API |
| 2 | UI Kütüphanesi | Vuetify | 3.x | Material Design bileşenleri |
| 3 | CSS Framework | Tailwind CSS | 3.x | Utility-first stil yönetimi |
| 4 | State Management | Pinia | 2.x | Vue store yönetimi |
| 5 | Router | Vue Router | 4.x | Sayfa yönlendirme |
| 6 | Grafik Kütüphanesi | Chart.js + vue-chartjs | — | Performans grafikleri |
| 7 | Sürükle-Bırak | vuedraggable (SortableJS) | — | Kanban drag & drop |
| 8 | HTTP İstemci | Axios | — | API iletişimi |
| 9 | Backend | Go (Golang) | 1.21+ | REST API |
| 10 | Database | PostgreSQL | 15+ | İlişkisel veritabanı |
| 11 | Auth | JWT | — | Token bazlı kimlik doğrulama |
| 12 | Real-time | WebSocket (gorilla/websocket) | — | Anlık güncellemeler |
| 13 | Container | Docker + Docker Compose | — | Geliştirme ve dağıtım ortamı |
| 14 | Test (Frontend) | Vitest + Vue Test Utils | — | Unit ve component testleri |
| 15 | Test (Backend) | Go testing + httptest | — | API ve unit testleri |
| 16 | E2E Test | Playwright | — | Uçtan uca testler |

---

## 3. Kurallar

### 3.1 YAPILANLAR (Do's)

- [x] ES Modules kullan (import/export)
- [x] TypeScript strict mode aktif
- [x] Composition API + `<script setup>` kullan
- [x] Her bileşen için karşılık gelen test dosyası yaz
- [x] AAA pattern: Arrange → Act → Assert
- [x] Temiz commit mesajları, her iş için ayrı branch
- [x] PR üzerinden merge, direkt main'e push yok
- [x] Kod değişikliğinden sonra mutlaka typecheck çalıştır
- [x] Her aşama sonunda testler çalıştırılır
- [x] Kritik değişikliklerden önce onay al
- [x] Hata yönetimi her API çağrısında yapılmalı (try/catch + kullanıcı bildirimi)
- [x] Go tarafında merkezi error handling middleware kullan
- [x] Responsive tasarım (mobil uyumlu)
- [x] Component bazlı mimari — tekrar kullanılabilir, küçük parçalar
- [x] Veritabanı işlemlerinde transaction kullan

### 3.2 YAPILMAYACAKLAR (Don'ts)

- ❌ **`any` tipi kesinlikle YASAK** — TypeScript'te hiçbir yerde kullanılmayacak
- ❌ **CommonJS (require/module.exports) YASAK** — sadece ES modules
- ❌ **Tek dosyada 300 satır limiti** — zorunlu hallerde maksimum 500 satır
- ❌ **`/core/auth/` dosyaları onaysız düzenlenmeyecek**
- ❌ **İnline CSS yazılmayacak** — Tailwind utility class veya Vuetify props kullan
- ❌ **Console.log production'da kalmayacak** — geliştirme sonrası temizlenecek
- ❌ **Direkt DOM manipülasyonu yok** — Vue reaktivitesini kullan
- ❌ **God component yok** — büyük bileşenler parçalanacak
- ❌ **Testsiz kod merge edilmeyecek**
- ❌ **Hardcoded değerler yok** — config/env dosyalarından okunacak
- ❌ **Inline style yok** — tüm stiller Tailwind veya Vuetify ile

### 3.3 KOD YAZIM PRENSİPLERİ (Kod Kalitesi)

> **Temel felsefe:** "Nasıl daha kısa, daha anlaşılır ve daha modern yazarım?" sorusunu sormadan kod yazılmaz.

#### Lifecycle & Fonksiyon Ayrımı
- **`onMounted` içinde doğrudan iş mantığı YAZILMAZ** — sadece fonksiyon çağrılır
- Her fonksiyon tek bir iş yapar, ismi o işi tam anlatır

```typescript
// ❌ YANLIŞ — onMounted içinde logic
onMounted(async () => {
  loading.value = true
  try {
    const { data } = await api.get('/tasks')
    tasks.value = data
  } finally {
    loading.value = false
  }
})

// ✅ DOĞRU — ayrı fonksiyon, onMounted sadece çağırır
async function loadTasks() {
  loading.value = true
  try {
    tasks.value = await taskService.list()
  } finally {
    loading.value = false
  }
}
onMounted(loadTasks)
```

#### Kısa & Semantik Yazım
- Gereksiz değişken tanımlama yapma — doğrudan kullan
- Ternary, optional chaining, nullish coalescing tercih et
- Computed property'leri tek satırda tut (mümkünse)

```typescript
// ❌ Uzun
const userName = computed(() => {
  if (authStore.user) {
    return authStore.user.name
  }
  return ''
})

// ✅ Kısa
const userName = computed(() => authStore.user?.name ?? '')
```

#### Modern Vue 3 Kalıpları
- `<script setup>` zorunlu — Options API yasak
- `defineProps` + `defineEmits` generic syntax kullan
- Composable fonksiyonlar `use` prefix ile, reactive state döndürür
- `watch` yerine `watchEffect` tercih et (bağımlılık açıksa)
- `ref` + `.value` yerine `reactive` kullanma — tutarlılık için sadece `ref`

```typescript
// ❌ Gereksiz reactive
const state = reactive({ loading: false, items: [] })

// ✅ Ayrı ref'ler — daha net
const loading = ref(false)
const items = ref<Item[]>([])
```

#### Fonksiyon İsimlendirme
| Prefix | Kullanım | Örnek |
|--------|----------|-------|
| `load` | Veri çekme (sayfa açılışı) | `loadTasks`, `loadBoard` |
| `handle` | Kullanıcı etkileşimi (event) | `handleSubmit`, `handleDelete` |
| `toggle` | Açma/kapama | `toggleFilter`, `toggleSidebar` |
| `format` | Veri dönüştürme | `formatDate`, `formatPrice` |
| `is/has/can` | Boolean computed | `isOverdue`, `hasSubtasks`, `canEdit` |

#### Template Kuralları
- `v-if` / `v-else` zincirleri 3'ü geçmez — component'e taşı
- Event handler tek satırsa inline yaz: `@click="show = false"`
- Event handler 2+ satırsa fonksiyona taşı: `@click="handleDelete"`
- `class` binding'de Tailwind + koşullu class için array syntax değil, object syntax kullan

---

## 4. Dosya & Bileşen Kuralları

| Kural | Değer |
|-------|-------|
| Maksimum dosya satırı | 300 (zorunlu: 500) |
| Bir bileşen = bir sorumluluk | Tekil görev prensibi |
| Component isimlendirme | PascalCase → `TaskCard.vue` |
| Composable isimlendirme | camelCase + use prefix → `useBoard.ts` `tamamen tüm proje ingilizce olarak yapılandırılmalı`|
| Go handler isimlendirme | PascalCase → `CreateTask` |
| Test dosyası | `*.spec.ts` veya `*.test.ts` |
| Migration dosyası | `timestamp_description.sql` |
| Store dosyası | `use[Name]Store.ts` |

---

## 5. Proje Yapısı

```
kanbanmaster/
├── app/                          # Vue.js frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── common/           # Button, Modal, Input, Toast
│   │   │   ├── board/            # BoardView, Column, TaskCard
│   │   │   ├── task/             # TaskDetail, TaskForm, Checklist
│   │   │   ├── team/             # TeamList, MemberCard, InviteModal
│   │   │   ├── dashboard/        # StatCard, Charts, PerformancePanel
│   │   │   ├── notification/     # NotificationBell, NotificationList
│   │   │   └── layout/           # Navbar, Sidebar, Footer
│   │   ├── pages/                # Sayfa bileşenleri
│   │   │   ├── LoginPage.vue
│   │   │   ├── RegisterPage.vue
│   │   │   ├── DashboardPage.vue
│   │   │   ├── BoardPage.vue
│   │   │   ├── TeamPage.vue
│   │   │   ├── ReportsPage.vue
│   │   │   └── ProfilePage.vue
│   │   ├── composables/          # use* fonksiyonları
│   │   │   ├── useAuth.ts
│   │   │   ├── useBoard.ts
│   │   │   ├── useTask.ts
│   │   │   ├── useNotification.ts
│   │   │   └── useWebSocket.ts
│   │   ├── stores/               # Pinia store modülleri
│   │   │   ├── useAuthStore.ts
│   │   │   ├── useBoardStore.ts
│   │   │   ├── useTaskStore.ts
│   │   │   ├── useTeamStore.ts
│   │   │   └── useNotificationStore.ts
│   │   ├── services/             # API çağrıları (Axios)
│   │   │   ├── api.ts            # Axios instance + interceptor
│   │   │   ├── authService.ts
│   │   │   ├── boardService.ts
│   │   │   ├── taskService.ts
│   │   │   ├── teamService.ts
│   │   │   └── reportService.ts
│   │   ├── router/               # Vue Router
│   │   │   └── index.ts
│   │   ├── types/                # TypeScript tipleri
│   │   │   ├── user.ts
│   │   │   ├── board.ts
│   │   │   ├── task.ts
│   │   │   ├── team.ts
│   │   │   └── notification.ts
│   │   └── utils/                # Yardımcı fonksiyonlar
│   │       ├── date.ts
│   │       ├── validation.ts
│   │       └── format.ts
│   ├── public/
│   └── package.json
├── cmd/                          # Go backend
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── board.go
│   │   │   ├── column.go
│   │   │   ├── task.go
│   │   │   ├── team.go
│   │   │   ├── report.go
│   │   │   ├── notification.go
│   │   │   └── dashboard.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── error.go
│   │   ├── routes/
│   │   │   └── router.go
│   │   └── main.go
│   ├── models/
│   │   ├── user.go
│   │   ├── organization.go
│   │   ├── team.go
│   │   ├── board.go
│   │   ├── column.go
│   │   ├── task.go
│   │   ├── notification.go
│   │   └── report.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── board.go
│   │   ├── task.go
│   │   ├── team.go
│   │   ├── report.go
│   │   ├── notification.go
│   │   └── performance.go
│   ├── websocket/
│   │   └── hub.go
│   └── config/
│       └── config.go
├── lib/                          # Paylaşılan helper'lar
├── db/
│   ├── migrations/
│   └── schema.sql
├── docker/
│   ├── Dockerfile.api
│   ├── Dockerfile.app
│   └── docker-compose.yml
├── docs/
│   └── api-patterns.md
├── tests/                        # E2E testler
└── CLAUDE.md
```

---

## 6. Veritabanı Şeması

### Tablolar

**users**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| email | VARCHAR(255) | Unique, not null |
| password_hash | VARCHAR(255) | Not null |
| name | VARCHAR(100) | Not null |
| avatar_url | TEXT | Nullable |
| created_at | TIMESTAMP | Default now() |
| updated_at | TIMESTAMP | Default now() |

**organizations**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| name | VARCHAR(200) | Not null |
| owner_id | UUID | FK → users |
| created_at | TIMESTAMP | Default now() |

**teams**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| name | VARCHAR(200) | Not null |
| organization_id | UUID | FK → organizations |
| created_at | TIMESTAMP | Default now() |

**team_members**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| team_id | UUID | FK → teams |
| user_id | UUID | FK → users |
| role | ENUM | 'owner', 'leader', 'member', 'viewer' |
| joined_at | TIMESTAMP | Default now() |

**boards**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| name | VARCHAR(200) | Not null |
| team_id | UUID | FK → teams |
| created_at | TIMESTAMP | Default now() |

**columns**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| board_id | UUID | FK → boards |
| name | VARCHAR(100) | Not null |
| position | INTEGER | Sıralama |
| created_at | TIMESTAMP | Default now() |

**tasks**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| column_id | UUID | FK → columns |
| title | VARCHAR(300) | Not null |
| description | TEXT | Nullable, Markdown destekli |
| creator_id | UUID | FK → users (oluşturan) |
| assignee_id | UUID | FK → users (atanan) |
| priority | ENUM | 'urgent', 'high', 'medium', 'low' |
| deadline | TIMESTAMP | Nullable |
| position | INTEGER | Sütun içi sıralama |
| created_at | TIMESTAMP | Default now() |
| updated_at | TIMESTAMP | Default now() |
| completed_at | TIMESTAMP | Nullable, tamamlanma zamanı |

**subtasks**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| task_id | UUID | FK → tasks |
| title | VARCHAR(300) | Not null |
| is_completed | BOOLEAN | Default false |
| created_at | TIMESTAMP | Default now() |

**labels**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| board_id | UUID | FK → boards |
| name | VARCHAR(50) | Not null |
| color | VARCHAR(7) | Hex renk kodu |

**task_labels**
| Alan | Tip | Açıklama |
|------|-----|----------|
| task_id | UUID | FK → tasks |
| label_id | UUID | FK → labels |

**comments**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| task_id | UUID | FK → tasks |
| user_id | UUID | FK → users |
| content | TEXT | Not null |
| created_at | TIMESTAMP | Default now() |

**task_delegations**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| task_id | UUID | FK → tasks |
| from_user_id | UUID | FK → users (devreden) |
| to_user_id | UUID | FK → users (devralan) |
| reason | TEXT | Devretme nedeni |
| delegated_at | TIMESTAMP | Default now() |

**report_requests**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| requester_id | UUID | FK → users (rapor isteyen) |
| target_user_id | UUID | FK → users (rapor istenen) |
| team_id | UUID | FK → teams |
| message | TEXT | İstek mesajı |
| response | TEXT | Nullable, cevap |
| status | ENUM | 'pending', 'submitted', 'reviewed' |
| created_at | TIMESTAMP | Default now() |
| responded_at | TIMESTAMP | Nullable |

**notifications**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| user_id | UUID | FK → users (alıcı) |
| type | ENUM | 'assigned', 'delegated', 'deadline', 'comment', 'report_request', 'completed', 'overdue' |
| title | VARCHAR(300) | Bildirim başlığı |
| message | TEXT | Bildirim içeriği |
| reference_id | UUID | İlgili task/report ID |
| is_read | BOOLEAN | Default false |
| created_at | TIMESTAMP | Default now() |

**task_activity_log**
| Alan | Tip | Açıklama |
|------|-----|----------|
| id | UUID | Primary key |
| task_id | UUID | FK → tasks |
| user_id | UUID | FK → users |
| action | VARCHAR(50) | 'created', 'moved', 'assigned', 'delegated', 'completed', 'commented' |
| details | JSONB | Ek bilgi (ör: hangi sütundan hangi sütuna) |
| created_at | TIMESTAMP | Default now() |

---

## 7. API Endpoints

### Auth
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/auth/register | Kayıt |
| POST | /api/auth/login | Giriş |
| POST | /api/auth/refresh | Token yenileme |
| GET | /api/auth/me | Mevcut kullanıcı bilgisi |

### Organization
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/organizations | Kurum oluştur |
| GET | /api/organizations | Kurumlarımı listele |
| GET | /api/organizations/:id | Kurum detay |
| PUT | /api/organizations/:id | Kurum güncelle |
| DELETE | /api/organizations/:id | Kurum sil |

### Team
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/teams | Takım oluştur |
| GET | /api/teams | Takımlarımı listele |
| GET | /api/teams/:id | Takım detay |
| PUT | /api/teams/:id | Takım güncelle |
| DELETE | /api/teams/:id | Takım sil |
| POST | /api/teams/:id/invite | Üye davet et |
| DELETE | /api/teams/:id/members/:userId | Üye çıkar |
| PATCH | /api/teams/:id/members/:userId/role | Üye rolünü değiştir |

### Board
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/boards | Tahta oluştur |
| GET | /api/boards | Tahtalarımı listele |
| GET | /api/boards/:id | Tahta detay (sütunlar + görevlerle) |
| PUT | /api/boards/:id | Tahta güncelle |
| DELETE | /api/boards/:id | Tahta sil |

### Column
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/boards/:boardId/columns | Sütun ekle |
| PUT | /api/columns/:id | Sütun güncelle |
| DELETE | /api/columns/:id | Sütun sil |
| PATCH | /api/columns/reorder | Sütun sıralamasını değiştir |

### Task
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/tasks | Görev oluştur |
| GET | /api/tasks/:id | Görev detay |
| PUT | /api/tasks/:id | Görev güncelle |
| DELETE | /api/tasks/:id | Görev sil |
| PATCH | /api/tasks/:id/move | Görev taşı (sütun değiştir) |
| PATCH | /api/tasks/:id/assign | Görev ata |
| POST | /api/tasks/:id/delegate | Görev devret |
| GET | /api/tasks/:id/activity | Görev aktivite geçmişi |

### Subtask
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/tasks/:taskId/subtasks | Alt görev ekle |
| PATCH | /api/subtasks/:id/toggle | Alt görev tamamla/geri al |
| DELETE | /api/subtasks/:id | Alt görev sil |

### Comment
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/tasks/:taskId/comments | Yorum ekle |
| GET | /api/tasks/:taskId/comments | Yorumları listele |
| DELETE | /api/comments/:id | Yorum sil |

### Report
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /api/reports/request | Rapor iste |
| GET | /api/reports/requests | Gelen rapor istekleri |
| POST | /api/reports/requests/:id/respond | Rapor isteğine cevap ver |
| GET | /api/reports/requests/sent | Gönderilen rapor istekleri |

### Dashboard & Performans
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| GET | /api/dashboard/summary | Genel özet istatistikler |
| GET | /api/dashboard/team/:teamId | Takım bazlı istatistikler |
| GET | /api/dashboard/user/:userId/performance | Kişi performans verileri |
| GET | /api/dashboard/team/:teamId/performance | Takım performans verileri |
| GET | /api/dashboard/overdue | Geciken görevler |
| GET | /api/dashboard/weekly-report | Haftalık rapor |
| GET | /api/dashboard/monthly-report | Aylık rapor |

### Notification
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| GET | /api/notifications | Bildirimlerimi listele |
| PATCH | /api/notifications/:id/read | Bildirimi okundu işaretle |
| PATCH | /api/notifications/read-all | Tümünü okundu işaretle |
| WS | /ws/notifications | WebSocket — anlık bildirimler |

---

## 8. Özellik Detayları

### 8.1 Görev Atama & Devretme

Takım lideri şunları yapabilir:
- Bir görevi herhangi bir takım üyesine atayabilir
- Atanmış bir görevi başka birine devredebilir
- Devretme sırasında neden belirtir
- Devir geçmişi `task_delegations` tablosunda tutulur
- Atanan ve devralan kişiye otomatik bildirim gider

### 8.2 Rapor İsteme Sistemi

- Takım lideri, atadığı kişiden rapor talep edebilir
- Rapor isteği bildirim olarak gider
- Üye, rapor isteğine metin ile cevap verir
- Lider cevabı inceleyip "reviewed" olarak işaretler
- Durumlar: `pending` → `submitted` → `reviewed`

### 8.3 Aktif Bilgilendirme (Notification) Sistemi

Otomatik bildirim tetikleyicileri:
| Olay | Bildirim Alıcısı |
|------|------------------|
| Görev atandı | Atanan kişi |
| Görev devredildi | Yeni atanan kişi |
| Deadline 24 saat kala | Atanan kişi |
| Deadline geçti (gecikme) | Atanan kişi + takım lideri |
| Göreve yorum yapıldı | Görev sahibi + atanan kişi |
| Rapor istendi | Hedef kişi |
| Rapor cevaplandı | İsteyen kişi |
| Görev tamamlandı | Görevi oluşturan + takım lideri |

### 8.4 Performans Görselleştirme

**Kişi Bazlı Performans Paneli:**
- Tamamlanan görev sayısı (haftalık/aylık) → **Bar Chart**
- Zamanında vs gecikmeli tamamlama oranı → **Donut/Pie Chart**
- Zaman içindeki performans trendi → **Line Chart**
- Verimlilik skoru (0-100) → **Gauge/Radial Chart**
  - Hesaplama: `(zamanında_tamamlanan / toplam_atanan) × 100`

**Takım Bazlı Dashboard:**
- Takım üyeleri iş yükü dağılımı → **Horizontal Bar Chart**
- Sütun bazlı görev dağılımı → **Stacked Bar Chart**
- Haftalık tamamlanma trendi → **Area Chart**
- Geciken görevler listesi → **Tablo + kırmızı vurgu**
- En verimli üyeler sıralaması → **Leaderboard**

**Aylık/Haftalık Rapor Görünümü:**
- Toplam atanan görev
- Tamamlanan görev
- Geciken görev
- Ortalama tamamlama süresi
- Performans değişimi (önceki döneme kıyasla ↑↓)

---

## 9. UI/UX Tasarım Sistemi

> **21st.dev Magic MCP** entegrasyonu ile profesyonel, modern ve tutarlı UI bileşenleri kullanılacaktır.
> Tüm arayüz bileşenleri Magic MCP üzerinden üretilecek ve Vuetify + Tailwind ile şekillendirilecektir.

### 9.1 Tasarım Prensipleri

| Prensip | Açıklama |
|---------|----------|
| **Clean & Minimal** | Gereksiz görsel eleman yok, beyaz alan (whitespace) bolca kullanılır |
| **Tutarlılık** | Tüm sayfalarda aynı spacing, tipografi ve renk sistemi |
| **Erişilebilirlik** | WCAG 2.1 AA uyumlu, kontrast oranları yeterli |
| **Mikro Etkileşimler** | Hover, focus, aktif durumlar için akıcı geçişler (transition 200-300ms) |
| **Veri Yoğunluğu** | Dashboard ve tablolarda bilgi kaybı olmadan kompakt görünüm |

### 9.2 Renk Sistemi

| Renk | Kullanım | Hex |
|------|----------|-----|
| Primary | Ana butonlar, aktif menü, linkler | `#1E88E5` (Blue 600) |
| Secondary | İkincil aksiyonlar, badge'ler | `#7C4DFF` (Deep Purple A200) |
| Success | Tamamlanan görevler, onay | `#43A047` (Green 600) |
| Warning | Deadline yaklaşan, dikkat | `#FB8C00` (Orange 600) |
| Error | Geciken görevler, hata | `#E53935` (Red 600) |
| Urgent | Acil öncelik | `#D50000` (Red A700) |
| Background | Sayfa arka planı | `#F5F5F5` (Grey 100) |
| Surface | Kart ve panel arka planı | `#FFFFFF` |
| Text Primary | Ana metin | `#212121` (Grey 900) |
| Text Secondary | Açıklama, meta bilgi | `#757575` (Grey 600) |

### 9.3 Tipografi

| Eleman | Font | Boyut | Ağırlık |
|--------|------|-------|---------|
| Heading 1 | Inter | 28px | 700 (Bold) |
| Heading 2 | Inter | 22px | 600 (Semi-bold) |
| Heading 3 | Inter | 18px | 600 |
| Body | Inter | 14px | 400 (Regular) |
| Caption | Inter | 12px | 400 |
| Button | Inter | 14px | 500 (Medium) |
| Code/Mono | JetBrains Mono | 13px | 400 |

### 9.4 Spacing & Grid

- **Base unit:** 4px
- **Spacing scale:** 4, 8, 12, 16, 24, 32, 48, 64px
- **Border radius:** Kartlar 12px, Butonlar 8px, Input 8px, Chip 16px
- **Gölge (elevation):** Vuetify elevation-1 (kartlar), elevation-3 (modal/dropdown)
- **Max content width:** 1440px
- **Sidebar width:** 260px (kapatılabilir)

### 9.5 Sayfa Bazlı UI Tasarım Detayları

#### Login / Register Sayfası
- Ortada tek kart, split layout (sol: form, sağ: görsel/branding)
- Gradient arka plan (`Primary → Secondary` soft geçiş)
- Form validation — anlık feedback (inline error mesajları)
- "Hoş geldiniz" / "Hesap oluşturun" başlık + alt metin
- Social login butonları (gelecek için placeholder)

#### Dashboard Sayfası
- **Üst bölüm:** 4 adet StatCard (Toplam Görev, Tamamlanan, Devam Eden, Geciken)
  - Her kart: ikon + sayı + önceki döneme kıyasla değişim (↑↓ %)
  - Renk kodlu: success/warning/error
- **Orta bölüm:** 2 sütunlu grid
  - Sol: Haftalık tamamlanma trendi (Area Chart)
  - Sağ: Görev dağılımı (Donut Chart)
- **Alt bölüm:** Takım performans tablosu + geciken görevler listesi
- Kişiye özel selamlama: "Merhaba, {isim}" + bugünün tarihi

#### Kanban Tahtası Sayfası
- **Üst bar:** Tahta adı + filtre butonları + arama + "Yeni Görev" butonu
- **Sütunlar:** Yatay scroll, her sütun başlığında görev sayısı badge
- **Görev Kartı tasarımı:**
  - Üst: Öncelik renk şeridi (4px üst border)
  - Başlık (bold, max 2 satır truncate)
  - Etiketler (renkli chip'ler, max 3 görünür + "+2 more")
  - Alt bölüm: Avatar (atanan kişi) + deadline (ikon + tarih) + alt görev sayısı (✓ 3/5)
  - Hover: hafif elevation artışı + gölge geçişi
- **Boş sütun:** Dashed border + "Görev ekle" placeholder
- **Drag ghost:** Opaklık %70, hafif döndürme (2deg)

#### Görev Detay Modalı
- Sağdan slide-in panel (drawer) — tam sayfa modal değil
- Genişlik: 480px
- Bölümler: Başlık, Açıklama (Markdown editor), Atama, Öncelik, Deadline, Etiketler, Alt Görevler, Yorumlar, Aktivite Geçmişi
- Her bölüm collapse/expand yapılabilir
- Sağ üst: Kapat (X) + Üç nokta menüsü (devret, sil, arşivle)

#### Takım Yönetim Sayfası
- Sol: Takım listesi (kartlar halinde)
- Sağ: Seçili takımın üye tablosu
- Her üye satırı: Avatar + İsim + Rol badge + Aksiyon butonları
- Üst: "Üye Davet Et" butonu → modal (email ile davet)

#### Performans / Raporlar Sayfası
- **Kişi seçici:** Dropdown ile takım üyesi seçimi
- **Üst kartlar:** Verimlilik skoru (Gauge), Toplam görev, Zamanında tamamlanan, Geciken
- **Grafikler bölümü:**
  - Haftalık/Aylık toggle ile geçiş
  - Tamamlama trendi (Line Chart — son 12 hafta/ay)
  - Zamanında vs gecikmeli (Stacked Bar Chart)
  - Öncelik bazlı dağılım (Horizontal Bar)
- **Takım karşılaştırma:** Leaderboard tablosu — sıralama, isim, skor, tamamlanan, geciken
- **Isı haritası (Heatmap):** Haftalık aktivite yoğunluğu (GitHub contribution graph benzeri)

#### Bildirim Paneli
- Navbar'da bell ikonu + okunmamış sayısı (kırmızı badge)
- Tıklandığında dropdown panel (max 400px genişlik)
- Her bildirim: ikon (tip bazlı) + başlık + zaman (relative: "2 saat önce") + okundu durumu
- Alt kısım: "Tümünü gör" linki → ayrı bildirim sayfası
- Okunmamışlar hafif mavi arka plan ile vurgulu

#### Profil Sayfası
- Avatar yükleme alanı (circular crop)
- İsim, email düzenleme formu
- Şifre değiştirme bölümü (ayrı kart)
- Bildirim tercihleri (toggle switch'ler)

### 9.6 Bileşen Kütüphanesi (Component Library)

| Bileşen | Kullanım | Kaynak |
|---------|----------|--------|
| AppButton | Tüm butonlar (primary, secondary, outlined, text, icon) | Vuetify v-btn + Tailwind |
| AppCard | Kart yapıları (stat, görev, profil) | Vuetify v-card |
| AppModal | Dialog ve drawer paneller | Vuetify v-dialog / v-navigation-drawer |
| AppInput | Form input alanları (text, email, password, textarea) | Vuetify v-text-field |
| AppSelect | Dropdown seçiciler | Vuetify v-select / v-autocomplete |
| AppChip | Etiketler, öncelik badge, rol badge | Vuetify v-chip |
| AppAvatar | Kullanıcı avatar görselleri | Vuetify v-avatar |
| AppToast | Bildirim toast/snackbar | Vuetify v-snackbar |
| AppTable | Veri tabloları (üye listesi, geciken görevler) | Vuetify v-data-table |
| AppSkeleton | Loading durumları | Vuetify v-skeleton-loader |
| AppEmptyState | Boş durum ekranları | Custom (ikon + mesaj + CTA butonu) |
| StatCard | Dashboard istatistik kartı | Custom (ikon + değer + değişim) |
| TaskCard | Kanban görev kartı | Custom (öncelik + etiket + avatar + deadline) |
| ChartWrapper | Chart.js grafik sarmalayıcı | vue-chartjs + custom props |
| PerformanceGauge | Verimlilik skoru göstergesi | Chart.js doughnut (gauge mode) |

### 9.7 Animasyon & Geçişler

| Eleman | Animasyon | Süre |
|--------|-----------|------|
| Sayfa geçişi | Fade + slide-up | 250ms ease |
| Modal açılış | Fade backdrop + slide-right (drawer) | 300ms ease-out |
| Kart hover | Elevation artışı + scale(1.01) | 200ms ease |
| Drag & drop | Placeholder highlight + ghost opacity | 150ms |
| Toast | Slide-up + fade-in | 200ms ease |
| Grafik yükleme | Sayısal artış animasyonu (count-up) | 600ms |
| Chip ekleme/silme | Scale-in / scale-out | 150ms |
| Sidebar toggle | Slide-left / width transition | 250ms ease |

### 9.8 Responsive Breakpoints

| Breakpoint | Genişlik | Davranış |
|------------|----------|----------|
| Mobile | < 640px | Tek sütun, sidebar gizli, bottom nav |
| Tablet | 640–1024px | 2 sütun grid, sidebar overlay |
| Desktop | 1024–1440px | Tam layout, sidebar sabit |
| Wide | > 1440px | Max-width container, ortalanmış |

### 9.9 Karanlık Mod (Dark Mode) — V2.0

| Eleman | Light | Dark |
|--------|-------|------|
| Background | `#F5F5F5` | `#121212` |
| Surface | `#FFFFFF` | `#1E1E1E` |
| Text Primary | `#212121` | `#E0E0E0` |
| Text Secondary | `#757575` | `#9E9E9E` |
| Border | `#E0E0E0` | `#333333` |

### 9.10 Genel UX Kuralları

| Alan | Yaklaşım |
|------|----------|
| Boş durum | İlk kez giren kullanıcıya rehber ekran ("İlk tahtanızı oluşturun") |
| Loading | Skeleton loader (iskelet ekran), spinner sadece inline aksiyonlarda |
| Hata | Kullanıcı dostu mesajlar + yeniden deneme butonu |
| Drag & Drop | Akıcı animasyonlarla sürükle-bırak |
| Mobil | Responsive — mobilde sütunlar yatay kaydırmalı |
| Çoklu Görünüm | Kanban / Liste / Takvim (V2.0) |
| Klavye kısayolları | N: yeni görev, F: filtre, /: arama, Esc: modal kapat |
| Feedback | Her aksiyon sonrası toast bildirimi (başarılı/hata) |

---

## 10. Hata Yönetimi

| Katman | Yöntem |
|--------|--------|
| Frontend API | Axios interceptor + try/catch → toast/snackbar bildirimi |
| Frontend UI | Vue error boundary, fallback bileşenleri |
| Backend API | Merkezi error handler middleware → `{error: string, code: number}` |
| Database | Transaction kullanımı, constraint violation yakalama |
| Validation | Frontend: Vuetify form rules / Backend: request validation |

---

## 11. Test Stratejisi

| Tür | Araç | Kapsam |
|-----|------|--------|
| Unit (Frontend) | Vitest | Composable, store, util fonksiyonları |
| Component | Vitest + Vue Test Utils | Bileşen render ve etkileşim |
| Unit (Backend) | Go testing | Handler, service, model |
| API | Go httptest | Endpoint doğrulama |
| E2E | Playwright | Kritik kullanıcı akışları |

**Kural:** Her faz sonunda ilgili testler yazılır ve geçer. Testsiz kod merge edilmez.

### 11.1 Hızlı Test Komutları

> **KRİTİK KURAL:** Her kod değişikliğinden sonra ilgili dosyanın testi MUTLAKA çalıştırılmalıdır.
> Değişiklik yapılıp test edilmeden bir sonraki adıma geçilmez.

#### Frontend Test Komutları

```bash
# Tüm frontend testlerini çalıştır
cd app && npm test

# Tek bir dosyanın testini çalıştır (hızlı test)
cd app && npx vitest run src/utils/date.spec.ts

# Belirli bir klasördeki testleri çalıştır
cd app && npx vitest run src/stores/

# Watch modunda test (geliştirirken)
cd app && npx vitest watch src/components/board/

# Coverage raporu ile
cd app && npx vitest run --coverage

# TypeScript type-check (her değişiklik sonrası zorunlu)
cd app && npx vue-tsc --noEmit
```

#### Backend Test Komutları

```bash
# Tüm Go testlerini çalıştır
go test ./...

# Tek bir paketi test et (hızlı test)
go test ./cmd/api/handlers/ -v

# Tek bir test fonksiyonunu çalıştır
go test ./cmd/services/ -run TestCreateTask -v

# Race condition kontrolü ile
go test ./... -race

# Coverage raporu ile
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Benchmark testleri
go test ./cmd/services/ -bench=. -benchmem
```

#### Değişiklik Sonrası Zorunlu Adımlar

Her kod değişikliğinden sonra sırasıyla:

1. **Typecheck** → `cd app && npx vue-tsc --noEmit` (frontend) veya `go build ./...` (backend)
2. **İlgili test** → Değiştirilen dosyanın testi çalıştırılır
3. **Build kontrolü** → `cd app && npx vite build` (frontend) veya `go build ./cmd/api` (backend)

```bash
# Tek komutla frontend hızlı doğrulama
cd app && npx vue-tsc --noEmit && npx vitest run && npx vite build

# Tek komutla backend hızlı doğrulama
go build ./... && go test ./... && go vet ./...
```

#### Test Dosyası İsimlendirme

| Dosya | Test Dosyası |
|-------|-------------|
| `src/utils/date.ts` | `src/utils/date.spec.ts` |
| `src/stores/useAuthStore.ts` | `src/stores/useAuthStore.spec.ts` |
| `src/components/board/TaskCard.vue` | `src/components/board/TaskCard.spec.ts` |
| `cmd/api/handlers/auth.go` | `cmd/api/handlers/auth_test.go` |
| `cmd/services/task.go` | `cmd/services/task_test.go` |

---

## 12. Geliştirme Aşamaları (Fazlar)

### FAZ 1 — Temel Altyapı & Proje İskeleti
> **Amaç:** Projeyi ayağa kaldırmak, geliştirme ortamını hazırlamak.

- [ ] Git repo oluştur, branch stratejisi belirle (main, develop, feature/*)
- [ ] Docker Compose yapılandırması (API + DB + Frontend)
- [ ] PostgreSQL bağlantısı ve migration sistemi kur
- [ ] Go API iskeleti (router, middleware, config, health-check endpoint)
- [ ] Vue projesi kurulumu (Vite + Vuetify + Tailwind + Router + Pinia)
- [ ] Axios instance + interceptor yapısı
- [ ] Temel layout bileşenleri (Navbar, Sidebar, Footer)
- [ ] Environment değişkenleri (.env) yapılandırması
- [ ] **Çıktı:** `docker-compose up` ile tüm servisler ayağa kalkıyor
- [ ] **Test:** Health-check endpoint çalışıyor, DB bağlantısı aktif

---

### FAZ 2 — Kimlik Doğrulama (Auth)
> **Amaç:** Kullanıcı giriş/çıkış sistemini tamamlamak.

- [ ] `users` tablosu migration
- [ ] Register endpoint (şifre hashleme: bcrypt)
- [ ] Login endpoint (JWT access + refresh token üretimi)
- [ ] Refresh token endpoint
- [ ] Auth middleware (JWT doğrulama)
- [ ] Auth store (Pinia) — login, logout, token yönetimi
- [ ] Login sayfası
- [ ] Register sayfası
- [ ] Route guard (korumalı sayfalar yönlendirmesi)
- [ ] **Çıktı:** Kullanıcı kayıt olup giriş yapabiliyor, korumalı sayfalara erişebiliyor
- [ ] **Test:** Auth akışı uçtan uca çalışıyor, geçersiz token reddediliyor

---

### FAZ 3 — Kurum & Takım Yönetimi
> **Amaç:** Organizasyon yapısını kurmak, takım ve üye yönetimini sağlamak.

- [ ] `organizations`, `teams`, `team_members` tabloları migration
- [ ] Organization CRUD endpoint'leri
- [ ] Team CRUD endpoint'leri
- [ ] Üye davet etme, çıkarma, rol değiştirme endpoint'leri
- [ ] Rol bazlı yetkilendirme middleware (owner, leader, member, viewer)
- [ ] Kurum yönetim sayfası
- [ ] Takım yönetim sayfası
- [ ] Üye listesi ve davet modalı
- [ ] Team store (Pinia)
- [ ] **Çıktı:** Kurum oluşturup takım ekleyebiliyor, üye davet edebiliyor
- [ ] **Test:** CRUD işlemleri, yetkilendirme kontrolleri

---

### FAZ 4 — Kanban Tahtası (Çekirdek)
> **Amaç:** Projenin kalbini oluşturmak — görev yönetimi ve Kanban tahtası.

- [ ] `boards`, `columns`, `tasks`, `subtasks` tabloları migration
- [ ] Board CRUD endpoint'leri
- [ ] Column CRUD + reorder endpoint'leri
- [ ] Task CRUD + move + assign endpoint'leri
- [ ] Subtask CRUD + toggle endpoint'leri
- [ ] Kanban tahtası görünümü (sütunlar + görev kartları)
- [ ] vuedraggable ile sürükle-bırak
- [ ] Görev detay modalı (başlık, açıklama, atama, öncelik, deadline, alt görevler)
- [ ] Board store + Task store (Pinia)
- [ ] **Çıktı:** Tam çalışan Kanban tahtası — görev oluşturma, taşıma, atama
- [ ] **Test:** Drag & drop, görev CRUD, sütun sıralama

---

### FAZ 5 — Görev Devretme & Rapor İsteme
> **Amaç:** Takım liderlerinin yönetim araçlarını eklemek.

- [ ] `task_delegations`, `report_requests`, `task_activity_log` tabloları migration
- [ ] Görev devretme endpoint'i (delegate)
- [ ] Rapor isteme/cevaplama endpoint'leri
- [ ] Aktivite log kayıt sistemi (her görev eylemi loglanır)
- [ ] Görev devretme UI (modal + neden alanı)
- [ ] Rapor isteme/cevaplama UI
- [ ] Görev aktivite geçmişi görünümü
- [ ] **Çıktı:** Lider görev devredebiliyor, rapor isteyip cevap alabiliyor
- [ ] **Test:** Devretme akışı, rapor durumları, aktivite logu

---

### FAZ 6 — Etiketler, Yorumlar & Filtreleme
> **Amaç:** Görevleri zenginleştirmek ve aranabilir kılmak.

- [ ] `labels`, `task_labels`, `comments` tabloları migration
- [ ] Label CRUD endpoint'leri
- [ ] Comment CRUD endpoint'leri
- [ ] Filtreleme endpoint'leri (kişi, öncelik, etiket, tarih)
- [ ] Etiket yönetimi UI (renk seçimi, ekleme/çıkarma)
- [ ] Yorum alanı UI (görev detayında)
- [ ] Filtre paneli UI (sidebar veya dropdown)
- [ ] Arama çubuğu (serbest metin)
- [ ] **Çıktı:** Görevlere etiket/yorum eklenebiliyor, filtrelenebiliyor
- [ ] **Test:** Filtreleme kombinasyonları, yorum CRUD

---

### FAZ 7 — Bildirim Sistemi & Real-time
> **Amaç:** Kullanıcıları anlık bilgilendirmek.

- [ ] `notifications` tablosu migration
- [ ] Bildirim oluşturma servisi (olay tetikleyicileri)
- [ ] Bildirim endpoint'leri (listeleme, okundu işaretleme)
- [ ] WebSocket hub (gorilla/websocket)
- [ ] Anlık bildirim gönderimi (WebSocket üzerinden)
- [ ] Bildirim ikonu (bell) + dropdown listesi
- [ ] Okundu/okunmadı durumu
- [ ] Deadline yaklaşma uyarıları (cron job veya scheduler)
- [ ] Notification store (Pinia)
- [ ] **Çıktı:** Kullanıcılar anlık bildirim alıyor, geçmiş bildirimleri görebiliyor
- [ ] **Test:** Bildirim tetikleyicileri, WebSocket bağlantısı, okundu durumu

---

### FAZ 8 — Dashboard & Performans Görselleştirme
> **Amaç:** Yöneticilere ve kullanıcılara performans verilerini sunmak.

- [ ] Dashboard istatistik endpoint'leri
- [ ] Kişi performans endpoint'i (tamamlama oranı, gecikme, trend)
- [ ] Takım performans endpoint'i
- [ ] Haftalık/aylık rapor endpoint'leri
- [ ] Dashboard sayfası — genel özet kartları
- [ ] Kişi performans paneli (Bar, Pie, Line, Gauge chart'lar)
- [ ] Takım performans paneli (Stacked Bar, Area, Leaderboard)
- [ ] Geciken görevler tablosu (kırmızı vurgulu)
- [ ] Haftalık/aylık karşılaştırma görünümü (↑↓ değişim)
- [ ] **Çıktı:** Görsel performans paneli — grafiklerle kişi/takım analizi
- [ ] **Test:** Doğru veri hesaplaması, grafik render, edge case'ler

---

### FAZ 9 — Son Rötuşlar & Production Hazırlık
> **Amaç:** Projeyi production'a hazır hale getirmek.

- [ ] Responsive tasarım kontrolü (mobil, tablet, desktop)
- [ ] Hata sayfaları (404, 500, 403)
- [ ] Loading state'leri ve boş durum ekranları
- [ ] Profil sayfası (isim, avatar, şifre değiştirme)
- [ ] Production Docker build optimizasyonu
- [ ] Environment bazlı config (dev, staging, prod)
- [ ] API rate limiting
- [ ] Console.log temizliği
- [ ] Son kod gözden geçirme (300 satır limiti kontrolü)
- [ ] E2E testler (kritik akışlar)
- [ ] **Çıktı:** Production'a deploy edilebilir, tam test edilmiş uygulama
- [ ] **Test:** E2E test suite geçiyor, production build sorunsuz

---

## 13. Roller & Yetkiler Matrisi

| İşlem | Owner | Leader | Member | Viewer |
|-------|-------|--------|--------|--------|
| Kurum ayarları | ✅ | ❌ | ❌ | ❌ |
| Takım oluştur/sil | ✅ | ❌ | ❌ | ❌ |
| Takım ayarları | ✅ | ✅ | ❌ | ❌ |
| Üye davet/çıkar | ✅ | ✅ | ❌ | ❌ |
| Tahta oluştur/sil | ✅ | ✅ | ❌ | ❌ |
| Sütun ekle/sil | ✅ | ✅ | ❌ | ❌ |
| Görev oluştur | ✅ | ✅ | ✅ | ❌ |
| Görev ata | ✅ | ✅ | ❌ | ❌ |
| Görev devret | ✅ | ✅ | ❌ | ❌ |
| Görev taşı | ✅ | ✅ | ✅ | ❌ |
| Rapor iste | ✅ | ✅ | ❌ | ❌ |
| Rapor cevapla | ✅ | ✅ | ✅ | ❌ |
| Yorum yap | ✅ | ✅ | ✅ | ❌ |
| Dashboard görüntüle | ✅ | ✅ | ✅ | ✅ |
| Performans paneli | ✅ | ✅ | 🔸 | ❌ |

> 🔸 Member sadece kendi performansını görebilir

---

## 14. Ücretsiz Deployment (Canlıya Çıkma)

> **Platform: Render.com** — Frontend (static), Backend (Go), PostgreSQL hepsi ücretsiz tier'da.

### Adımlar

1. **GitHub'a push et:**
   ```bash
   git remote add origin https://github.com/KULLANICI/kanbanmaster.git
   git push -u origin main
   ```

2. **Render.com'a git:** https://render.com → GitHub ile giriş yap

3. **Blueprint ile deploy:**
   - "New" → "Blueprint" → GitHub repo'yu seç
   - `render.yaml` dosyası otomatik algılanır
   - 3 servis oluşturulur: API + Frontend + PostgreSQL

4. **İlk kurulum sonrası:**
   - PostgreSQL bağlantısı otomatik
   - `JWT_SECRET` otomatik üretilir
   - Schema'yı yükle: Render Shell → `psql $DATABASE_URL < db/schema.sql`

### Alternatif Ücretsiz Platformlar

| Platform | Frontend | Backend | DB | Not |
|----------|----------|---------|----|-----|
| **Render.com** | Static (free) | Go (free) | PostgreSQL (free) | En kolay, blueprint ile |
| **Vercel + Railway** | Vercel (free) | Railway (free) | Railway PostgreSQL | Ayrı deploy |
| **Netlify + Fly.io** | Netlify (free) | Fly.io (free) | Fly.io PostgreSQL | Docker gerekli |
| **Cloudflare Pages + Supabase** | CF Pages (free) | Workers (free) | Supabase (free) | Edge tabanlı |

### Render.com Ücretsiz Limitler
- Web servisleri 15 dakika inaktivitede uyur (ilk istek ~30sn)
- 750 saat/ay (yeterli)
- PostgreSQL 1GB, 90 gün sonra silinir (yenilenmeli)
- Custom domain desteği var

---

*Son güncelleme: 2026-03-23*
