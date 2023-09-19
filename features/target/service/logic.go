package service

import (
	"be_golang/klp3/features/target"
	usernodejs "be_golang/klp3/features/userNodejs"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
)

type targetService struct {
	targetRepo target.TargetDataInterface
	validate   *validator.Validate
}

func New(repo target.TargetDataInterface) target.TargetServiceInterface {
	return &targetService{
		targetRepo: repo,
		validate:   validator.New(),
	}
}

// Create implements target.TargetServiceInterface.
func (s *targetService) Create(input target.TargetEntity) (string, error) {
	userPembuat, err := usernodejs.GetByIdUser(input.UserIDPembuat)
	if err != nil {
		log.Printf("Error get detail user: %s", err.Error())
		return "", err
	}
	userPenerima, err := usernodejs.GetByIdUser(input.UserIDPenerima)
	if err != nil {
		log.Printf("Error get detail user: %s", err.Error())
		return "", err
	}
	err = s.validate.Struct(input)
	if err != nil {
		log.Printf("Error validate: %s", err.Error())
		return "", errors.New("error validate, konten target, due date required")
	}

	if userPembuat.Jabatan == "c-level" {
		if userPenerima.Jabatan != "karyawan" && userPenerima.Jabatan != "manager" {
			return "", errors.New("c-level hanya bisa membuat target untuk karyawan atau manager")
		}
	}

	if userPembuat.Jabatan == "manager" {
		if userPenerima.Jabatan != "karyawan" {
			return "", errors.New("manager hanya bisa membuat target untuk karyawan")
		}
	}

	if userPembuat.Jabatan != "c-level" && userPembuat.Jabatan != "manager" {
		return "", errors.New("jabatan Anda tidak memiliki izin untuk membuat target")
	}

	targetID, err := s.targetRepo.Insert(input)
	if err != nil {
		log.Printf("Error creating target: %s", err.Error())
		return "", err
	}
	log.Println("Target created successfully")
	return targetID, nil
}

// GetAll implements target.TargetServiceInterface.
func (s *targetService) GetAll(userID string, param target.QueryParam) (bool, []target.TargetEntity, error) {
	var total_page int64
	nextPage := true

	count, data, err := s.targetRepo.SelectAll(userID, param)
	if err != nil {
		return true, nil, err
	}
	if param.ExistOtherPage {
		total_page = count / int64(param.LimitPerPage)
		if count%int64(param.LimitPerPage) != 0 {
			total_page += 1
		}

		if param.Page == int(total_page) {
			nextPage = false
		}
	}
	return nextPage, data, nil
}

// GetById implements target.TargetServiceInterface.
func (s *targetService) GetById(targetID string, userID string) (target.TargetEntity, error) {
	result, err := s.targetRepo.Select(targetID, userID)
	if err != nil {
		return target.TargetEntity{}, err
	}
	return result, nil
}

// UpdateById implements target.TargetServiceInterface.
func (s *targetService) UpdateById(targetID string, userID string, targetData target.TargetEntity) error {
	// Dapatkan peran pengguna
	user, err := usernodejs.GetByIdUser(userID)
	if err != nil {
		return err
	}

	// Dapatkan target yang akan diperbarui
	existingTarget, err := s.targetRepo.Select(targetID, userID)
	if err != nil {
		return err
	}

	// Dapatkan pengguna dengan ID sesuai existingTarget.UserIDPenerima
	userTarget, err := usernodejs.GetByIdUser(existingTarget.UserIDPenerima)
	if err != nil {
		return err
	}

	// Inisialisasi variabel yang menunjukkan apakah pembaruan diizinkan
	allowedToUpdate := false

	// Pemeriksaan peran pengguna
	if user.Jabatan == "c-level" {
		allowedToUpdate = true
	}

	if user.Jabatan == "manager" {
		// Pemeriksaan apakah manajer dapat mengedit target karyawan atau target milik diri sendiri
		if userTarget.Jabatan == "karyawan" || existingTarget.UserIDPenerima == userID {
			allowedToUpdate = true
		}
	}

	if user.Jabatan == "karyawan" {
		// Pemeriksaan apakah karyawan dapat mengedit target milik diri sendiri
		if existingTarget.UserIDPenerima == userID {
			allowedToUpdate = true
		}
	}

	// Periksa izin pembaruan
	if !allowedToUpdate {
		return errors.New("anda tidak memiliki izin untuk mengedit target ini")
	}

	// Lakukan pembaruan hanya jika diizinkan
	err = s.targetRepo.Update(targetID, userID, targetData)
	if err != nil {
		return err
	}
	return nil
}

// DeleteById implements target.TargetServiceInterface.
func (s *targetService) DeleteById(targetID string, userID string) error {
	err := s.targetRepo.Delete(targetID, userID)
	if err != nil {
		return err
	}
	return nil
}
