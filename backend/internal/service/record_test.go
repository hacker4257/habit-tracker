package service

import (
	"testing"

	"habit-tracker/internal/model"
)

type mockRepository struct {
	records []model.Record
	nextID  int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		records: []model.Record{},
		nextID:  1,
	}
}

func (m *mockRepository) Create(record *model.Record) error {
	record.ID = m.nextID
	m.nextID++
	m.records = append(m.records, *record)
	return nil
}

func (m *mockRepository) GetByID(id int64) (*model.Record, error) {
	for _, r := range m.records {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) GetAll() ([]model.Record, error) {
	return m.records, nil
}

func (m *mockRepository) Update(record *model.Record) error {
	for i, r := range m.records {
		if r.ID == record.ID {
			m.records[i] = *record
			return nil
		}
	}
	return nil
}

func (m *mockRepository) Delete(id int64) error {
	for i, r := range m.records {
		if r.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockRepository) GetStats() (*model.Stats, error) {
	total := 0
	for _, r := range m.records {
		total += r.Duration
	}
	return &model.Stats{
		TotalRecords:  len(m.records),
		TotalDuration: total,
	}, nil
}

func TestRecordService_Create(t *testing.T) {
	repo := newMockRepository()
	svc := NewRecordService(repo)

	tests := []struct {
		name    string
		req     *model.CreateRecordRequest
		wantErr bool
	}{
		{
			name: "valid record",
			req: &model.CreateRecordRequest{
				Date:     "2024-01-15",
				Content:  "Test content",
				Duration: 30,
				Notes:    "Test notes",
			},
			wantErr: false,
		},
		{
			name: "missing date",
			req: &model.CreateRecordRequest{
				Content:  "Test content",
				Duration: 30,
			},
			wantErr: true,
		},
		{
			name: "missing content",
			req: &model.CreateRecordRequest{
				Date:     "2024-01-15",
				Duration: 30,
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			req: &model.CreateRecordRequest{
				Date:     "2024-01-15",
				Content:  "Test content",
				Duration: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := svc.Create(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && record == nil {
				t.Error("Create() returned nil record")
			}
		})
	}
}

func TestRecordService_GetAll(t *testing.T) {
	repo := newMockRepository()
	svc := NewRecordService(repo)

	// Create some records
	svc.Create(&model.CreateRecordRequest{
		Date:     "2024-01-15",
		Content:  "Test 1",
		Duration: 30,
	})
	svc.Create(&model.CreateRecordRequest{
		Date:     "2024-01-16",
		Content:  "Test 2",
		Duration: 45,
	})

	records, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(records) != 2 {
		t.Errorf("GetAll() returned %d records, want 2", len(records))
	}
}

func TestRecordService_GetStats(t *testing.T) {
	repo := newMockRepository()
	svc := NewRecordService(repo)

	svc.Create(&model.CreateRecordRequest{
		Date:     "2024-01-15",
		Content:  "Test 1",
		Duration: 30,
	})
	svc.Create(&model.CreateRecordRequest{
		Date:     "2024-01-16",
		Content:  "Test 2",
		Duration: 45,
	})

	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats.TotalRecords != 2 {
		t.Errorf("GetStats() TotalRecords = %d, want 2", stats.TotalRecords)
	}

	if stats.TotalDuration != 75 {
		t.Errorf("GetStats() TotalDuration = %d, want 75", stats.TotalDuration)
	}
}
