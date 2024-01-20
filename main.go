package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// konstant yang digunakan untuk penyambungan database
const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

// struktur data yang digunakan
type Kurban struct {
	KurbanID     int    `json:"id"`
	KurbanName   string `json:"kurbanName"`
	KurbanType   string `json:"kurbanType"`
	KurbanWeight int    `json:"kurbanWeight"`
	KurbanPrice  int    `json:"kurbanPrice"`
}

// struktur data untuk pesan pasca-operasi
type JsonResponse struct {
	Type    string   `json:"type"`
	Data    []Kurban `json:"data"`
	Message string   `json:"message"`
}

func main() {
	//buat link router
	router := mux.NewRouter()
	router.HandleFunc("/kurban", getKurbanAll).Methods("GET")
	router.HandleFunc("/kurban/{id}", getKurban).Methods("GET")
	router.HandleFunc("/kurban", createKurban).Methods("POST")
	router.HandleFunc("/kurban/{id}", updateKurban).Methods("PUT")
	router.HandleFunc("/kurban/{id}", deleteKurban).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// koneksi database
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	//buat tabel jika belum ada
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS kurban (id SERIAL PRIMARY KEY, kurbanName TEXT NOT NULL, kurbanType TEXT NOT NULL, kurbanWeight INT NOT NULL, kurbanPrice INT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// message handler
func printMessage(message string) {
	fmt.Println(message)
	fmt.Println("")
}

// error checker
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// # getKurbanAll: GET
//
// Menampilkan semua data yang ada di tabel.
//
// Expected responses:
//   - jika tabel kosong: null data
//   - jika ada tabel: sukses, respons semua data dalam bentuk JSON
func getKurbanAll(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Mengambil daftar kurban...")

	var response = JsonResponse{}

	rows, err := db.Query("SELECT * FROM kurban")
	if err != nil {
		response = JsonResponse{Type: "error", Message: "!!! Tabel kurban kosong!"}
	} else {
		var kurbanResult []Kurban
		for rows.Next() {
			var k Kurban
			err = rows.Scan(&k.KurbanID, &k.KurbanName, &k.KurbanType, &k.KurbanWeight, &k.KurbanPrice)
			checkErr(err)

			kurbanResult = append(kurbanResult, k)
		}
		response = JsonResponse{Type: "success", Data: kurbanResult}
	}
	json.NewEncoder(w).Encode(response)
}

// # getKurban: GET
//
// Menampilkan data yang sesuai dengan ID yang diberikan sebagai parameter.
//
// Expected responses:
//   - jika param ID kosong: error, "!!! ID kurban belum dimasukkan!"
//   - jika ada tabel: sukses, respons semua data valid dalam bentuk JSON
func getKurban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var response = JsonResponse{}

	if id == "" {
		response = JsonResponse{Type: "error", Message: "!!! ID kurban belum dimasukkan!"}
	} else {
		db := setupDB()
		printMessage("Mengambil data entri kurban...")

		var k Kurban

		err := db.QueryRow("SELECT * FROM kurban WHERE id = $1", id).Scan(&k.KurbanID, &k.KurbanName, &k.KurbanType, &k.KurbanWeight, &k.KurbanPrice)
		if err != nil {
			response = JsonResponse{Type: "error", Message: "!! ID kurban salah atau tidak ditemukan!"}
		} else {
			var kurbanResult []Kurban
			kurbanResult = append(kurbanResult, k)

			response = JsonResponse{Type: "success", Data: kurbanResult}
		}
	}

	json.NewEncoder(w).Encode(response)
}

// # createKurban: POST
//
// Menambahkan entri baru ke dalam tabel, sesuai dengan parameter yang telah diberikan.
//
// Expected responses:
//   - jika param ID/data kosong: error, "!!! Data kurban belum dimasukkan!"
//   - jika param ID/data ada: sukses, "Data telah dimasukkan."
func createKurban(w http.ResponseWriter, r *http.Request) {
	kurbanName := r.FormValue("kurbanName")
	kurbanType := r.FormValue("kurbanType")
	kurbanWeight := r.FormValue("kurbanWeight")
	kurbanPrice := r.FormValue("kurbanPrice")

	var response = JsonResponse{}

	if kurbanName == "" || kurbanType == "" || kurbanWeight == "" || kurbanPrice == "" {
		response = JsonResponse{Type: "error", Message: "!!! Data kurban belum dimasukkan!"}
	} else {
		db := setupDB()
		printMessage("Memasukkan entri baru kurban...")
		fmt.Println("Memasukkan entri baru dengan nama " + kurbanName + ", tipe " + kurbanType + ", berat " + kurbanWeight + " dan harga " + kurbanPrice)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO kurban(kurbanName, kurbanType, kurbanWeight, kurbanPrice) VALUES($1, $2, $3, $4) RETURNING id;", kurbanName, kurbanType, kurbanWeight, kurbanPrice).Scan(&lastInsertID)
		if err != nil {
			response = JsonResponse{Type: "error", Message: "!! Data gagal dimasukkan!"}
		} else {
			response = JsonResponse{Type: "success", Message: "Data telah dimasukkan."}
		}
	}
	json.NewEncoder(w).Encode(response)
}

// # updateKurban: PUT
//
// Memperbarui entri spesifik dalam tabel, sesuai dengan parameter ID dan data baru yang telah diberikan.
//
// Expected responses:
//   - jika param ID/data kosong: error, "!!! Data kurban belum dimasukkan!"
//   - jika param ID/data ada: sukses, "Data telah diperbarui."
func updateKurban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	kurbanName := r.FormValue("kurbanName")
	kurbanType := r.FormValue("kurbanType")
	kurbanWeight := r.FormValue("kurbanWeight")
	kurbanPrice := r.FormValue("kurbanPrice")

	var response = JsonResponse{}

	if id == "" || kurbanName == "" || kurbanType == "" || kurbanWeight == "" || kurbanPrice == "" {
		response = JsonResponse{Type: "error", Message: "!!! Data kurban belum dimasukkan!"}
	} else {
		db := setupDB()
		printMessage("Memperbarui entri kurban...")
		fmt.Println("Memperbarui entri kurban dengan data baru: nama " + kurbanName + ", tipe " + kurbanType + ", berat " + kurbanWeight + " dan harga " + kurbanPrice)

		_, err := db.Exec("UPDATE kurban SET kurbanName = $1, kurbanType = $2, kurbanWeight = $3, kurbanPrice = $4 WHERE id = $5", kurbanName, kurbanType, kurbanWeight, kurbanPrice, id)
		if err != nil {
			response = JsonResponse{Type: "error", Message: "!! Data gagal diperbarui!"}
		} else {
			response = JsonResponse{Type: "success", Message: "Data telah diperbarui."}
		}
	}
	json.NewEncoder(w).Encode(response)
}

// # deleteKurban: DELETE
//
// Menghapus entri spesifik dalam tabel, sesuai dengan parameter ID yang telah diberikan.
//
// Expected responses:
//   - jika param ID kosong: error, "!!! ID kurban belum dimasukkan!"
//   - jika param ID ada: sukses, "Entri telah dihapus."
func deleteKurban(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var response = JsonResponse{}

	if id == "" {
		response = JsonResponse{Type: "error", Message: "!!! ID kurban belum dimasukkan!"}
	} else {
		db := setupDB()
		printMessage("Menghapus entri kurban...")

		err := db.QueryRow("SELECT id FROM kurban WHERE id = $1", id)
		if err != nil {
			response = JsonResponse{Type: "error", Message: "ID kurban salah atau tidak ditemukan!"}
		} else {
			_, err := db.Exec("DELETE FROM kurban WHERE id = $1", id)
			if err != nil {
				response = JsonResponse{Type: "error", Message: "!! Entri gagal dihapus!"}
			} else {
				response = JsonResponse{Type: "success", Message: "Entri telah dihapus."}
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}
