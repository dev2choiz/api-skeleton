package fixtures

import (
	"time"

	"github.com/dev2choiz/api-skeleton/entity"
)

var Users = []entity.User{
	{
		ID:        "38dfac55-8ff6-4c9d-8916-9016052f50cd",
		Username:  "geralt",
		Firstname: "Geralt",
		Lastname:  "de Riv",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Username:  "ciri",
		Firstname: "Cirilla",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
	{
		Username:  "yen",
		Firstname: "Yennefer",
		Lastname:  "Vengerberg",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
	{
		Username:  "fitz",
		Firstname: "Fitz",
		Lastname:  "Chevalerie",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
	{
		Username:  "abeille",
		Firstname: "Abeille",
		Lastname:  "Chevalerie",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
	{
		Username:  "richard",
		Firstname: "Richard",
		Lastname:  "Rahl",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
	{
		Username:  "belgarion",
		Firstname: "Garion",
		Lastname:  "Bel",
		Password:  "$2a$12$0tM0U3Yt7.ySsMKqdajmH.K2yQeSNlkI55rk9BDraEHp7HjGTFnG.",
	},
}
