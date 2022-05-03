/**
* @author mnunez
 */

package gosql

//var dfltBuilder = QueryBuilder{}

type Repository interface {
	GetByID(id string, models interface{}) (interface{}, error)
	DeleteByID(id string, models interface{}) error
}

func GetAll(models interface{}) (interface{}, error) {
	err := data.DB.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func Create(models interface{}) error {
	err := data.DB.Create(models).Error
	if err != nil {
		return err
	}
	return nil
}

func GetByID(id string, models interface{}) (interface{}, error) {
	err := data.DB.Where("id = ?", id).Find(models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func DeleteByID(id string, models interface{}) error {
	err := data.DB.Delete(models, id).Error
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

func rawQueryBuild(query string, models interface{}) (interface{}, error) {
	err := data.DB.Raw(query).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func execQueryBuild(query string) error {
	err := data.DB.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}
