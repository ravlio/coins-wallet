package db

import "github.com/jackc/pgtype"

type Model interface {
	Values(cols ...string) []interface{}
}

type ModelProp struct {
	prop  pgtype.Value
	value interface{}
}
type ModelSetter struct {
	props []ModelProp
}

func NewModelSetter() *ModelSetter {
	return &ModelSetter{props: make([]ModelProp, 0)}
}

func (ms *ModelSetter) Set(prop pgtype.Value, value interface{}) {
	ms.props = append(ms.props, ModelProp{prop: prop, value: value})
}

func (ms *ModelSetter) Apply() error {
	var err error

	for _, v := range ms.props {
		err = v.prop.Set(v)
		if err != nil {
			return err
		}
	}

	return nil
}
