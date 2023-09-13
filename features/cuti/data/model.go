package data

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cuti struct {
	ID        		uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
	DeletedAt 		gorm.DeletedAt `gorm:"index"`
	TipeCuti  		string `gorm:"type:enum('melahirkan','hari raya','tahunan');default:'tahunan';column:tipe_cuti"`
	Status    		string
	JumlahCuti 		int
	BatasanCuti 	int  `gorm:"column:batasan_cuti;default:90"`
	Description 	string
	Persetujuan 	time.Time
	StartCuti 		time.Time
	EndCuti 		time.Time
	UrlPendukung 	string
	UserID 			string
}