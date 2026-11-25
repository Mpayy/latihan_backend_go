-- Topik 4: Pagination (offset & keyset)

-- Konsep: Offset memakai LIMIT n OFFSET m: sederhana, tapi database tetap membaca lalu membuang m baris, jadi makin lambat di halaman dalam, dan data bisa duplikat atau terlewat kalau ada insert/delete antar halaman. Keyset memakai "mulai setelah baris terakhir yang dilihat": stabil dan cepat karena bisa memakai index.

-- Contoh keyset:

-- go
-- const q = `
-- SELECT id, total, created_at
-- FROM orders
-- WHERE (created_at, id) < ($1, $2)
-- ORDER BY created_at DESC, id DESC
-- LIMIT $3`

-- Latihan (skema: orders(id, user_id, total, created_at)):

-- Tulis query offset untuk halaman ke-3, 20 baris per halaman, order terbaru dulu. Kenapa ORDER BY wajib ada?

SELECT id, user_id, total, created_at
FROM orders
ORDER BY created_at DESC, id DESC
LIMIT 20
OFFSET 40;
-- Tanpa ORDER BY, SQL tidak menjamin data terurut, jadi halaman bisa saling tumpang tindih atau ada baris yang terlewat.

-- Kenapa keyset di atas memakai (created_at, id) dan bukan created_at saja?

-- karena created_at tidak dijamin unik. Kalau cursor hanya created_at < $1, semua baris lain dengan timestamp yang sama ikut terlewat di halaman berikutnya. id (unik) membuat posisi cursor menjadi tunggal dan pasti. Syarat penting: kolom di WHERE harus sama dan searah dengan ORDER BY. 

-- Sebutkan satu kekurangan keyset dibanding offset.
-- keyset tidak bisa langsung melompat kehalaman atau page yang ingin di tuju, karena keyset bergantung pada posisi cursor. sort key harus deterministik, dan mengubah urutan sort berarti mengubah query cursor.


-- Quiz: Pagination (offset & keyset)

-- Skema: orders(id, user_id, total, status, created_at).

-- 1. Tulis query keyset untuk halaman berikutnya dari daftar order milik user_id = $1 yang status = 'paid', urut created_at DESC, id DESC, 20 baris. Cursor: (created_at, id) baris terakhir halaman sebelumnya. Lalu tulis index yang paling cocok.

SELECT id, user_id, total, status, created_at
FROM orders
WHERE user_id = $1 AND status = 'paid' AND (created_at < $2 OR (created_at = $2 AND id < $3))
ORDER BY created_at DESC, id DESC
LIMIT 20;

-- dengan perbandingan tuple
SELECT id, user_id, total, status, created_at
FROM orders
WHERE user_id = $1 AND status = 'paid' AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT 20;

CREATE INDEX idx_orders_user_status_created ON orders(user_id, status, created_at DESC, id DESC);

-- 2. Seorang teman memakai WHERE created_at < $1 ORDER BY created_at DESC, id DESC LIMIT 20 dengan cursor hanya created_at. Skenario konkret apa yang membuat datanya hilang? Beri contoh angka.
-- data akan hilang saat skenario jumlah data yang created_atnya sama lebih banyak dari limit. jadi saat berganti halaman, saat cursor mengambil created_at di baris terakhir halaman sebelumnnya. query selanjutnya hanya mengambil created at di bawahnya, bukan sisa yang tadi kepotong.

-- contoh angka
-- Baris 1 s.d. 18: id = 100 sampai id = 83 (dibuat pukul 10:05:00 – 10:01:00) 
-- Baris 19: id = 82, created_at = 10:00:00
-- Baris 20: id = 81, created_at = 10:00:00  -> Batas akhir Halaman 1 (LIMIT 20)
-- Baris 21: id = 80, created_at = 10:00:00
-- Baris 22: id = 79, created_at = 10:00:00
-- Baris 23: id = 78, created_at = 09:59:00
-- saat kehalaman berikutnya, query akan mengambil data yang created_atnya < 10:00:00, maka id 80 dan 90 ga di ambil lagi dan akan hilang

-- 3. Tabel punya 5 juta baris. Jelaskan kenapa OFFSET 4000000 lambat, padahal LIMIT-nya hanya 20. Kapan kamu tetap memilih offset?
-- karena database akan mengambil 4000020 baris, lalu membuang 4000000 dan hanya menampilkan 20 baris terakhir. memilih offset hanya jika jumlah data sedikit (atau halaman dibatasi hanya beberapa halaman awal), ketika fitur membutuhkan lompat langsung ke nomor halaman acak (jump to page N), dan saat data tidak sering berubah (insert/delete) ketika pengguna berpindah halaman.