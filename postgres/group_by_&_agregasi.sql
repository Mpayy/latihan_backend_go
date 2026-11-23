-- Topik 2: GROUP BY + agregasi

-- Konsep: GROUP BY mengelompokkan baris berdasarkan nilai kolom, lalu fungsi agregat (COUNT, SUM, AVG, MIN, MAX) dihitung per kelompok. Kolom di SELECT harus ada di GROUP BY atau dibungkus agregat. WHERE memfilter baris sebelum dikelompokkan, HAVING memfilter kelompok setelah agregasi.

-- Contoh:

-- go
-- const q = `
-- SELECT user_id, COUNT(*) AS n, SUM(total) AS spent
-- FROM orders
-- WHERE created_at >= $1
-- GROUP BY user_id
-- HAVING SUM(total) > 1000`

-- Latihan (skema: users(id, name), orders(id, user_id, total NUMERIC)):

-- Tampilkan u.name dan jumlah order per user, termasuk user dengan 0 order (ingat pelajaran COUNT).

SELECT u.name, COUNT(o.id) AS jumlah_order
FROM users AS u
LEFT JOIN orders AS o ON u.id = u.user_id
GROUP BY u.id, u.name;

-- Dari query itu, tampilkan hanya user dengan total belanja > 500. Kenapa kondisi ini tidak bisa ditaruh di WHERE?

SELECT u.name, COUNT(o.id)
FROM users AS u
JOIN orders AS o u.id = o.user_id
GROUP BY u.id, u.name
HAVING SUM(o.total) > 500;

-- karena WHERE jalan sebelum GROUP BY, karena belum di kelompokan jadi gabisa SUM total per usernya.

-- SUM(o.total) untuk user tanpa order menghasilkan NULL, bukan 0. Bagaimana membuatnya jadi 0 di SQL? Dan kalau tanpa itu, tipe Go apa yang aman untuk Scan?

-- menggunakan COALESCE pada SUM(o.total) seperti ini COALESCE(SUM(o.total), 0) untuk membuat NULL jadi 0 di SQL.
-- tipe data yang aman di Go untuk Scan adalah pointer tipe data, untuk membedakan NULL atau 0
