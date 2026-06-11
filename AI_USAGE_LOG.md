# AI Usage Log

## 1. Tool yang Dipakai
* **Nama Tool:** Gemini
* **Versi:** Gemini (Google DeepMind - *Chat UI/Web Interface*)

## 2. Pattern Penggunaan
* **Architectural & Scaffolding (Boilerplate):** Digunakan pada jam pertama untuk merumuskan struktur *Clean Architecture* di Golang dan memastikan konfigurasi `docker-compose.yml` terisolasi dengan aman untuk *deployment* (Monorepo).
* **Debugging & Error Handling:** Menganalisis *stack trace error* spesifik, seperti menyelesaikan kendala `nil pointer dereference` saat inisiasi *unit test* di Golang dengan mengimplementasikan `go-sqlmock`.
* **Refactoring & Code Simplification:** Digunakan sebagai *reviewer* kode untuk menyederhanakan bahasa teknis dan alur logika agar kode lebih *maintainable* dan siap untuk *Live Defense*.

## 3. History Prompt Utama (Selama Fase 2)
Berikut adalah daftar *prompt* krusial yang mengarahkan pengembangan inti sistem:
1. *"Apa saja rancangan tabel database yang optimal untuk Sistem Pemantauan Log dengan arsitektur meminimalisir bottleneck, mempertimbangkan kita memakai MySQL InnoDB?"*
2. *"Tolong evaluasi konfigurasi `docker-compose.yml` saya, pastikan variabel environment DB_MYSQL_NAME sudah selaras dengan Go backend dan healthcheck-nya berjalan sempurna tanpa error spasi pada password."*
3. *"Saya mendapat error `panic: runtime error: invalid memory address` saat menjalankan `go test -v ./internal/usecase/` karena `*Tx` nil. Bagaimana cara setup mock database yang benar untuk skenario Pessimistic Locking di Go?"*
4. *"Mas, sepertinya gunakan bahasa sederhana saja, jangan pakai istilah 'concurrent skala besar tanpa bottleneck' kalau saya sendiri tidak mengerti."* (Menolak saran *copywriting* teknis AI).

## 4. Penilaian Judgment & Keputusan Kritis (*Impactful Interaction*)
Sesuai dengan prinsip *engineering* modern, asisten AI digunakan sebagai alat bantu, bukan pengambil keputusan. Salah satu interaksi paling berdampak terjadi ketika AI menyarankan penggunaan istilah *buzzword* teknis *"memproses beban concurrent skala besar tanpa bottleneck"* untuk mendeskripsikan kapabilitas sistem Go Echo.

**Keputusan (Judgment):** Saya **menolak** penggunaan frasa tersebut dan menginstruksikan AI untuk merombaknya menjadi narasi operasional yang logis dan membumi (*"memproses banyak laporan yang masuk secara bersamaan tanpa membuat sistem menjadi lambat"*). 

**Alasan Penolakan:** Sebagai *Software Engineer*, saya berpegang pada prinsip kehati-hatian: saya tidak akan mengimplementasikan logika, baris kode, maupun narasi arsitektur yang tidak dapat saya bedah dan pertanggungjawabkan secara mandiri di hadapan juri saat *Live Defense*. Sistem ini dibangun untuk keandalan dan kesederhanaan, bukan untuk memaksakan kompleksitas semu. Penggunaan AI difokuskan murni untuk akselerasi penulisan *boilerplate* dan *unit test setup*, sementara logika inti seperti *Pessimistic Locking* dan *HTTP Short-polling* sepenuhnya berada di bawah kendali arsitektural saya.