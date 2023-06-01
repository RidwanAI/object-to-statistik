package main

import (
	"net/http"

	gudangcontroller "github.com/RidwanAI/object-to-statistik/controllers"
)

func main() {
	http.HandleFunc("/", gudangcontroller.Index)
	http.HandleFunc("/login", gudangcontroller.Login)
	http.HandleFunc("/logout", gudangcontroller.Logout)
	http.HandleFunc("/register", gudangcontroller.Register)
	http.HandleFunc("/gudang", gudangcontroller.Index)
	http.HandleFunc("/gudang/index", gudangcontroller.Index)
	http.HandleFunc("/gudang/update", gudangcontroller.Add)
	http.HandleFunc("/gudang/ubah", gudangcontroller.Edit)
	http.HandleFunc("/gudang/hapus", gudangcontroller.Delete)

	http.ListenAndServe(":9000", nil)
}