package entities

type Gudang struct {
	Id            int64
	NamaProduk    string `validate:"required" label:"Nama Produk"`
	KodeProduk    string `validate:"required" label:"Kode Produk"`
	JenisProduk   string `validate:"required" label:"Jenis Produk"`
	TanggalUpdate string `validate:"required" label:"Tanggal Update Produk"`
	JumlahStok    string `validate:"required" label:"Jumlah Stok Produk"`
	NamaPemasok   string `validate:"required" label:"Nama Pemasok Produk"`
	AlamatPemasok string `validate:"required" label:"Alamat Pemasok Produk"`
}