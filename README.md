# Zacode Go Backend

Backend server untuk aplikasi Zacode menggunakan Go dengan arsitektur clean architecture.

## Struktur Proyek

```
/yourapp
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/          # load env, config global
│   │   └── config.go
│   │
│   ├── app/             # HTTP handler + routing (Gin)
│   │   ├── router.go
│   │   ├── auth_handler.go
│   │   ├── chat_handler.go
│   │   └── ...
│   │
│   ├── service/         # business logic / usecase
│   │   ├── auth_service.go
│   │   ├── chat_service.go
│   │   └── ...
│   │
│   ├── repository/      # DB access (gorm / raw SQL)
│   │   ├── user_repo.go
│   │   ├── chat_repo.go
│   │   └── ...
│   │
│   ├── model/           # struct model untuk DB
│   │   ├── user.go
│   │   ├── chat.go
│   │   └── ...
│   │
│   ├── websocket/       # ws hub, manager, client
│   │   ├── hub.go
│   │   ├── client.go
│   │   ├── ws_handler.go
│   │   └── ...
│   │
│   └── util/            # helper: jwt, hash, error, response
│       ├── jwt.go
│       ├── hash.go
│       └── response.go
│
├── pkg/                 # library reusable (optional)
│   └── logger/
│       └── logger.go
│
├── go.mod
├── .env
├── Dockerfile
└── docker-compose.yml
```

## Deskripsi Folder

### `cmd/server/`
Entry point aplikasi. Berisi `main.go` yang menginisialisasi dan menjalankan server.

### `internal/config/`
Konfigurasi aplikasi, termasuk loading environment variables dan setup global config.

### `internal/app/`
Layer HTTP handler dan routing menggunakan Gin framework.
- `router.go`: Setup routing dan middleware
- `*_handler.go`: HTTP handlers untuk setiap endpoint

### `internal/service/`
Business logic layer (use case layer). Berisi logika bisnis aplikasi.

### `internal/repository/`
Data access layer. Interface dan implementasi untuk akses database (GORM atau raw SQL).

### `internal/model/`
Struct model untuk database. Definisi struct yang digunakan untuk mapping database.

### `internal/websocket/`
WebSocket implementation untuk real-time communication.
- `hub.go`: WebSocket hub untuk manage connections
- `client.go`: WebSocket client implementation
- `ws_handler.go`: WebSocket handler

### `internal/util/`
Utility functions dan helpers:
- `jwt.go`: JWT token generation dan validation
- `hash.go`: Password hashing utilities
- `response.go`: Standard response formatter

### `pkg/logger/`
Reusable logger library yang bisa digunakan di seluruh aplikasi.

## Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL
- Redis (optional)
- RabbitMQ (optional)

### Installation

1. Clone repository
```bash
git clone <repository-url>
cd /go
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Update `.env` dengan konfigurasi yang sesuai

4. Install dependencies
```bash
go mod download
```

5. Run dengan Docker Compose
```bash
docker-compose up -d
```

6. Atau run secara lokal
```bash
go run cmd/server/main.go
```

## Environment Variables

Buat file `.env` dengan variabel berikut:

```env
# Server
PORT=5000
SERVER_HOST=0.0.0.0
CLIENT_URL=http://localhost:3000

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret_key

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=your_user
RABBITMQ_PASSWORD=your_password
```

## Development

### Run Development Server
```bash
go run cmd/server/main.go
```

### Build
```bash
go build -o bin/server cmd/server/main.go
```

### Run Tests
```bash
go test ./...
```

## Docker

### Build Image
```bash
docker build -t 
```

### Run with Docker Compose
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

## Services & Ports

Setelah menjalankan `docker-compose up -d`, services berikut akan tersedia:

- **API Server**: http://localhost:5000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **RabbitMQ Management UI**: http://localhost:15672
  - Username: `yourapp` (default)
  - Password: `password123` (default)
- **pgweb (Database UI)**: http://localhost:8081
  - Web-based PostgreSQL client untuk melihat dan mengelola database
  - Otomatis terhubung ke database yang dikonfigurasi

## Architecture

Aplikasi ini menggunakan **Clean Architecture** dengan layer separation:

1. **Handler Layer** (`internal/app/`): HTTP handlers dan routing
2. **Service Layer** (`internal/service/`): Business logic
3. **Repository Layer** (`internal/repository/`): Data access
4. **Model Layer** (`internal/model/`): Domain models

## License

MIT

# Flowchart Sistem Lelang Online dengan Logic Pilih Item Dulu

```
                    ┌─────────────────────────┐
                    │  User Membuka Aplikasi  │
                    └───────────┬─────────────┘
                                │
                                ▼
                        ┌───────────────┐
                        │ Status Login? │
                        └───┬───────┬───┘
                            │       │
                ┌───────────┘       └───────────┐
                │                               │
        Belum Login                       Sudah Login
                │                               │
                ▼                               ▼
        ┌───────────────┐               ┌────────────────┐
        │ Mode Guest    │               │ Dashboard User │
        └───────┬───────┘               └────┬───┬───┬───┘
                │                            │   │   │
                ▼                            │   │   │
        ┌─────────────────┐           ┌─────┘   │   └─────┐
        │ Lihat Daftar    │           │         │         │
        │ Lelang          │           ▼         ▼         ▼
        └────────┬────────┘    ┌──────────┐ ┌────────┐ ┌──────────┐
                 │             │ Lihat    │ │ Profil │ │ Riwayat  │
                 ▼             │ Lelang   │ │ & Saldo│ │ Bidding  │
        ┌──────────────────┐  └────┬─────┘ └───┬────┘ └────┬─────┘
        │ Pilih Item       │       │           │           │
        │ Lelang?          │◄──────┘           │           │
        └───┬──────────┬───┘                   │           │
            │          │                       │           │
        Tidak          Ya                      ▼           ▼
            │          │                ┌──────────────┐ ┌────────────┐
            │          ▼                │ Lihat Saldo  │ │ Lihat      │
            │  ┌────────────────┐      └──────┬───────┘ │ Riwayat    │
            │  │ Lihat Detail   │             │         └─────┬──────┘
            │  │ Lelang Item    │             ▼               │
            │  │ Terpilih       │      ┌───────────┐          │
            │  └────────┬───────┘      │ Top Up?   │          │
            │           │              └──┬────┬───┘          │
            │           ▼                 │    │              │
            │    ┌─────────────────┐  Tidak  Ya              │
            │    │ Ingin Ikut      │     │    │              │
            │    │ Bidding?        │     │    │              │
            │    └──┬──────────┬───┘     │    │              │
            │       │          │         │    │              │
            │   Tidak          Ya        │    │              │
            │       │          │         │    │              │
            └───────┘          ▼         │    │              │
                    ┌──────────────────┐ │    │              │
                    │ Diarahkan ke     │ │    │              │
                    │ Halaman Login    │ │    │              │
                    └─────────┬────────┘ │    │              │
                              │          │    │              │
                              ▼          │    │              │
                    ┌──────────────────┐ │    │              │
                    │ Halaman Login/   │ │    │              │
                    │ Register         │ │    │              │
                    └─────────┬────────┘ │    │              │
                              │          │    │              │
                              ▼          │    │              │
                      ┌────────────┐     │    │              │
                      │ Sudah Punya│     │    │              │
                      │ Akun?      │     │    │              │
                      └──┬─────┬───┘     │    │              │
                         │     │         │    │              │
                    Belum│     │Sudah    │    │              │
                         │     │         │    │              │
                         ▼     ▼         │    │              │
              ┌──────────┐   ┌──────┐   │    │              │
              │ Form     │   │ Form │   │    │              │
              │ Register │   │ Login│   │    │              │
              └────┬─────┘   └───┬──┘   │    │              │
                   │             │      │    │              │
                   ▼             │      │    │              │
         ┌────────────────┐     │      │    │              │
         │ Verifikasi     │     │      │    │              │
         │ Akun Email/    │     │      │    │              │
         │ Phone          │     │      │    │              │
         └────────┬───────┘     │      │    │              │
                  │             │      │    │              │
                  └──────┬──────┘      │    │              │
                         │             │    │              │
                         ▼             │    │              │
                ┌────────────────┐     │    │              │
                │ Login Berhasil │     │    │              │
                └───────┬────────┘     │    │              │
                        │              │    │              │
        ┌───────────────┘              │    │              │
        │                              │    │              │
        ▼                              │    │              │
┌────────────────┐                    │    │              │
│ Dashboard User │                    │    │              │
└────┬───────────┘                    │    │              │
     │                                │    │              │
     └──────────────┐                 │    │              │
                    │                 │    │              │
                    ▼                 │    │              │
            ┌──────────────┐          │    │              │
            │ Lihat Lelang │          │    │              │
            └──────┬───────┘          │    │              │
                   │                  │    │              │
                   ▼                  │    │              │
          ┌──────────────────┐        │    │              │
          │ Pilih Item       │◄───────┘    │              │
          │ Lelang?          │             │              │
          └──┬───────────┬───┘             │              │
             │           │                 │              │
         Tidak           Ya                │              │
             │           │                 │              │
             │           ▼                 │              │
             │   ┌────────────────┐        │              │
             │   │ Detail Lelang  │        │              │
             │   │ Item Terpilih  │        │              │
             │   └────────┬───────┘        │              │
             │            │                │              │
             │            ▼                │              │
             │   ┌────────────────┐        │              │
             │   │ Item Lelang    │        │              │
             │   │ Dipilih        │        │              │
             │   └────────┬───────┘        │              │
             │            │                │              │
             │            ▼                │              │
             │   ┌────────────────┐        │              │
             │   │ Ingin Bidding  │        │              │
             │   │ Item Ini?      │        │              │
             │   └──┬─────────┬───┘        │              │
             │      │         │            │              │
             │  Tidak         Ya           │              │
             │      │         │            │              │
             └──────┘         ▼            │              │
                      ┌────────────┐       │              │
                      │ Cek Saldo  │       │              │
                      └──┬─────┬───┘       │              │
                         │     │           │              │
                  Saldo  │     │  Saldo    │              │
                  Cukup  │     │  Tidak    │              │
                         │     │  Cukup    │              │
                         │     │           │              │
                         │     ▼           │              │
                         │  ┌──────────────────┐          │
                         │  │ Notifikasi:      │          │
                         │  │ Saldo Tidak Cukup│          │
                         │  └────────┬─────────┘          │
                         │           │                    │
                         │           ▼                    │
                         │  ┌──────────────┐◄─────────────┘
                         │  │ Halaman      │
                         │  │ Deposit      │
                         │  └──────┬───────┘
                         │         │
                         │         ▼
                         │  ┌──────────────────┐
                         │  │ Pilih Metode     │
                         │  │ Pembayaran       │
                         │  └────┬─────────────┘
                         │       │
                         │       ▼
                         │  ┌──────────────┐
                         │  │ Metode       │
                         │  │ Pembayaran?  │
                         │  └┬──┬───┬───┬──┘
                         │   │  │   │   │
                         │   ▼  ▼   ▼   ▼
                         │  ┌───────────────────────────┐
                         │  │ Transfer Bank / VA /      │
                         │  │ E-Wallet / Kartu Kredit   │
                         │  └────────────┬──────────────┘
                         │               │
                         │               ▼
                         │      ┌─────────────────┐
                         │      │ Lakukan         │
                         │      │ Pembayaran      │
                         │      └────────┬────────┘
                         │               │
                         │               ▼
                         │      ┌─────────────────┐
                         │      │ Menunggu        │
                         │      │ Konfirmasi      │
                         │      └────────┬────────┘
                         │               │
                         │               ▼
                         │      ┌─────────────────┐
                         │      │ Saldo Terupdate │
                         │      └────────┬────────┘
                         │               │
                         │               ▼
                         │      ┌─────────────────┐
                         │      │ Kembali ke      │
                         │      │ Item Terpilih   │
                         │      └────────┬────────┘
                         │               │
                         └───────────────┘
                         │
                         ▼
                ┌────────────────────┐
                │ Input Nominal Bid  │
                │ untuk Item Terpilih│
                └──────────┬─────────┘
                           │
                           ▼
                ┌────────────────────┐
                │ Konfirmasi Bid     │
                │ untuk Item Ini     │
                └──────────┬─────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ Validasi Bid │
                    └──┬───────┬───┘
                       │       │
              Bid      │       │  Bid
              Terlalu  │       │  Valid
              Rendah   │       │
                       │       │
                       ▼       ▼
            ┌──────────────┐  ┌─────────────────┐
            │ Error: Bid   │  │ Proses Bidding  │
            │ Minimal Tidak│  │ Item Terpilih   │
            │ Terpenuhi    │  └────────┬────────┘
            └──────┬───────┘           │
                   │                   ▼
                   │          ┌─────────────────┐
                   │          │ Potong Saldo    │
                   │          │ Sesuai Bid      │
                   │          └────────┬────────┘
                   │                   │
                   │                   ▼
                   │          ┌─────────────────┐
                   │          │ Bidding Berhasil│
                   │          │ untuk Item Ini! │
                   │          └────────┬────────┘
                   │                   │
                   └───────────────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │ Status Lelang   │
                              │ Item Ini?       │
                              └──┬──────────┬───┘
                                 │          │
                         Masih   │          │  Lelang
                      Berlangsung│          │  Selesai
                                 │          │
                                 ▼          ▼
                      ┌──────────────────┐  ┌──────────────┐
                      │ Monitor Status   │  │ Pemenang     │
                      │ Bid Item Terpilih│  │ Item Ini?    │
                      └────────┬─────────┘  └──┬───────┬───┘
                               │               │       │
                               ▼               │       │
                      ┌──────────────────┐  Anda     Anda
                      │ Status Bid Anda  │  Menang   Kalah
                      │ untuk Item Ini?  │     │       │
                      └──┬───────────┬───┘     │       │
                         │           │         │       │
                  Terkalahkan  Masih │         │       │
                         │     Tertinggi       │       │
                         │           │         │       │
                         ▼           ▼         ▼       ▼
            ┌───────────────────┐ ┌─────────┐ ┌───────────────┐ ┌─────────────┐
            │ Notifikasi: Bid   │ │Notifikasi│ │Notifikasi     │ │Notifikasi & │
            │ Anda Telah        │ │Anda Masih│ │Pemenang Item  │ │Refund Saldo │
            │ Dilampaui         │ │Tertinggi │ │Terpilih &     │ └──────┬──────┘
            └─────────┬─────────┘ └────┬─────┘ │Proses Pembayar│        │
                      │                │       └──────┬────────┘        │
                      ▼                │              │                 │
              ┌───────────────┐        │              ▼                 │
              │ Bid Lagi      │        │     ┌─────────────────┐       │
              │ Item Ini?     │        │     │ Pembayaran Penuh│       │
              └──┬────────┬───┘        │     │ Item Menang     │       │
                 │        │            │     └────────┬────────┘       │
             Tidak        Ya           │              │                │
                 │        │            │              ▼                │
                 ▼        │            │     ┌─────────────────┐       │
         ┌───────────────┐│            │     │ Proses          │       │
         │ Lihat Item    ││            │     │ Pengiriman Item │       │
         │ Lain?         ││            │     └────────┬────────┘       │
         └──┬────────┬───┘│            │              │                │
            │        │    │            │              ▼                │
        Tidak        Ya   │            │     ┌─────────────────┐       │
            │        │    │            │     │ Transaksi       │       │
            ▼        │    │            │     │ Selesai         │       │
        ┌────────┐   │    │            │     └─────────────────┘       │
        │Selesai │   │    │            │                               │
        └────────┘   │    │            ▼                               │
                     │    │     ┌──────────────┐                       │
                     │    │     │ Menunggu     │                       │
                     │    │     │ Lelang Selesai│                      │
                     │    │     └──────────────┘                       │
                     │    │                                            │
                     │    └────────────────────────────────────────────┘
                     │
                     └──────────────────────────────────────────────────┐
                                                                        │
                                                                        ▼
                                                               ┌────────────────┐
                                                               │ Lihat Lelang   │
                                                               └────────────────┘


KETERANGAN WARNA (dalam implementasi):
- Hijau Muda: Start & Transaksi Selesai
- Merah Muda: Proses Selesai/Keluar
- Hijau: Bidding Berhasil, Login Berhasil, Saldo Terupdate
- Emas: Notifikasi Pemenang
- Biru Muda: Item Lelang Dipilih
- Kuning: Node Pilihan Item
```

**KEY FEATURES:**
1. ✅ User HARUS memilih item lelang tertentu sebelum bisa bidding
2. ✅ Jelas menunjukkan "Item Terpilih" di setiap tahap bidding
3. ✅ Flow lengkap dari guest mode hingga transaksi selesai
4. ✅ Proses deposit dan top up saldo terintegrasi
5. ✅ Monitoring real-time status bid untuk item yang dipilih
6. ✅ Notifikasi menang/kalah dengan proses pembayaran

// Database Schema untuk Sistem Lelang Online
// Paste kode ini ke dbdiagram.io

Table users {
  user_id int [pk, increment]
  email varchar(255) [unique, not null]
  password_hash varchar(255) [not null]
  full_name varchar(255) [not null]
  phone varchar(20)
  id_card_number varchar(50) [unique]
  id_card_type enum('KTP', 'SIM', 'PASSPORT')
  address text
  city varchar(100)
  province varchar(100)
  postal_code varchar(10)
  is_verified boolean [default: false]
  verification_token varchar(255)
  verification_date timestamp
  balance decimal(15,2) [default: 0]
  status enum('active', 'suspended', 'blocked') [default: 'active']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    email
    phone
    id_card_number
  }
}

Table sellers {
  seller_id int [pk, increment]
  seller_name varchar(255) [not null]
  seller_type enum('bank', 'government', 'company', 'individual')
  address text
  phone varchar(20)
  email varchar(255)
  contact_person varchar(255)
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table organizers {
  organizer_id int [pk, increment]
  organizer_name varchar(255) [not null]
  organizer_code varchar(50) [unique]
  organizer_type enum('KPKNL', 'bank', 'private')
  address text
  city varchar(100)
  province varchar(100)
  phone varchar(20)
  email varchar(255)
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table item_categories {
  category_id int [pk, increment]
  category_name varchar(100) [not null]
  parent_category_id int
  description text
  created_at timestamp [default: `now()`]
  
  indexes {
    parent_category_id
  }
}

Table auction_items {
  item_id int [pk, increment]
  lot_code varchar(50) [unique, not null]
  item_name varchar(255) [not null]
  category_id int [not null, ref: > item_categories.category_id]
  seller_id int [not null, ref: > sellers.seller_id]
  organizer_id int [not null, ref: > organizers.organizer_id]
  item_type enum('movable', 'immovable') [not null]
  sub_type varchar(100)
  description text
  detailed_description text
  ownership_proof varchar(255)
  ownership_number varchar(255)
  ownership_date date
  ownership_holder_name varchar(255)
  limit_price decimal(15,2) [not null]
  deposit_amount decimal(15,2) [not null]
  starting_price decimal(15,2)
  current_highest_bid decimal(15,2)
  increment_amount decimal(15,2)
  auction_method enum('open_bidding', 'closed_bidding', 'tender')
  status enum('draft', 'published', 'ongoing', 'closed', 'cancelled') [default: 'draft']
  view_count int [default: 0]
  bid_count int [default: 0]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    lot_code
    category_id
    seller_id
    organizer_id
    status
  }
}

Table item_properties {
  property_id int [pk, increment]
  item_id int [not null, ref: > auction_items.item_id]
  address text
  village varchar(100)
  district varchar(100)
  city varchar(100)
  province varchar(100)
  land_area decimal(10,2)
  building_area decimal(10,2)
  land_area_unit varchar(10) [default: 'M2']
  building_area_unit varchar(10) [default: 'M2']
  latitude decimal(10,8)
  longitude decimal(11,8)
  year_built int
  certificate_type varchar(50)
  certificate_number varchar(100)
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    item_id
  }
}

Table item_images {
  image_id int [pk, increment]
  item_id int [not null, ref: > auction_items.item_id]
  image_url varchar(500) [not null]
  image_type enum('main', 'gallery', 'document')
  display_order int [default: 0]
  caption text
  created_at timestamp [default: `now()`]
  
  indexes {
    item_id
    image_type
  }
}

Table auction_schedules {
  schedule_id int [pk, increment]
  item_id int [not null, ref: > auction_items.item_id]
  registration_start timestamp
  registration_end timestamp
  deposit_deadline timestamp [not null]
  auction_start timestamp [not null]
  auction_end timestamp [not null]
  announcement_date timestamp
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    item_id
    auction_start
    auction_end
  }
}

Table bids {
  bid_id int [pk, increment]
  item_id int [not null, ref: > auction_items.item_id]
  user_id int [not null, ref: > users.user_id]
  bid_amount decimal(15,2) [not null]
  bid_type enum('manual', 'auto', 'proxy')
  bid_status enum('active', 'outbid', 'winning', 'won', 'lost', 'cancelled') [default: 'active']
  is_highest boolean [default: false]
  bid_time timestamp [default: `now()`]
  ip_address varchar(45)
  user_agent text
  
  indexes {
    item_id
    user_id
    bid_status
    bid_time
    (item_id, bid_amount)
  }
}

Table deposits {
  deposit_id int [pk, increment]
  user_id int [not null, ref: > users.user_id]
  item_id int [ref: > auction_items.item_id]
  deposit_amount decimal(15,2) [not null]
  deposit_type enum('auction_deposit', 'balance_topup')
  payment_method enum('bank_transfer', 'virtual_account', 'ewallet', 'credit_card')
  payment_proof varchar(500)
  transaction_id varchar(100) [unique]
  bank_name varchar(100)
  account_number varchar(50)
  account_name varchar(255)
  status enum('pending', 'verified', 'rejected', 'refunded') [default: 'pending']
  verified_by int [ref: > users.user_id]
  verified_at timestamp
  notes text
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    user_id
    item_id
    transaction_id
    status
  }
}

Table transactions {
  transaction_id int [pk, increment]
  user_id int [not null, ref: > users.user_id]
  transaction_type enum('deposit', 'bid', 'refund', 'payment', 'withdrawal')
  amount decimal(15,2) [not null]
  balance_before decimal(15,2)
  balance_after decimal(15,2)
  reference_type enum('deposit', 'bid', 'auction', 'refund')
  reference_id int
  description text
  status enum('pending', 'completed', 'failed') [default: 'pending']
  created_at timestamp [default: `now()`]
  
  indexes {
    user_id
    transaction_type
    reference_id
    created_at
  }
}

Table auction_winners {
  winner_id int [pk, increment]
  item_id int [not null, ref: > auction_items.item_id]
  user_id int [not null, ref: > users.user_id]
  bid_id int [not null, ref: > bids.bid_id]
  winning_amount decimal(15,2) [not null]
  payment_deadline timestamp
  payment_status enum('unpaid', 'partial', 'paid', 'overdue') [default: 'unpaid']
  payment_date timestamp
  payment_proof varchar(500)
  delivery_status enum('pending', 'processing', 'shipped', 'delivered', 'completed')
  delivery_address text
  tracking_number varchar(100)
  notes text
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    item_id
    user_id
    payment_status
    delivery_status
  }
}

Table notifications {
  notification_id int [pk, increment]
  user_id int [not null, ref: > users.user_id]
  notification_type enum('bid_outbid', 'bid_winning', 'auction_ending', 'auction_won', 'auction_lost', 'payment_reminder', 'deposit_verified', 'general')
  title varchar(255) [not null]
  message text [not null]
  related_item_id int [ref: > auction_items.item_id]
  related_bid_id int [ref: > bids.bid_id]
  is_read boolean [default: false]
  read_at timestamp
  created_at timestamp [default: `now()`]
  
  indexes {
    user_id
    is_read
    notification_type
    created_at
  }
}

Table user_favorites {
  favorite_id int [pk, increment]
  user_id int [not null, ref: > users.user_id]
  item_id int [not null, ref: > auction_items.item_id]
  created_at timestamp [default: `now()`]
  
  indexes {
    (user_id, item_id) [unique]
    user_id
    item_id
  }
}

Table activity_logs {
  log_id int [pk, increment]
  user_id int [ref: > users.user_id]
  activity_type enum('login', 'logout', 'view_item', 'bid', 'deposit', 'profile_update')
  description text
  ip_address varchar(45)
  user_agent text
  created_at timestamp [default: `now()`]
  
  indexes {
    user_id
    activity_type
    created_at
  }
}

Table admin_users {
  admin_id int [pk, increment]
  username varchar(100) [unique, not null]
  password_hash varchar(255) [not null]
  full_name varchar(255) [not null]
  email varchar(255) [unique, not null]
  role enum('super_admin', 'admin', 'verifier', 'moderator')
  permissions json
  is_active boolean [default: true]
  last_login timestamp
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  
  indexes {
    username
    email
    role
  }
}

// Relationships are defined using ref in column definitions above


graph TD
    Start([User Membuka Aplikasi]) --> CheckAuth{Status Login?}
    
    CheckAuth -->|Belum Login| GuestMode[Mode Guest/Tamu]
    CheckAuth -->|Sudah Login| UserDashboard[Dashboard User]
    
    GuestMode --> BrowseAuction[Lihat Daftar Lelang]
    BrowseAuction --> SelectItem{Pilih Item Lelang?}
    
    SelectItem -->|Tidak| BrowseAuction
    SelectItem -->|Ya| ViewDetail[Lihat Detail Lelang<br/>Item Terpilih]
    
    ViewDetail --> WantToBid{Ingin Ikut Bidding?}
    WantToBid -->|Tidak| BrowseAuction
    WantToBid -->|Ya| MustLogin[Diarahkan ke Halaman Login]
    
    MustLogin --> LoginPage[Halaman Login/Register]
    LoginPage --> Register{Sudah Punya Akun?}
    
    Register -->|Belum| RegisterForm[Form Registrasi]
    RegisterForm --> VerifyAccount[Verifikasi Akun<br/>Email/Phone]
    VerifyAccount --> LoginSuccess
    
    Register -->|Sudah| LoginForm[Form Login]
    LoginForm --> LoginSuccess[Login Berhasil]
    
    LoginSuccess --> UserDashboard
    
    UserDashboard --> MenuChoice{Pilih Menu}
    MenuChoice --> BrowseAuction2[Lihat Lelang]
    MenuChoice --> Profile[Profil & Saldo]
    MenuChoice --> History[Riwayat Bidding]
    
    BrowseAuction2 --> SelectItem2{Pilih Item Lelang?}
    SelectItem2 -->|Tidak| BrowseAuction2
    SelectItem2 -->|Ya| ViewDetail2[Detail Lelang<br/>Item Terpilih]
    
    ViewDetail2 --> ItemSelected[Item Lelang Dipilih]
    ItemSelected --> ReadyToBid{Ingin Bidding<br/>Item Ini?}
    
    ReadyToBid -->|Tidak| BrowseAuction2
    ReadyToBid -->|Ya| CheckBalance{Cek Saldo}
    
    CheckBalance -->|Saldo Cukup| PlaceBid[Input Nominal Bid<br/>untuk Item Terpilih]
    CheckBalance -->|Saldo Tidak Cukup| NeedDeposit[Notifikasi:<br/>Saldo Tidak Cukup]
    
    NeedDeposit --> DepositPage[Halaman Deposit]
    DepositPage --> ChoosePayment[Pilih Metode Pembayaran]
    ChoosePayment --> PaymentMethod{Metode Pembayaran}
    
    PaymentMethod --> BankTransfer[Transfer Bank]
    PaymentMethod --> VirtualAccount[Virtual Account]
    PaymentMethod --> EWallet[E-Wallet]
    PaymentMethod --> CreditCard[Kartu Kredit]
    
    BankTransfer --> MakePayment[Lakukan Pembayaran]
    VirtualAccount --> MakePayment
    EWallet --> MakePayment
    CreditCard --> MakePayment
    
    MakePayment --> WaitConfirm[Menunggu Konfirmasi]
    WaitConfirm --> BalanceUpdated[Saldo Terupdate]
    
    BalanceUpdated --> ReturnToItem[Kembali ke<br/>Item Terpilih]
    ReturnToItem --> CheckBalance
    
    PlaceBid --> ConfirmBid[Konfirmasi Bid<br/>untuk Item Ini]
    ConfirmBid --> ValidateBid{Validasi Bid}
    
    ValidateBid -->|Bid Terlalu Rendah| BidError[Error: Bid Minimal<br/>Tidak Terpenuhi]
    BidError --> PlaceBid
    
    ValidateBid -->|Bid Valid| ProcessBid[Proses Bidding<br/>Item Terpilih]
    ProcessBid --> DeductBalance[Potong Saldo<br/>Sesuai Bid]
    DeductBalance --> BidSuccess[Bidding Berhasil<br/>untuk Item Ini!]
    
    BidSuccess --> WaitAuction{Status Lelang<br/>Item Ini}
    WaitAuction -->|Masih Berlangsung| MonitorBid[Monitor Status Bid<br/>Item Terpilih]
    
    MonitorBid --> BidStatus{Status Bid Anda<br/>untuk Item Ini}
    BidStatus -->|Terkalahkan| Notification1[Notifikasi: Bid Anda<br/>Telah Dilampaui]
    BidStatus -->|Masih Tertinggi| Notification2[Notifikasi: Anda<br/>Masih Tertinggi]
    
    Notification1 --> WantRebid{Bid Lagi<br/>Item Ini?}
    WantRebid -->|Ya| PlaceBid
    WantRebid -->|Tidak| BrowseOther{Lihat Item Lain?}
    BrowseOther -->|Ya| BrowseAuction2
    BrowseOther -->|Tidak| End1([Selesai])
    
    Notification2 --> End2([Menunggu Lelang Selesai])
    
    WaitAuction -->|Lelang Selesai| AuctionEnd{Pemenang<br/>Item Ini?}
    AuctionEnd -->|Anda Menang| WinNotif[Notifikasi Pemenang<br/>Item Terpilih<br/>& Proses Pembayaran]
    AuctionEnd -->|Anda Kalah| LoseNotif[Notifikasi & Refund Saldo]
    
    WinNotif --> PaymentFull[Pembayaran Penuh<br/>Item Menang]
    PaymentFull --> DeliveryProcess[Proses Pengiriman<br/>Item]
    DeliveryProcess --> End3([Transaksi Selesai])
    
    LoseNotif --> End4([Selesai])
    
    Profile --> ViewBalance[Lihat Saldo]
    ViewBalance --> TopUp{Top Up?}
    TopUp -->|Ya| DepositPage
    TopUp -->|Tidak| UserDashboard
    
    History --> ViewHistory[Lihat Riwayat<br/>Item yang Pernah di-Bid]
    ViewHistory --> UserDashboard
    
    style Start fill:#e1f5e1
    style End1 fill:#ffe1e1
    style End2 fill:#ffe1e1
    style End3 fill:#e1f5e1
    style End4 fill:#ffe1e1
    style BidSuccess fill:#90ee90
    style WinNotif fill:#ffd700
    style LoginSuccess fill:#90ee90
    style BalanceUpdated fill:#90ee90
    style ItemSelected fill:#87ceeb
    style SelectItem fill:#ffeb99
    style SelectItem2 fill:#ffeb99