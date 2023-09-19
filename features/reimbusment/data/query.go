package data

import (
	"be_golang/klp3/features/reimbusment"
	usernodejs "be_golang/klp3/features/userNodejs"
	"be_golang/klp3/helper"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ReimbusmentData struct {
	db *gorm.DB
}

// Delete implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) Delete(id string) error {
	var inputModel Reimbursement
	tx := repo.db.Where("id=?", id).Delete(&inputModel)
	if tx.Error != nil {
		return errors.New("delete error reimbursement")
	}
	if tx.RowsAffected == 0 {
		return errors.New("row not affected")
	}
	return nil
}

// SelectAll implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) SelectAll(param reimbusment.QueryParams) (int64, []reimbusment.ReimbursementEntity, error) {
	var inputModel []Reimbursement
	var total_reimbursement int64

	query := repo.db

	if param.IsClassDashboard {
		offset := (param.Page - 1) * param.ItemsPerPage
		fmt.Println("offset",offset)
		if param.SearchName != "" {
			query = query.Where("description like ?", "%"+param.SearchName+"%")
		}
		tx := query.Find(&inputModel)
		if tx.Error != nil {
			return 0, nil, errors.New("failed get all reimbursement")
		}
		total_reimbursement = tx.RowsAffected
		query = query.Offset(offset).Limit(param.ItemsPerPage)
	}
	if param.SearchName != "" {
		query = query.Where("description like ?", "%"+param.SearchName+"%")
	}
	tx := query.Find(&inputModel)
	if tx.Error != nil {
		return 0, nil, errors.New("error get all reimbursement")
	}

	dataPengguna, errUser := usernodejs.GetAllUser()
	if errUser != nil {
		return 0, nil, errUser
	}
	var dataUser []User
	for _, value := range dataPengguna {
		dataUser = append(dataUser, PenggunaToUser(value))
	}
	var userEntity []reimbusment.UserEntity
	for _, value := range dataUser {
		userEntity = append(userEntity, UserToEntity(value))
	}
	fmt.Println("user entity",userEntity)
	var reimbushPengguna []ReimbursementPengguna
	for _, value := range inputModel {
		reimbushPengguna = append(reimbushPengguna, ModelToPengguna(value))
	}
	fmt.Println("reimb",reimbushPengguna)
	var reimbushEntity []reimbusment.ReimbursementEntity
	for i := 0; i < len(userEntity); i++ {
		for j := 0; j < len(reimbushPengguna); j++ {
			if userEntity[i].ID == reimbushPengguna[j].UserID {
				reimbushPengguna[j].User = User(userEntity[i])
				reimbushEntity = append(reimbushEntity, PenggunaToEntity(reimbushPengguna[j]))
			}
		}
	}
	fmt.Println("reimbursement user",reimbushEntity)
	return total_reimbursement, reimbushEntity, nil
}

// SelectAllKaryawan implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) SelectAllKaryawan(idUser string, param reimbusment.QueryParams) (int64, []reimbusment.ReimbursementEntity, error) {

	var inputModel []Reimbursement
	var total_reimbursement int64

	query := repo.db

	if param.IsClassDashboard {
		offset := (param.Page - 1) * param.ItemsPerPage
		if param.SearchName != "" {
			query = query.Where("user_id=? and description like ?", idUser, "%"+param.SearchName+"%")
		}
		tx := query.Find(&inputModel)
		if tx.Error != nil {
			return 0, nil, errors.New("failed get all reimbursement")
		}
		total_reimbursement = tx.RowsAffected
		query = query.Offset(offset).Limit(param.ItemsPerPage)
	}
	if param.SearchName != "" {
		query = query.Where("user_id=? and description like ?", idUser, "%"+param.SearchName+"%")
	}
	tx := query.Find(&inputModel)
	if tx.Error != nil {
		return 0, nil, errors.New("error get all reimbursement karyawan")
	}

	dataUser, errUser := usernodejs.GetByIdUser(idUser)
	if errUser != nil {
		return 0, nil, errUser
	}
	pengguna := PenggunaToUser(dataUser)
	userEntity := UserToEntity(pengguna)

	var reimbushPengguna []ReimbursementPengguna
	for _, value := range inputModel {
		reimbushPengguna = append(reimbushPengguna, ModelToPengguna(value))
	}
	var reimbushEntity []reimbusment.ReimbursementEntity
	for _, value := range reimbushPengguna {
		if value.UserID == userEntity.ID {
			value.User = User(userEntity)
			reimbushEntity = append(reimbushEntity, PenggunaToEntity(value))
		}
	}
	return total_reimbursement, reimbushEntity, nil
}

// UpdateKaryawan implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) UpdateKaryawan(input reimbusment.ReimbursementEntity, id string) error {
	inputModel := EntityToModel(input)
	tx := repo.db.Model(&Reimbursement{}).Where("id=? and user_id=?", id, input.UserID).Updates(inputModel)
	if tx.Error != nil {
		return errors.New("update data reimbursment error, hanya boleh mengedit reimbursment sendiri")
	}
	if tx.RowsAffected == 0 {
		return errors.New("row not affected, hanya dapat mengedit reimbursement sendiri")
	}
	return nil
}

// SelectById implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) SelectById(id string) (reimbusment.ReimbursementEntity, error) {
	var inputModel Reimbursement
	tx := repo.db.Where("id=?", id).First(&inputModel)
	if tx.Error != nil {
		return reimbusment.ReimbursementEntity{}, errors.New("error get batasan reimbursment")
	}
	output := ModelToEntity(inputModel)
	return output, nil
}

// Update implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) Update(input reimbusment.ReimbursementEntity, id string) error {
	inputModel := EntityToModel(input)
	tx := repo.db.Model(&Reimbursement{}).Where("id=?", id).Updates(inputModel)
	if tx.Error != nil {
		return errors.New("update data reimbursment")
	}
	fmt.Println("input admin", input)
	fmt.Println("input model", inputModel)
	if tx.RowsAffected == 0 {
		return errors.New("row not affected")
	}
	return nil
}

// Insert implements reimbusment.ReimbusmentDataInterface.
func (repo *ReimbusmentData) Insert(input reimbusment.ReimbursementEntity) error {
	idUUID, errUUID := helper.GenerateUUID()
	if errUUID != nil {
		return errors.New("failed generated uuid")
	}
	inputModel := EntityToModel(input)
	inputModel.ID = idUUID
	tx := repo.db.Create(&inputModel)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("row not affected")
	}
	return nil
}

func New(db *gorm.DB) reimbusment.ReimbusmentDataInterface {
	return &ReimbusmentData{
		db: db,
	}
}
