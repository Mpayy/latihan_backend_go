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


-- Quiz: GROUP BY + agregasi

-- Skema: users(id, name), orders(id, user_id, total NUMERIC, status, created_at).

-- 1. Tampilkan status dan jumlah order per status, hanya untuk order tahun 2026, dan hanya status dengan lebih dari 100 order. Tulis query-nya, lalu sebutkan kondisi mana yang masuk WHERE dan mana yang masuk HAVING, beserta alasannya.

SELECT status, COUNT(id) AS jumlah_order
FROM orders
WHERE created_at >= '2026-01-01 00:00:00+07' AND created_at <  '2027-01-01 00:00:00+07'
GROUP BY status
HAVING COUNT(id) > 100;

-- 2. Tampilkan nama user beserta total belanja order berstatus 'paid', termasuk user yang tidak punya order paid (total = 0). Hati-hati dengan posisi filter status (ingat pelajaran JOIN) dan NULL dari SUM.

SELECT u.id, u.name, COALESCE(SUM(o.total), 0) AS total_belanja
FROM users AS u
LEFT JOIN orders AS o ON u.id = o.user_id AND o.status = 'paid'
GROUP BY u.id, u.name;


-- 3. Query ini error. Kenapa, dan bagaimana memperbaikinya?

-- sql
-- SELECT u.id, u.name, COUNT(o.id)
-- FROM users u LEFT JOIN orders o ON o.user_id = u.id
-- GROUP BY u.id;

-- Tebakan: apakah error juga di PostgreSQL? Jelaskan alasan jawabanmu.

-- karena kolom non agregrat di select ga masuk ke group by, dibeberapa database akan error. cara memperbaikinya dengan memasukan u.name ke group by
-- di postgres tidak akan error, karena postgres tau kalo nilai u.id adalah primary key dan pasti unik, yang membuat u.name gapapa kalo ga di masukin ke group by karena ada jaminan unik di u.id