/**
* @author mnunez
 */

package gosql

//var dfltBuilder = QueryBuilder{}

type Repository interface {
	GetByID(storage Data, id string, models interface{}) (interface{}, error)
	DeleteByID(storage Data, id string, models interface{}) error
	GetAll(storage Data, models interface{}) (interface{}, error)
	Create(storage Data, models interface{}) error
}

func GetAll(storage Data, models interface{}) (interface{}, error) {
	err := storage.DB.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func Create(storage Data, models interface{}) error {
	err := storage.DB.Create(models).Error
	if err != nil {
		return err
	}
	return nil
}

func GetByID(storage Data, id string, models interface{}) (interface{}, error) {
	err := storage.DB.Where("id = ?", id).Find(models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func DeleteByID(storage Data, id string, models interface{}) error {
	err := storage.DB.Delete(models, id).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateByID(modelsUpdate interface{}, id int) (interface{}, error) {
	if err := data.DB.Where("id = ?", id).First(modelsUpdate).Error; err != nil {
		return modelsUpdate, err
	}
	if err := data.DB.Model(modelsUpdate).Updates(modelsUpdate).Error; err != nil {
		return modelsUpdate, err
	}
	return modelsUpdate, nil
}

func rawQueryBuild(storage Data, query string, models interface{}) (interface{}, error) {
	err := storage.DB.Raw(query).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func execQueryBuild(storage Data, query string) error {
	err := storage.DB.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}
