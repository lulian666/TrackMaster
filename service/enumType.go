package service

import (
	"TrackMaster/model"
	"TrackMaster/third_party/jet"
	"gorm.io/gorm"
)

type EnumTypeService interface {
	SyncEnumType(p *model.Project) error
}

type enumTypeService struct {
	DB *gorm.DB
}

func NewEnumTypeService(db *gorm.DB) EnumTypeService {
	return &enumTypeService{DB: db}
}

func (s enumTypeService) SyncEnumType(p *model.Project) error {
	// project 是否存在
	err := p.Get(s.DB)
	if err != nil {
		return err
	}

	types, err := jet.GetEnumTypes(p.ID)
	if err != nil {
		return err
	}

	typeIDs := make([]string, len(types))
	for i := range types {
		typeIDs[i] = types[i].ID
	}

	var existingTypes []model.Type
	s.DB.Where("id IN (?)", typeIDs).Find(&existingTypes)

	// sync type
	for i := range types {
		var t model.Type
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
			err = t.Create(s.DB)
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
		s.DB.Where("type_id = ?", t.ID).Where("id IN (?)", enumValueIDs).Find(&existingEnumValues)

		createList := make([]model.EnumValue, 0, len(enumValues))
		updateList := make([]model.EnumValue, 0, len(enumValues))

		for m := range enumValues {
			var ev model.EnumValue
			for n := range existingEnumValues {
				if existingEnumValues[n].ID == enumValues[m].ID {
					ev.ID = existingEnumValues[n].ID
					ev.TypeId = t.ID
					// 需要更新的记录
					if ev.Name != existingEnumValues[n].Name {
						ev.Name = existingEnumValues[n].Name
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
			result := s.DB.Save(updateList)
			if result.Error != nil {
				return err
			}
		}

		// 批量创建
		if len(createList) > 0 {
			result := s.DB.Create(createList)
			if result.Error != nil {
				return err
			}
		}
	}

	return nil
}
