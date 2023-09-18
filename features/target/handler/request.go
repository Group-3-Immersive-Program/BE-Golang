package handler

import (
	"be_golang/klp3/features/target"
)

type TargetRequest struct {
	KontenTarget   string `json:"konten_target,omitempty" form:"konten_target"`
	Status         string `json:"status,omitempty" form:"status"`
	DevisiID       string `json:"devisi_id,omitempty" form:"devisi_id"`
	UserIDPembuat  string `json:"user_id_pembuat,omitempty" form:"user_id_pembuat"`
	UserIDPenerima string `json:"user_id_penerima,omitempty" form:"user_id_penerima"`
	Due_Date       string `json:"due_date,omitempty" form:"due_date"`
	Proofs         string `json:"proofs,omitempty" form:"proofs"`
	// User           UserEntity
}

type TargetReqPenerima struct {
	Status string `json:"status,omitempty" form:"status"`
	Proofs string `json:"proofs,omitempty" form:"proofs"`
}

func TargetRequestToEntity(req TargetRequest) target.TargetEntity {
	return target.TargetEntity{
		KontenTarget:   req.KontenTarget,
		Status:         req.Status,
		DevisiID:       req.DevisiID,
		UserIDPembuat:  req.UserIDPembuat,
		UserIDPenerima: req.UserIDPenerima,
		Due_Date:       req.Due_Date,
		Proofs:         req.Proofs,
	}
}

func TargetReqPenerimaToEntity(req TargetReqPenerima) target.TargetEntity {
	return target.TargetEntity{
		Status: req.Status,
		Proofs: req.Proofs,
	}
}
