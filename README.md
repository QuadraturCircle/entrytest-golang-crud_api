# Entry Test: Simple CRUD REST API in Go/PostgreSQL

Contoh penggunaan Go dalam pembuatan CRUD REST API, dengan sampel data hewan kurban. Database yang digunakan adalah PostgreSQL.

## Dependensi

- [pq](https://github.com/lib/pq)
- [mux](https://github.com/gorilla/mux)

## Instalasi dan Penggunaan

- Unduh dari github (*Download as ZIP*).
- Ekstrak di tempat folder yang diinginkan.
- Ubah `const` untuk menuju ke instalasi PostgreSQL yang dimiliki (*host*, *port*, *username*, *password*, *database name*).
- Jalankan `go run main.go`.
- Gunakan [Postman](https://www.postman.com/) untuk melakukan pengujian.

## Daftar *Request*

- `getKurbanAll`: `GET`, `/kurban`, mengambil semua data yang ada di dalam tabel
- `getKurban`: `GET`, `/kurban{id}`, mengambil data spesifik sesuai parameter ID yang diberikan
- `createKurban`: `POST`, `/kurban`, menambahkan entri data ke dalam tabel sesuai dengan parameter data yang diberikan
- `updateKurban`: `PUT`, `/kurban{id}`, memperbarui entri spesifik di dalam tabel sesuai parameter ID dan data yang diberikan
- `deleteKurban`: `DELETE`, `/kurban{id}`, menghapus entri spesifik di dalam tabel sesuai parameter ID yang diberikan

## Tabel

```sql
CREATE TABLE IF NOT EXISTS kurban (
    id SERIAL PRIMARY KEY, 
    kurbanName TEXT NOT NULL, 
    kurbanType TEXT NOT NULL, 
    kurbanWeight INT NOT NULL, 
    kurbanPrice INT NOT NULL
)
```
