package repo

import (
	"database/sql"
	"encoding/json"
	model "withNats/pkg/model"
	st "withNats/pkg/structs"
)



type NatsDataMemoryRepository struct {
	Data   *sql.DB
}


func NewMemoryRepo(db *sql.DB) *NatsDataMemoryRepository {
	return &NatsDataMemoryRepository{
		Data:   db,
	}
}

func (repo *NatsDataMemoryRepository) FindNatsData(ID string) (*model.NatsData, error) {
	nData := &model.NatsData{}
	row := repo.Data.QueryRow(`SELECT * FROM orders WHERE OrderUID  = $1`, ID)
	err := row.Scan(&nData)
	if err != nil {
		return nil, model.ErrNoNatsData
	}
	return nData, nil
}

func (repo *NatsDataMemoryRepository) Add(nData *model.NatsData) error {
	byte, err := json.Marshal(nData)
	data := st.OrderJSON{Order_uid: nData.OrderUID, DataJSON: string(byte)}
	_, err = repo.Data.Exec( `insert into orders (order_id, data) values ($1,$2)`, data.Order_uid, data.DataJSON)
	return err
}

