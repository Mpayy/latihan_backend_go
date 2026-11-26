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