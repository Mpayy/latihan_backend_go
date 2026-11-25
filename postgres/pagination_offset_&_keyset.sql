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