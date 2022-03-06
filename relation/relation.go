package relation

import "github.com/go-qbit/model"

type relationOpts struct {
	required         bool
	alias, backAlias string
}

type relationOptsFunc func(opts *relationOpts)

func WithRequired(required bool) relationOptsFunc {
	return func(opts *relationOpts) {
		opts.required = required
	}
}

func WithAlias(alias string) relationOptsFunc {
	return func(opts *relationOpts) {
		opts.alias = alias
	}
}

func WithBackAlias(alias string) relationOptsFunc {
	return func(opts *relationOpts) {
		opts.backAlias = alias
	}
}

func AddOneToOne(model1, model2 model.IModel, opts ...relationOptsFunc) {
	o := &relationOpts{}
	for _, optFunc := range opts {
		optFunc(o)
	}

	model1.AddRelation(model.Relation{
		ExtModel:         model2,
		RelationType:     model.RELATION_ONE_TO_ONE,
		LocalFieldsNames: model1.GetPKFieldsNames(),
		FkFieldsNames:    model2.GetPKFieldsNames(),
		IsRequired:       o.required,
	}, "", nil)

	model2.AddRelation(model.Relation{
		ExtModel:         model1,
		RelationType:     model.RELATION_ONE_TO_ONE,
		LocalFieldsNames: model2.GetPKFieldsNames(),
		FkFieldsNames:    model1.GetPKFieldsNames(),
		IsRequired:       true,
		IsBack:           true,
	}, "", nil)
}

func AddManyToOne(model1, model2 model.IModel, opts ...relationOptsFunc) {
	o := &relationOpts{}
	for _, optFunc := range opts {
		optFunc(o)
	}

	fkFieldsNames := make([]string, len(model2.GetPKFieldsNames()))
	fkFields := make([]model.IFieldDefinition, len(fkFieldsNames))
	for i, pkFieldName := range model2.GetPKFieldsNames() {
		fkName := "fk_"
		if o.alias != "" {
			fkName += o.alias
		} else {
			fkName += model2.GetId()
		}

		fkName += "_" + pkFieldName

		fkFieldsNames[i] = fkName
		fkFields[i] = model2.GetFieldDefinition(pkFieldName).CloneForFK(fkName, "FK field", o.required)
	}

	model1.AddRelation(model.Relation{
		ExtModel:         model2,
		RelationType:     model.RELATION_MANY_TO_ONE,
		LocalFieldsNames: fkFieldsNames,
		FkFieldsNames:    model2.GetPKFieldsNames(),
		IsRequired:       o.required,
	}, o.alias, fkFields)

	model2.AddRelation(model.Relation{
		ExtModel:         model1,
		RelationType:     model.RELATION_ONE_TO_MANY,
		LocalFieldsNames: model1.GetPKFieldsNames(),
		FkFieldsNames:    fkFieldsNames,
		IsBack:           true,
	}, o.backAlias, nil)
}

func AddManyToMany(model1, model2 model.IModel, storage model.IStorage, opts ...relationOptsFunc) {
	o := &relationOpts{}
	for _, optFunc := range opts {
		optFunc(o)
	}

	junctionPkFields := make([]string, 0, len(model1.GetPKFieldsNames())+len(model2.GetPKFieldsNames()))
	junctionFields := make([]model.IFieldDefinition, 0, len(junctionPkFields))

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

	junctionModel := storage.NewModel("_junction__"+model1.GetId()+"__"+model2.GetId(), junctionFields, model.BaseModelOpts{
		PkFieldsNames: junctionPkFields,
	})

	model1.AddRelation(model.Relation{
		ExtModel:                 model2,
		RelationType:             model.RELATION_MANY_TO_MANY,
		IsRequired:               true,
		LocalFieldsNames:         model1.GetPKFieldsNames(),
		FkFieldsNames:            model2.GetPKFieldsNames(),
		JunctionModel:            junctionModel,
		JunctionLocalFieldsNames: fk1Fields,
		JunctionFkFieldsNames:    fk2Fields,
	}, "", nil)

	model2.AddRelation(model.Relation{
		ExtModel:                 model1,
		RelationType:             model.RELATION_MANY_TO_MANY,
		IsRequired:               true,
		LocalFieldsNames:         model2.GetPKFieldsNames(),
		FkFieldsNames:            model1.GetPKFieldsNames(),
		JunctionModel:            junctionModel,
		JunctionLocalFieldsNames: fk2Fields,
		JunctionFkFieldsNames:    fk1Fields,
		IsBack:                   true,
	}, "", nil)

	junctionModel.AddRelation(model.Relation{
		ExtModel:         model1,
		RelationType:     model.RELATION_MANY_TO_ONE,
		IsRequired:       true,
		LocalFieldsNames: fk1Fields,
		FkFieldsNames:    model1.GetPKFieldsNames(),
	}, "", nil)

	junctionModel.AddRelation(model.Relation{
		ExtModel:         model2,
		RelationType:     model.RELATION_MANY_TO_ONE,
		IsRequired:       true,
		LocalFieldsNames: fk2Fields,
		FkFieldsNames:    model2.GetPKFieldsNames(),
	}, "", nil)
}

func AddManyToManyUsingTable(model1, model2, junction model.IModel, opts ...relationOptsFunc) {
	o := &relationOpts{}
	for _, optFunc := range opts {
		optFunc(o)
	}

	junctionPkFields := make([]string, 0, len(model1.GetPKFieldsNames())+len(model2.GetPKFieldsNames()))
	junctionFields := make([]model.IFieldDefinition, 0, len(junctionPkFields))

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

	//model1.AddRelation(model.Relation{
	//	ExtModel:                 model2,
	//	RelationType:             model.RELATION_MANY_TO_MANY,
	//	IsRequired:               true,
	//	LocalFieldsNames:         model1.GetPKFieldsNames(),
	//	FkFieldsNames:            model2.GetPKFieldsNames(),
	//	JunctionModel:            junction,
	//	JunctionLocalFieldsNames: fk1Fields,
	//	JunctionFkFieldsNames:    fk2Fields,
	//}, "", nil)
	//
	//model2.AddRelation(model.Relation{
	//	ExtModel:                 model1,
	//	RelationType:             model.RELATION_MANY_TO_MANY,
	//	IsRequired:               true,
	//	LocalFieldsNames:         model2.GetPKFieldsNames(),
	//	FkFieldsNames:            model1.GetPKFieldsNames(),
	//	JunctionModel:            junction,
	//	JunctionLocalFieldsNames: fk2Fields,
	//	JunctionFkFieldsNames:    fk1Fields,
	//	IsBack:                   true,
	//}, "", nil)

	junction.AddRelation(model.Relation{
		ExtModel:         model1,
		RelationType:     model.RELATION_MANY_TO_ONE,
		IsRequired:       true,
		PkFieldsNames:    junctionPkFields,
		LocalFieldsNames: fk1Fields,
		FkFieldsNames:    model1.GetPKFieldsNames(),
	}, "", junctionFields)

	junction.AddRelation(model.Relation{
		ExtModel:         model2,
		RelationType:     model.RELATION_MANY_TO_ONE,
		IsRequired:       true,
		LocalFieldsNames: fk2Fields,
		FkFieldsNames:    model2.GetPKFieldsNames(),
	}, "", nil)
}
