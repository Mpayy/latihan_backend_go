-- Topik 1: JOIN

-- Konsep: JOIN menggabungkan baris dua tabel berdasarkan kondisi kolom, biasanya FK = PK. INNER JOIN hanya mengembalikan baris yang cocok di kedua sisi. LEFT JOIN mempertahankan semua baris tabel kiri, dan kolom kanan berisi NULL jika tidak ada pasangan. Pilih LEFT JOIN kalau baris tanpa relasi tetap harus muncul.

-- Contoh:

-- go
-- const q = `
-- SELECT o.id, o.total, u.name
-- FROM orders o
-- JOIN users u ON u.id = o.user_id
-- WHERE o.created_at >= $1`
-- rows, err := db.QueryContext(ctx, q, since)

-- Latihan (skema: users(id, name), orders(id, user_id, total NUMERIC)):

-- Tulis query untuk menampilkan semua user beserta id dan total order-nya, termasuk user yang belum punya order.

SELECT u.id, u.name, ,o.id, o.total
FROM users AS u
LEFT JOIN orders AS o
ON u.id = o.user_id;

-- Tulis query untuk menampilkan hanya user yang belum pernah order. Petunjuk: manfaatkan NULL dari LEFT JOIN.

SELECT u.name
FROM users AS u
LEFT JOIN orders AS o
ON u.id = o.user_id AND o.id IS NULL;

-- Di Go, kalau kolom o.total bisa NULL, tipe apa yang kamu pakai untuk Scan, dan kenapa float64 biasa bermasalah?

-- tipe data yang saya gunakan adalah *float64 atau menggunakan sql.NullFloat64 dan float64 biasa bermasalah karena tidak bisa membedakan NULL dengan 0, dan saat Scan NULL ke float64 akan error


-- Quiz: JOIN

-- 1. Skema: users(id, name), orders(id, user_id, total). Tulis query untuk menampilkan user yang punya order dengan total > 100, tanpa duplikat nama meskipun user punya banyak order seperti itu. Jangan pakai DISTINCT di semua kolom secara asal; jelaskan pilihanmu.

SELECT DISTINCT u.id, u.name
FROM users AS u
JOIN orders AS o ON u.id = o.user_id
WHERE o.total > 100;

-- with subquery
SELECT u.id, u.name
FROM users AS u
WHERE EXISTS (
    SELECT 1 
    FROM orders AS o 
    WHERE u.id = o.user_id AND o.total > 100
);

-- 2. Apa hasil SELECT COUNT(*) dan SELECT COUNT(o.id) pada query ini, kalau users berisi 3 baris dan hanya 1 user yang punya 2 order?

-- sql
-- SELECT ... FROM users u LEFT JOIN orders o ON o.user_id = u.id

-- Sebutkan angka keduanya dan kenapa berbeda.

-- hasil dari SELECT COUNT(*) adalah 4 dan hasil dari SELECT COUNT(o.id) adalah 2. hasilnya berbeda karena SELECT COUNT(*) menghitung semua baris, dan SELECT COUNT(o.id) hanya mmenghitung baris yang o.id nya tidak NULL.

-- 3. Kamu menulis FROM orders o JOIN users u ON u.id = o.user_id WHERE u.active = true AND o.total > 100. Kalau diganti LEFT JOIN, apakah hasilnya berubah? Jelaskan.

-- tidak berubah karena filter di table kanan, maka akan menetralkan left join. karena user yang ga punya order akan di buang