package entities

type User struct {
	Id          int64
	NamaLengkap string `validate:"required" label:"Nama Lengkap"`
	Email       string `validate:"required,email,isunique=login-email"`
	Username    string `validate:"required,gte=3"`
	Password    string `validate:"required,gte=8"`
	Cpassword   string `validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}