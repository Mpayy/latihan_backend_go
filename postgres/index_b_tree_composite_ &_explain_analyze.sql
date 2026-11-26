-- Topik 6: Index (B-tree, composite) + EXPLAIN ANALYZE

-- Konsep: Index B-tree adalah struktur terurut yang membuat pencarian, ORDER BY, dan range tidak perlu membaca seluruh tabel (Seq Scan). Harganya: INSERT/UPDATE jadi lebih lambat dan memakai ruang. Composite index (a, b) hanya efisien kalau query memakai a (prefix kiri), bukan hanya b. EXPLAIN ANALYZE menjalankan query sungguhan dan menampilkan plan, waktu aktual, dan jumlah baris.

-- Contoh:

-- go
-- const q = `
-- EXPLAIN (ANALYZE, BUFFERS)
-- SELECT id, total FROM orders
-- WHERE user_id = $1 AND status = 'paid'
-- ORDER BY created_at DESC LIMIT 20`

-- Untuk query ini: CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at DESC);

-- Latihan (tabel orders berisi 1 juta baris):

-- Query WHERE user_id = 5 lambat dan plan menunjukkan Seq Scan. Tulis CREATE INDEX-nya. Setelah itu, apa yang kamu harapkan berubah di plan?

CREATE INDEX idx_orders_user_id ON orders (user_id);
-- yang saya harapkan query menggunakan index yang saya buat, dan merubah Seq Scan menjadi Index Scan, dan Execution Time turun drastis.

-- Dengan index (user_id, status), apakah WHERE status = 'paid' saja bisa memakainya secara efisien? Jelaskan alasannya.
-- Aturannya prefix kiri: tanpa user_id, Postgres tidak bisa menelusuri tree secara efisien. Secara teknis planner masih boleh memindai seluruh index (full index scan), tapi itu jarang lebih murah dari Seq Scan. Kalau query status saja memang sering dan selektif, buat index terpisah (atau partial index)

-- Di output EXPLAIN ANALYZE, dua angka apa yang kamu bandingkan untuk tahu planner salah memperkirakan?
-- estimasi rows dengan actual rows * loops berbeda jauh.

-- Quiz: Index (B-tree, composite) + EXPLAIN ANALYZE

-- Skema: orders(id, user_id, total, status, created_at), 1 juta baris.

-- 1. Index: (user_id, created_at). Mana dari query berikut yang bisa memakainya secara efisien, dan kenapa?

-- a. WHERE user_id = 5 ORDER BY created_at DESC LIMIT 10
-- b. WHERE created_at > '2026-01-01'
-- c. WHERE user_id = 5 AND created_at > '2026-01-01'
-- d. WHERE user_id > 5 ORDER BY created_at

-- yang bisa ,memakainya secara efisien adalah query a, dan c
-- a karena user_id sudah di lock dengan sama dengan(=), maka created_at nya boleh di taruh di ORDER BY
-- c karena memakai semua kolom di index ada di WHERE

-- kalo b dia ga sesuai aturan prefix, kalo d bisa tapi tidak efisien, karena hanya bisa di gunakan untuk user > 5, tapi index tersebut tidak bisa menghilangkan langkah Sort untuk ORDER BY created_at

-- 2. Query WHERE LOWER(email) = 'a@x.com' tetap Seq Scan padahal ada index biasa di email. Kenapa, dan apa solusinya?
-- karena data email di index bukan lowercase, jadi data base tidak memakai index tersebut, solusinya adalah membuat email di index dengan LOWER(email)

-- 3. Potongan plan:

-- Index Scan using idx_orders_user_id on orders (cost=0.4..8.5 rows=1 ...) (actual time=0.02..45.3 rows=48000 loops=1)

-- Apa yang janggal, apa kemungkinan penyebabnya, dan apa langkah pertamamu? Lalu: kapan planner justru memilih Seq Scan padahal index ada, dan kenapa itu bukan selalu masalah?

-- estimasi rows meleset jauh dari actual rows * loops, kemungkinan penyebabnya adalah statistik tabel sudah belum di update atau baru menjalankan EXPLAIN ANALYZE setelah banyak perubahan data. planner akan tetap memilih seq scan walau bisa mengunakan index scan, saat ukuran table masih kecil, kardinalitas pada kolom tersebut rendah, dan saat table masih sedikit