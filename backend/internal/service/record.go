package service

import (
	"errors"

	"habit-tracker/internal/model"
	"habit-tracker/internal/repository"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input")
)

type RecordService interface {
	Create(req *model.CreateRecordRequest) (*model.Record, error)
	GetByID(id int64) (*model.Record, error)
	GetAll() ([]model.Record, error)
	Update(id int64, req *model.UpdateRecordRequest) (*model.Record, error)
	Delete(id int64) error
	GetStats() (*model.Stats, error)
}

type recordService struct {
	repo repository.RecordRepository
}

func NewRecordService(repo repository.RecordRepository) RecordService {
	return &recordService{repo: repo}
}

func (s *recordService) Create(req *model.CreateRecordRequest) (*model.Record, error) {
	if req.Date == "" || req.Content == "" || req.Duration < 1 {
		return nil, ErrInvalidInput
	}

	record := &model.Record{
		Date:     req.Date,
		Content:  req.Content,
		Duration: req.Duration,
		Notes:    req.Notes,
	}

	if err := s.repo.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *recordService) GetByID(id int64) (*model.Record, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}
	return record, nil
}

func (s *recordService) GetAll() ([]model.Record, error) {
	records, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if records == nil {
		return []model.Record{}, nil
	}
	return records, nil
}

func (s *recordService) Update(id int64, req *model.UpdateRecordRequest) (*model.Record, error) {
	if req.Date == "" || req.Content == "" || req.Duration < 1 {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrRecordNotFound
	}

	existing.Date = req.Date
	existing.Content = req.Content
	existing.Duration = req.Duration
	existing.Notes = req.Notes

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *recordService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return ErrRecordNotFound
	}
	return nil
}

func (s *recordService) GetStats() (*model.Stats, error) {
	return s.repo.GetStats()
}
