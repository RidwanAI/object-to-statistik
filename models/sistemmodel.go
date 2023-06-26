package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RidwanAI/object-to-statistik/config"
	"github.com/RidwanAI/object-to-statistik/entities"
)

type GudangModel struct {
	conn *sql.DB
}

func NewGudangModel() *GudangModel {
	conn, err := config.DBConnect()
	if err != nil {
		panic(err)
	}
	
	return &GudangModel{
		conn: conn,
	}
}

func (g *GudangModel) GetData() ([]entities.Gudang, error ) {

	rows, err := g.conn.Query("select id, nama_produk, jumlah_stok from sistem")
    if err != nil {
        return []entities.Gudang{}, err
    }
    defer rows.Close()

	var visual []entities.Gudang
    for rows.Next() {
        var data entities.Gudang 
        rows.Scan(&data.Id, 
			&data.NamaProduk, 
			&data.JumlahStok)
        visual = append(visual, data)
    }

    return visual, nil

}

func (g *GudangModel) BacaData() ([]entities.Gudang, error) {
	rows, err := g.conn.Query("select * from sistem")
	if err != nil {
		return []entities.Gudang{}, err
	}
	defer rows.Close()

	var dataGudang []entities.Gudang
	for rows.Next() {
		var dGudang entities.Gudang
		rows.Scan(&dGudang.Id,
			&dGudang.NamaProduk,
			&dGudang.KodeProduk,
			&dGudang.JenisProduk,
			&dGudang.TanggalUpdate,
			&dGudang.JumlahStok,
			&dGudang.NamaPemasok,
			&dGudang.AlamatPemasok)

		tanggal_lahir, _ := time.Parse("2006-01-02", dGudang.TanggalUpdate)
		dGudang.TanggalUpdate = tanggal_lahir.Format("02-01-2006")
	
		dataGudang = append(dataGudang, dGudang)
	}

	return dataGudang, nil

}

func (g *GudangModel) UpdateStok(gudang entities.Gudang) bool {
	
	result, err := g.conn.Exec("insert into sistem (nama_produk, kode_produk, jenis_produk, tanggal_update, jumlah_stok, nama_pemasok, alamat_pemasok) values(?,?,?,?,?,?,?)",
		gudang.NamaProduk, gudang.KodeProduk, gudang.JenisProduk, gudang.TanggalUpdate, gudang.JumlahStok, gudang.NamaPemasok, gudang.AlamatPemasok)

	if err != nil {
		fmt.Println(err)
		return false
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId > 0

}

func (g *GudangModel) CariQueri(id int64, gudang *entities.Gudang) error {

	return g.conn.QueryRow("select * from sistem where id = ?", id).Scan(
		&gudang.Id,
		&gudang.NamaProduk,
		&gudang.KodeProduk,
		&gudang.JenisProduk,
		&gudang.TanggalUpdate,
		&gudang.JumlahStok,
		&gudang.NamaPemasok,
		&gudang.AlamatPemasok,
	)
}

func (g *GudangModel) EditStok(gudang entities.Gudang) error {
	
	_, err := g.conn.Exec(
		"update sistem set nama_produk = ?, kode_produk = ?, jenis_produk = ?, tanggal_update = ?, jumlah_stok = ?, nama_pemasok = ?, alamat_pemasok = ? where id = ?",
		gudang.NamaProduk, gudang.KodeProduk, gudang.JenisProduk, gudang.TanggalUpdate, gudang.JumlahStok, gudang.NamaPemasok, gudang.AlamatPemasok, gudang.Id)

	if err != nil {
		return err
	}

	return nil

}

func (g *GudangModel) HapusStok(id int64) {
	
	g.conn.Exec("delete from sistem where id = ?", id)

}
