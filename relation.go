package model

type RelationType int

const (
	RELATION_ONE_TO_ONE RelationType = iota
	RELATION_ONE_TO_MANY
	RELATION_MANY_TO_ONE
	RELATION_MANY_TO_MANY
)

type Relation struct {
	ExtModel                 IModel
	RelationType             RelationType
	LocalFieldsNames         []string
	FkFieldsNames            []string
	JunctionModel            IModel
	JunctionLocalFieldsNames []string
	JunctionFkFieldsNames    []string
	IsRequired               bool
	IsBack                   bool
}

func (r RelationType) String() string {
	switch r {
	case RELATION_ONE_TO_ONE:
		return "OneToOne"
	case RELATION_ONE_TO_MANY:
		return "OneToMany"
	case RELATION_MANY_TO_ONE:
		return "ManyToOne"
	case RELATION_MANY_TO_MANY:
		return "ManyToMany"
	default:
		return "Unknown"
	}
}

func AddOneToOneRelation(model1, model2 IModel, required bool) {
	model1.AddRelation(Relation{
		ExtModel:         model2,
		RelationType:     RELATION_ONE_TO_ONE,
		LocalFieldsNames: model1.GetPKFieldsNames(),
		FkFieldsNames:    model2.GetPKFieldsNames(),
		IsRequired:       required,
	}, nil)

	model2.AddRelation(Relation{
		ExtModel:         model1,
		RelationType:     RELATION_ONE_TO_ONE,
		LocalFieldsNames: model2.GetPKFieldsNames(),
		FkFieldsNames:    model1.GetPKFieldsNames(),
		IsRequired:       true,
		IsBack:           true,
	}, nil)
}

func AddManyToOneRelation(model1, model2 IModel, required bool) {
	fkFieldsNames := make([]string, len(model2.GetPKFieldsNames()))
	fkFields := make([]IFieldDefinition, len(fkFieldsNames))
	for i, pkFieldName := range model2.GetPKFieldsNames() {
		fkName := "fk_" + model2.GetId() + "_" + pkFieldName
		fkFieldsNames[i] = fkName
		fkFields[i] = model2.GetFieldDefinition(pkFieldName).CloneForFK(fkName, "FK field", required)
	}

	model1.AddRelation(Relation{
		ExtModel:         model2,
		RelationType:     RELATION_MANY_TO_ONE,
		LocalFieldsNames: fkFieldsNames,
		FkFieldsNames:    model1.GetPKFieldsNames(),
		IsRequired:       required,
	}, fkFields)

	model2.AddRelation(Relation{
		ExtModel:         model1,
		RelationType:     RELATION_ONE_TO_MANY,
		LocalFieldsNames: model1.GetPKFieldsNames(),
		FkFieldsNames:    fkFieldsNames,
		IsBack:           true,
	}, nil)
}

func AddManyToManyRelation(model1, model2 IModel, storage IStorage) {
	junctionPkFields := make([]string, 0, len(model1.GetPKFieldsNames())+len(model2.GetPKFieldsNames()))
	junctionFields := make([]IFieldDefinition, 0, len(junctionPkFields))

	fk1Fields := make([]string, len(model1.GetPKFieldsNames()))
	for i, pkFieldName := range model1.GetPKFieldsNames() {
		fk1Fields[i] = "fk_" + model1.GetId() + "_" + pkFieldName
		junctionPkFields = append(junctionPkFields, fk1Fields[i])
		junctionFields = append(junctionFields, model1.GetFieldDefinition(pkFieldName).CloneForFK(fk1Fields[i], "FK field", true))
	}

	fk2Fields := make([]string, len(model2.GetPKFieldsNames()))
	for i, pkFieldName := range model2.GetPKFieldsNames() {
		fk2Fields[i] = "fk_" + model2.GetId() + "_" + pkFieldName
		junctionPkFields = append(junctionPkFields, fk2Fields[i])
		junctionFields = append(junctionFields, model2.GetFieldDefinition(pkFieldName).CloneForFK(fk2Fields[i], "FK field", true))
	}

	junctionModel := storage.NewModel("_junction__"+model1.GetId()+"__"+model2.GetId(), junctionFields, junctionPkFields)

	model1.AddRelation(Relation{
		ExtModel:                 model2,
		RelationType:             RELATION_MANY_TO_MANY,
		LocalFieldsNames:         model1.GetPKFieldsNames(),
		FkFieldsNames:            model2.GetPKFieldsNames(),
		JunctionModel:            junctionModel,
		JunctionLocalFieldsNames: fk1Fields,
		JunctionFkFieldsNames:    fk2Fields,
	}, nil)

	model2.AddRelation(Relation{
		ExtModel:                 model1,
		RelationType:             RELATION_MANY_TO_MANY,
		LocalFieldsNames:         model2.GetPKFieldsNames(),
		FkFieldsNames:            model1.GetPKFieldsNames(),
		JunctionModel:            junctionModel,
		JunctionLocalFieldsNames: fk2Fields,
		JunctionFkFieldsNames:    fk1Fields,
		IsBack:                   true,
	}, nil)

	junctionModel.AddRelation(Relation{
		ExtModel:         model1,
		RelationType:     RELATION_MANY_TO_ONE,
		LocalFieldsNames: fk1Fields,
		FkFieldsNames:    model1.GetPKFieldsNames(),
	}, nil)

	junctionModel.AddRelation(Relation{
		ExtModel:         model2,
		RelationType:     RELATION_MANY_TO_ONE,
		LocalFieldsNames: fk2Fields,
		FkFieldsNames:    model2.GetPKFieldsNames(),
	}, nil)
}
