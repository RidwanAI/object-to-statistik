package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/RidwanAI/object-to-statistik/config"
	"github.com/RidwanAI/object-to-statistik/entities"
	"github.com/RidwanAI/object-to-statistik/libraries"
	"github.com/RidwanAI/object-to-statistik/models"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Username string
	Password string
}

var userModel = models.NewUserModel()
var gudangModel = models.NewGudangModel()
var validation = libraries.NewValidation()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSIONS_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			gudang, _ := gudangModel.BacaData()

			data := map[string]interface{}{
				"nama_lengkap": session.Values["nama_lengkap"],
				"dGudang": gudang,
			}

			temp, err := template.ParseFiles("views/index.html")

			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)
		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		r.ParseForm()
		UserInput := &Input{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		errorMessages := validation.Struct(UserInput)

		if errorMessages != nil {

			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/login.html")
			temp.Execute(w, data)

		} else {

			var user entities.User
			userModel.Where(&user, "username", UserInput.Username)

			var message error
			if user.Username == "" {
				message = errors.New("Username atau Password salah!")
			} else {
				// pengecekan password
				errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(UserInput.Password))
				if errPassword != nil {
					message = errors.New("Username atau Password salah!")
				}
			}

			if message != nil {

				data := map[string]interface{}{
					"error": message,
				}

				temp, _ := template.ParseFiles("views/login.html")
				temp.Execute(w, data)
			} else {
				// set session
				session, _ := config.Store.Get(r, config.SESSIONS_ID)

				session.Values["loggedIn"] = true
				session.Values["email"] = user.Email
				session.Values["username"] = user.Username
				session.Values["nama_lengkap"] = user.NamaLengkap

				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}

	}

}

func Logout(w http.ResponseWriter, r *http.Request)  {

	session, _ := config.Store.Get(r, config.SESSIONS_ID)

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func Register(w http.ResponseWriter, r *http.Request)  {
	
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		user := entities.User {
			NamaLengkap: r.Form.Get("nama_lengkap"),
			Email: r.Form.Get("email"),
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
			Cpassword: r.Form.Get("cpassword"),
		}

		errorMessages := validation.Struct(user)

		if errorMessages != nil {

			data := map[string]interface{}{
				"validation": errorMessages,
				"user": 	  user,
			}

			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPassword)

			userModel.Create(user)

			data := map[string]interface{} {
				"pesan": "Register berhasil",
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		}
	}

}

func Add(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSIONS_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			
			if r.Method == http.MethodGet {
				temp, err := template.ParseFiles("views/update.html")
				if err != nil {
					panic(err)
				}
				temp.Execute(w, nil)
			} else if r.Method == http.MethodPost {
		
				r.ParseForm()
		
				var gudang entities.Gudang
				gudang.NamaProduk = r.Form.Get("nama_produk")
				gudang.KodeProduk = r.Form.Get("kode_produk")
				gudang.JenisProduk = r.Form.Get("jenis_produk")
				gudang.TanggalUpdate = r.Form.Get("tanggal_update")
				gudang.JumlahStok = r.Form.Get("jumlah_stok")
				gudang.NamaPemasok = r.Form.Get("nama_pemasok")
				gudang.AlamatPemasok = r.Form.Get("alamat_pemasok")
		
				var data = make(map[string]interface{})
		
				vErrors := validation.Struct(gudang)
		
				if vErrors != nil {
					data["gudang"] = gudang
					data["validation"] = vErrors
				} else {
					data["pesan"] = "Data Produk Berhasil Disimpan"
					gudangModel.UpdateStok(gudang)
				}
		
				temp, _ := template.ParseFiles("views/update.html")
				temp.Execute(w, data)
			}

		}
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSIONS_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			
			if r.Method == http.MethodGet {

				queryString := r.URL.Query()
				id, _ := strconv.ParseInt(queryString.Get("id"), 10, 64)
		
				var gudang entities.Gudang
				gudangModel.CariQueri(id, &gudang)
		
				data := map[string]interface{}{
					"gudang": gudang,
				}
		
				temp, err := template.ParseFiles("views/ubah.html")
				if err != nil {
					panic(err)
				}
				temp.Execute(w, data)
			} else if r.Method == http.MethodPost {
		
				r.ParseForm()
		
				var gudang entities.Gudang
				gudang.Id, _ = strconv.ParseInt(r.Form.Get("id"), 10, 64)
				gudang.NamaProduk = r.Form.Get("nama_produk")
				gudang.KodeProduk = r.Form.Get("kode_produk")
				gudang.JenisProduk = r.Form.Get("jenis_produk")
				gudang.TanggalUpdate = r.Form.Get("tanggal_update")
				gudang.JumlahStok = r.Form.Get("jumlah_stok")
				gudang.NamaPemasok = r.Form.Get("nama_pemasok")
				gudang.AlamatPemasok = r.Form.Get("alamat_pemasok")
		
				var data = make(map[string]interface{})
		
				vErrors := validation.Struct(gudang)
		
				if vErrors != nil {
					data["gudang"] = gudang
					data["validation"] = vErrors
				} else {
					data["pesan"] = "Data Produk Berhasil Diperbarui"
					gudangModel.EditStok(gudang)
				}
		
				temp, _ := template.ParseFiles("views/ubah.html")
				temp.Execute(w, data)
			}

		}
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, _ := strconv.ParseInt(queryString.Get("id"), 10, 64)

	gudangModel.HapusStok(id)

	http.Redirect(w, r, "/gudang", http.StatusSeeOther)

}