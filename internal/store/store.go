package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
)

type Record map[string]interface{}

type Spec struct {
	Id      string
	Columns map[string]string
}

type Data struct {
	Id      int64
	Content Record
}

type Store struct {
	Pool *pgxpool.Pool
}

func (s *Store) GetSpec(ctx context.Context, id string) (*Spec, error) {
	spec := &Spec{
		Id:      id,
		Columns: make(map[string]string),
	}

	rows, err := s.Pool.Query(ctx, "select column_name, data_type from spec where id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		var dataType string

		if err := rows.Scan(&columnName, &dataType); err != nil {
			return nil, err
		}

		spec.Columns[columnName] = dataType
	}

	return spec, nil
}

func (s *Store) Create(ctx context.Context, specId string, record Record) (*Data, error) {
	var id int64
	if err := s.Pool.QueryRow(ctx, "insert into data(spec_id, content) values($1, $2) returning id", specId, record).Scan(&id); err != nil {
		return nil, err
	}

	return &Data{
		Id:      id,
		Content: record,
	}, nil
}

func (s *Store) Update(ctx context.Context, id int64, record Record) (*Data, error) {
	if _, err := s.Pool.Exec(ctx, "update data set content=$2 where id=$1", id, record); err != nil {
		return nil, err
	}

	return &Data{
		Id:      id,
		Content: record,
	}, nil
}

func (s *Store) Get(ctx context.Context, specId string, key string, value string) ([]*Data, error) {
	data := []*Data{}

	rows, err := s.Pool.Query(ctx, fmt.Sprintf("select id, content from data where spec_id=$1 and content->>'%s'=$2;", key), specId, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var record Record

		if err := rows.Scan(&id, &record); err != nil {
			return nil, err
		}

		data = append(data, &Data{
			Id:      id,
			Content: record,
		})
	}
	return data, nil
}

func (s *Store) Validate(ctx context.Context, specId string, record Record) error {
	spec, err := s.GetSpec(ctx, specId)
	if err != nil {
		return err
	}

	for key, value := range record {
		dataType := spec.Columns[key]
		if dataType == "" {
			return fmt.Errorf("data type not found for column %s", key)
		}

		switch dataType {
		case "TEXT":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("value %+v of type %s is not string", value, reflect.TypeOf(value))
			}
		case "INTEGER":
			// golang unmarshal generic json number to float64
			// so need to do more specific check to ensure it is int
			f, ok := value.(float64)
			if !ok {
				return fmt.Errorf("value %+v of type %s is not integer", value, reflect.TypeOf(value))
			}
			if f != float64(int64(f)) {
				return fmt.Errorf("value %+v is not inteer", value)
			}
		case "BOOLEAN":
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("value %+v of type %s is not boolean", value, reflect.TypeOf(value))
			}
		default:
			return fmt.Errorf("unsupported datatype %s", dataType)
		}
	}

	return nil
}
