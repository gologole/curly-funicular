package service

import (
	"cmd/main.go/configs"
	"cmd/main.go/internal/cache"
	"cmd/main.go/internal/db"
	"cmd/main.go/models"
	"context"
	"encoding/json"
)

type Service interface {
	Getfile(name string, userid string) (models.File, error)
	UploadFile(file models.File) error
	DeleteFile(name string, userid string) error
	GetFilesByUserID(userid string) ([]models.SimpleFileView, error)
}

type service struct {
	cfg   *configs.Config
	db    db.Database
	cache cache.Cache
}

func NewService(cfg *configs.Config, db db.Database, cache cache.Cache) Service {
	return &service{cfg: cfg,
		db:    db,
		cache: cache}
}

func (s *service) Getfile(name string, userid string) (models.File, error) {

	file, is, err := s.cache.Get(context.Background(), name+userid)
	var filestruct models.File

	json.Unmarshal([]byte(file), &filestruct)

	if !is {
		filestruct, err = s.db.GetFile(name, userid)
		if err != nil {
			return models.File{}, err
		}
		return filestruct, nil
	}

	go func() {
		v, _ := json.Marshal(filestruct)
		err := s.cache.Set(context.Background(), name+userid, string(v), 500)
		if err != nil {
			//s.ELKapi.SendData()
		}
	}()
	//s.ELKapi.SendData()
	return filestruct, nil
}

func (s *service) UploadFile(file models.File) error {
	file.Size = len(file.Data)
	err := s.db.UploadFile(file)
	return err
	//s.ELKapi.SendData()
}

func (s *service) DeleteFile(name string, userid string) error {
	s.cache.Delete(context.Background(), name+userid)
	return s.db.DeleteFile(name, userid)

}

func (s *service) GetFilesByUserID(userid string) ([]models.SimpleFileView, error) {
	return s.db.GetFilesByUserID(userid)
}
