/**
* @author mnunez
 */

package gosql

//var dfltBuilder = QueryBuilder{}

func GetAll(ret []interface{}) ([]interface{}, error) {
	err := data.DB.Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func Create(ret *interface{}) error {
	err := data.DB.Create(&ret).Error
	if err != nil {
		return err
	}
	return nil
}

func GetByID(id string, ret interface{}) (interface{}, error) {
	err := data.DB.Where("id = ?", id).First(&ret).Error
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func DeleteByID(id string, ret interface{}) error {
	err := data.DB.Delete(ret, id).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateByID(retUpdate *interface{}, id int) (interface{}, error) {
	if err := data.DB.Where("id = ?", id).First(&retUpdate).Error; err != nil {
		return retUpdate, err
	}
	if err := data.DB.Model(&retUpdate).Updates(retUpdate).Error; err != nil {
		return retUpdate, err
	}
	return retUpdate, nil
}
