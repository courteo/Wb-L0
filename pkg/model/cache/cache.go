package cache

import model "withNats/pkg/model"

type NatsDataCacheRepository struct {
	data   []model.NatsData
}


func (repo *NatsDataCacheRepository) FindNatsData(ID string) (*model.NatsData, error) {
	for _, val := range repo.data {
		if val.OrderUID == ID {
			return &val, nil
		}
	}
	return nil, model.ErrNoNatsData
}

func (repo *NatsDataCacheRepository) Add(nData *model.NatsData) error {
	repo.data = append(repo.data, *nData)

	return nil
}