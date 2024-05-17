package contact

import (
	"database/sql"
	"log"
	"log/slog"
	"time"
)

type Contact struct {
	ContactID   int       `json:"contact_id"`
	Name        string    `json:"name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	GenderID    int       `json:"gender_id"`
	DOB         string    `json:"dob"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	PhotoPath   string    `json:"photo_path"`
	CreatedDate time.Time `json:"created_date"`
	CreatedBy   string    `json:"created_by"`
}

func GetContactData() ([]Contact, error) {
	db, err := sql.Open("sqlite3", "./graphqldb.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("Select * FROM contact LIMIT 10")
	if err != nil {
		return nil, err
	}

	var result []Contact

	defer rows.Close()

	for rows.Next() {
		var tempContact Contact
		if err := rows.Scan(&tempContact.ContactID, &tempContact.Name, &tempContact.FirstName, &tempContact.LastName, &tempContact.GenderID, &tempContact.DOB, &tempContact.Email, &tempContact.Phone, &tempContact.Address, &tempContact.PhotoPath, &tempContact.CreatedDate, &tempContact.CreatedBy); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		result = append(result, tempContact)
	}

	// jsonString, err := json.Marshal(result)
	// if err != nil {
	// 	return nil, err
	// }

	return result, nil
}
