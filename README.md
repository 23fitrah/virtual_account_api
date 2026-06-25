# Virtual Account Payment System — Project Documentation

## 1. Project Overview

Project ini bertujuan untuk membangun sistem **Virtual Account (VA)** yang memungkinkan nasabah atau pelanggan melakukan pembayaran melalui nomor rekening unik yang di-generate khusus untuk setiap transaksi atau setiap pelanggan. Sistem ini menjadi jembatan antara aplikasi merchant/bisnis dengan pihak bank atau payment gateway, sehingga proses rekonsiliasi pembayaran dapat dilakukan secara otomatis dan real-time.

Dokumentasi ini disusun dari sudut pandang System Analyst, mencakup business requirement, business flow, system flow, hingga kebutuhan fungsional dan non-fungsional yang menjadi acuan bagi tim development.

## 2. Business Background

Saat ini proses pembayaran manual (transfer bank biasa) menyulitkan proses rekonsiliasi karena petugas finance harus mencocokkan satu per satu mutasi rekening dengan invoice yang ada. Hal ini menyebabkan:

- Proses konfirmasi pembayaran lambat (manual checking)
- Risiko human error dalam pencocokan data
- Customer experience kurang baik karena status order tidak update otomatis

Dengan Virtual Account, setiap transaksi/customer mendapatkan nomor rekening unik, sehingga sistem dapat melakukan auto-reconciliation begitu pembayaran masuk.

## 3. Objectives

- Menyediakan nomor Virtual Account unik untuk setiap transaksi atau pelanggan
- Melakukan otomatisasi proses konfirmasi pembayaran (callback/notification dari bank)
- Mengintegrasikan sistem internal dengan API bank/payment gateway
- Menyediakan dashboard monitoring status pembayaran
- Memastikan keamanan data transaksi sesuai standar perbankan

## 4. Scope

**In Scope**
- Generate VA number (closed/open payment)
- Integrasi dengan bank/payment gateway (create VA, inquiry, payment notification)
- Callback handler untuk update status pembayaran
- Auto-expire VA number sesuai waktu yang ditentukan
- Reporting & reconciliation

**Out of Scope**
- Proses refund manual ke rekening sumber
- Integrasi dengan metode pembayaran selain VA (e-wallet, QRIS, dll) — akan dibahas di project terpisah

## 5. Actors

| Actor | Deskripsi |
|---|---|
| Customer | Pihak yang melakukan pembayaran melalui VA |
| Merchant/Internal System | Sistem yang melakukan request pembuatan VA |
| Bank/Payment Gateway | Pihak ketiga penyedia layanan VA |
| Finance/Admin | Melakukan monitoring dan rekonsiliasi transaksi |

## 6. Business Flow

1. Customer melakukan order/transaksi pada sistem merchant
2. Sistem internal mengirimkan request pembuatan VA ke Bank/Payment Gateway
3. Bank mengembalikan nomor VA beserta masa berlakunya
4. Customer melakukan pembayaran ke nomor VA tersebut melalui ATM/mobile banking
5. Bank mengirimkan notifikasi pembayaran (callback) ke sistem internal
6. Sistem internal memverifikasi dan mengupdate status transaksi menjadi "Paid"
7. Sistem mengirimkan notifikasi ke customer (email/SMS/push notification)

## 7. System Flow / Sequence

```
Customer -> Internal System : Create Order
Internal System -> Bank API : Request Create VA
Bank API -> Internal System : Response (VA Number, Expired Date)
Internal System -> Customer : Display VA Number
Customer -> Bank Channel : Pay via VA
Bank -> Internal System : Callback Notification (Payment Success)
Internal System -> Internal System : Update Order Status
Internal System -> Customer : Send Payment Confirmation
```

## 8. Functional Requirements

| ID | Requirement | Priority |
|---|---|---|
| FR-01 | Sistem dapat melakukan generate VA number otomatis | High |
| FR-02 | Sistem dapat menerima callback notification dari bank | High |
| FR-03 | Sistem dapat melakukan validasi signature/checksum pada callback | High |
| FR-04 | Sistem dapat menampilkan status transaksi (Pending/Paid/Expired) | High |
| FR-05 | Sistem dapat melakukan auto-expire VA sesuai waktu yang ditentukan | Medium |
| FR-06 | Sistem dapat melakukan retry mechanism jika callback gagal diproses | Medium |
| FR-07 | Sistem menyediakan reporting transaksi harian | Low |

## 9. Non-Functional Requirements

| ID | Requirement |
|---|---|
| NFR-01 | Response time API create VA maksimal 2 detik |
| NFR-02 | Sistem harus mendukung high availability (uptime ≥ 99.9%) |
| NFR-03 | Callback dari bank harus diproses secara idempotent (anti duplicate processing) |
| NFR-04 | Semua komunikasi API harus menggunakan HTTPS/TLS |
| NFR-05 | Data transaksi sensitif harus terenkripsi (at-rest & in-transit) |
| NFR-06 | Sistem harus memiliki audit trail untuk setiap perubahan status transaksi |

## 10. Data Model (Simplified ERD)

**Table: virtual_accounts**
- id
- order_id (FK)
- va_number
- bank_code
- amount
- status (PENDING, PAID, EXPIRED) x
- expired_at
- created_at

**Table: payment_callbacks**
- id
- va_id (FK)
- raw_payload
- signature
- processed_at
- status

**Table: orders**
- id
- customer_id
- total_amount
- status
- created_at

## 11. API List (Contoh)

| Method | Endpoint | Deskripsi |
|---|---|---|
| POST | /api/v1/virtual-accounts | Membuat VA baru untuk sebuah order |
| GET | /api/v1/virtual-accounts/{id} | Mengecek status VA |
| POST | /api/v1/callback/payment | Endpoint untuk menerima notifikasi pembayaran dari bank |
| GET | /api/v1/transactions | Mendapatkan daftar transaksi (untuk reporting) |

## 12. Assumptions & Constraints

- Pihak bank/payment gateway sudah menyediakan API dokumentasi resmi (sandbox & production)
- Format dan struktur callback payload mengikuti standar dari masing-masing bank
- Proses settlement dana ke rekening merchant dilakukan oleh pihak bank, di luar tanggung jawab sistem internal

## 13. Tech Stack (Contoh)

- Backend: Golang / C# .NET / Laravel-Lumen
- Database: PostgreSQL / MySQL
- Message Broker: RabbitMQ / Kafka (untuk async callback processing)
- Containerization: Docker
- CI/CD: GitLab CI / GitHub Actions
- Monitoring: Prometheus + Grafana

## 14. Glossary

| Term | Deskripsi |
|---|---|
| VA | Virtual Account, nomor rekening unik untuk menerima pembayaran |
| Callback | Notifikasi yang dikirim bank ke sistem internal saat pembayaran terjadi |
| Reconciliation | Proses pencocokan data transaksi dengan mutasi rekening |
| Idempotent | Proses yang aman dijalankan berulang kali tanpa menimbulkan efek ganda |
