package track

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

func SyncEnumType(p *model.Project) *pkg.Error {
	types, err := jet.GetEnumTypes(p.ID)
	if err != nil {
		return err
	}

	typeIDs := make([]string, len(types))
	for i := range types {
		typeIDs[i] = types[i].ID
	}

	var existingTypes []model.Type
	initializer.DB.Where("id IN (?)", typeIDs).Find(&existingTypes)

	// sync type
	for i := range types {
		t := model.Type{}
		for j := range existingTypes {
			if existingTypes[j].ID == types[i].ID {
				t = existingTypes[j]
			}
		}

		// create new ones
		// note: normally we only have 3 types, and they're always the same
		// so here we don't consider updating
		if t.ID == "" {
			t.ID = types[i].ID
			t.ProjectID = p.ID
			t.Type = types[i].Name
			err = t.Create(initializer.DB)
			if err != nil {
				return err
			}
		}

		// sync enum value in this type (types[i] / t)
		enumValues := types[i].EnumValues
		enumValueIDs := make([]string, len(enumValues))
		for k := range enumValues {
			enumValueIDs[k] = enumValues[k].ID
		}

		var existingEnumValues []model.EnumValue
		initializer.DB.Where("type_id = ?", t.ID).Where("id IN (?)", enumValueIDs).Find(&existingEnumValues)

		createList := make([]model.EnumValue, 0, len(enumValues))
		updateList := make([]model.EnumValue, 0, len(enumValues))

		for m := range enumValues {
			ev := model.EnumValue{}
			for n := range existingEnumValues {
				if existingEnumValues[n].ID == enumValues[m].ID {
					ev.ID = existingEnumValues[n].ID
					ev.TypeId = t.ID
					// 需要更新的记录
					if enumValues[m].Name != existingEnumValues[n].Name {
						ev.Name = enumValues[m].Name
						updateList = append(updateList, ev)
					}
				}
			}

			// 需要创建的记录
			if ev.ID == "" {
				ev.TypeId = t.ID
				ev.ID = enumValues[m].ID
				ev.Name = enumValues[m].Name
				createList = append(createList, ev)
			}
		}

		// 批量更新
		if len(updateList) > 0 {
			result := initializer.DB.Save(updateList)
			if result.Error != nil {
				return pkg.NewError(pkg.ServerError, result.Error.Error())
			}
		}

		// 批量创建
		if len(createList) > 0 {
			result := initializer.DB.Create(createList)
			if result.Error != nil {
				return pkg.NewError(pkg.ServerError, result.Error.Error())
			}
		}
	}

	return nil
}

func LocateValue(field jet.Field, db *gorm.DB) ([]string, *pkg.Error) {
	if len(field.Values) > 0 {
		// 根据type id和id去拿值
		values := make([]string, 0, len(field.Values))
		for _, v := range field.Values {
			value := model.EnumValue{}
			result := db.Where("type_id = ?", field.Type.ID).Where("id = ?", v).Find(&value)
			if result.Error != nil {
				return nil, pkg.NewError(pkg.ServerError, result.Error.Error())
			}
			if value.ID != "" {
				values = append(values, value.Name)
			}
		}
		// 如果遍历完数组里的值后，一个value都没有找到，就说明里面的值并非id
		// 那里面的值就是我们要取的值本身，直接赋值就好了
		if len(values) == 0 {
			values = append(values, field.Values...)
		}
		return values, nil
	} else {
		// 没有值
		return []string{""}, nil
	}
}
