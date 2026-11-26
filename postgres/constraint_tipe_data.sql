-- Topik 5: Constraint & tipe data

-- Konsep: Constraint memindahkan aturan data ke database, jadi data tetap valid meski ada bug di aplikasi atau ada proses lain yang menulis: PRIMARY KEY (unik + NOT NULL), FOREIGN KEY (relasi harus ada induknya), UNIQUE, CHECK (kondisi bebas), NOT NULL. Tipe data menentukan ketepatan: BIGINT untuk ID/hitungan besar, NUMERIC(p,s) untuk uang (eksak), TIMESTAMPTZ untuk waktu (disimpan sebagai UTC).

-- Contoh:

-- go
-- const ddl = `
-- CREATE TABLE wallets (
--   id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   user_id    BIGINT NOT NULL REFERENCES users(id),
--   balance    NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
--   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--   UNIQUE (user_id)
-- )`

-- Latihan:

-- Rancang tabel orders(id, user_id, total, status, created_at). Tentukan tipe, NOT NULL, FK, dan satu CHECK untuk total. Tambahkan CHECK agar status hanya boleh 'pending', 'paid', atau 'cancelled'.

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    total NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled'))
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

CREATE INDEX idx_orders_user ON orders(user_id);

-- Kenapa TIMESTAMP (tanpa TZ) berisiko dibanding TIMESTAMPTZ?
-- Karena TIMESTAMP tidak menyimpan zona waktu, jadi waktu di tulis apa adanya di db. jadi bisa terjadi abigu saat kita menarik datanya. TIMESTAMPTZ menyimpan zona waktu, jadi saat menyimpan akan di convert ke UTC 

-- Di wallet-mu, CHECK (balance >= 0) ada. Kenapa itu tetap berguna meski aplikasi sudah mengecek saldo sebelum debit?
-- untuk menjaga konsisten data level database, jika terdapat bug atau race condition di aplikasi
-- jika dua request debit bersamaan sama-sama lolos cek saldo di aplikasi (race), lalu keduanya menulis. CHECK menjadi jaring pengaman terakhir: salah satu UPDATE gagal, bukan saldo jadi negatif. Tapi CHECK bukan pengganti locking. Ia hanya mencegah data rusak, tidak mengatur urutan.


-- Quiz: Constraint & tipe data

-- 1. Rancang tabel wallets(id, user_id, balance, currency, created_at) dan wallet_transactions(id, wallet_id, amount, type, idempotency_key, created_at). Tentukan tipe, NOT NULL, FK, dan constraint yang menurutmu perlu. Syarat: satu user hanya punya satu wallet per currency, type hanya 'credit' atau 'debit', amount harus positif, dan idempotency_key tidak boleh dipakai dua kali dalam wallet yang sama.

CREATE TABLE IF NOT EXISTS wallets (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    balance NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(balance >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR' CHECK (currency IN ('IDR', 'USD', 'SGD')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, currency)
);

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    amount NUMERIC(18,2) NOT NULL CHECK (amount > 0),
    type TEXT NOT NULL CHECK (type IN ('credit', 'debit')),
    idempotency_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(wallet_id, idempotency_key)
);

-- 2. Kolom balance kamu pakai NUMERIC(18,2). Jelaskan kenapa float64 bermasalah untuk uang, dan tipe apa di Go yang kamu pakai untuk Scan kolom ini.

-- karena float64 memakai representasi biner (IEEE 754) yang tidak dapat menyimpan banyak pecahan desimal basis 10 secara tepat, sehingga memicu pembulatan saat perhitungan dan perbandingan uang
-- untuk scan tipe data NUMERIC saya akan pakai package shopspring/decimal 

-- 3. Kamu menambah kolom deleted_at TIMESTAMPTZ untuk soft delete, lalu UNIQUE (user_id, currency) mulai bermasalah: wallet yang sudah di-soft-delete menghalangi pembuatan wallet baru. Apa solusinya di level index? (Petunjuk: ini menyentuh partial index, cukup jelaskan idenya.)

-- Solusinya memakai partial unique index pada (user_id, currency) dengan syarat WHERE deleted_at IS NULL, sehingga aturan unik hanya berlaku pada wallet yang masih aktif (deleted_at IS NULL)