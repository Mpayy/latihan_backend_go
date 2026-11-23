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


-- Quiz: Subquery sederhana

-- Skema: users(id, name), orders(id, user_id, total NUMERIC, status, created_at).

-- 1. Tampilkan user yang punya minimal satu order 'paid' dan tidak punya satu pun order 'cancelled'. Tulis dengan EXISTS dan NOT EXISTS. Jelaskan kenapa tidak perlu DISTINCT.

SELECT u.id, u.name
FROM users AS u
WHERE EXISTS (
    SELECT 1 FROM orders AS o WHERE u.id = o.user_id AND o.status = 'paid'
) AND NOT EXISTS (
    SELECT 1 FROM orders AS o WHERE u.id = o.user_id AND o.status = 'cancelled'
);

-- karena sumber baris berasal dari user, karena tidak memakai JOIN. dan WHERE EXISTS dan NOT EXISTS berhenti begitu menemukan baris yang cocok

-- 2. Tampilkan order milik tiap user yang total-nya di atas rata-rata order user itu sendiri. Subquery ini correlated. Tulis query-nya, lalu jelaskan apa bedanya dengan subquery tanpa korelasi seperti total > (SELECT AVG(total) FROM orders).

SELECT o.id, o.user_id, o.total
FROM orders AS o
WHERE o.total > (
    SELECT AVG(o2.total) FROM orders AS o2 WHERE o2.user_id = o.user_id
);

-- bedanya adalah query tersebut mengambil rata rata dari semua user bukan per user

-- 3. Prediksi dulu, lalu jelaskan: kalau orders.user_id tidak punya NULL, apakah WHERE u.id IN (SELECT user_id FROM orders) dan WHERE EXISTS (...) memberi hasil yang sama? Lalu, kapan NOT IN menjadi berbahaya, dan apa penggantinya?

-- ya hasilnya tetap sama kalo user_id NOT NULL, NOT IN berbahaya saat kolom user_id bisa NULL, pengantinya adalah menggunakan NOT EXISTS
