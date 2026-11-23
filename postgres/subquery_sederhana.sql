-- Topik 3: Subquery sederhana

-- Konsep: Subquery adalah `SELECT` di dalam query lain. Ada tiga bentuk yang sering dipakai: scalar (mengembalikan satu nilai, dipakai di perbandingan), list (`IN (...)`), dan `EXISTS` (hanya mengecek ada atau tidaknya baris). Subquery correlated merujuk kolom query luar, jadi secara logis dievaluasi per baris luar.

-- Contoh:
-- ```go
-- const q = `
-- SELECT u.id, u.name
-- FROM users u
-- WHERE EXISTS (
--   SELECT 1 FROM orders o
--   WHERE o.user_id = u.id AND o.total > $1)`
-- ```

-- Latihan (skema: `users(id, name)`, `orders(id, user_id, total NUMERIC)`):

-- 1. Tampilkan order yang `total`-nya di atas rata-rata semua order (scalar subquery).

SELECT id, user_id, total
FROM orders
WHERE total > (
    SELECT AVG(total) FROM orders
);

-- 2. Tampilkan user yang belum punya order memakai `NOT EXISTS`.

SELECT u.id, u.name
FROM users AS u
WHERE NOT EXISTS (
    SELECT 1 FROM orders AS o ON u.id = o.user_id
)

-- 3. Kalau `orders.user_id` bisa NULL, apa hasil `WHERE u.id NOT IN (SELECT user_id FROM orders)` dan kenapa? Prediksi dulu sebelum menjawab.

-- hasilnya selalu kosong, karena saat mencocokan data di in akan membadingkan dengan null yang hasilnya unknow, dan where hanya mengembalikan yang true
