package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Piece represents a piece of the container.
type Piece struct {
	Key       string
	Singleton bool
	Value     interface{}
}

// Container contains all piece.
type Container map[string]*Piece

var instance Container
var once sync.Once

// Instance will return the container instance.
func Instance() Container {
	once.Do(func() {
		instance = Container{}
	})

	return instance
}

// Make will resolve the given identifier key and return the value
// if it exists.
func (c Container) Make(key string, params ...interface{}) (interface{}, error) {
	piece := c[key]

	if piece == nil {
		return nil, fmt.Errorf("Identifier `%s` does not exist", key)
	}

	value := piece.Value

	switch reflect.ValueOf(value).Kind() {
	case reflect.Func:
		typ := reflect.ValueOf(value).Type()
		names := getFuncNameParameters(typ.String())

		for index, val := range params {
			params[index] = reflect.ValueOf(val)
		}

		for _, name := range names {
			if !c.Contains(name) {
			}

			if val, _ := c.Make(name); val != nil {
				params = append(params, reflect.ValueOf(val))
			}
		}

		l := len(params)
		items := []reflect.Value{}

		for i := 0; i < typ.NumIn(); i++ {
			if i < l {
				items = append(items, params[i].(reflect.Value))
			} else {
				items = append(items, reflect.Value{})
			}
		}

		res := reflect.ValueOf(value).Call(items)

		return res[0].Interface(), nil
	case reflect.Struct:
	case reflect.Ptr:
		v := reflect.ValueOf(value)

		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		typ := v.Type()

		if v.Kind() != reflect.Struct {
			break
		}

		for i := 0; i < v.NumField(); i++ {
			if typ.Field(i).Tag == "inject" {
				f := v.Field(i)
				name := getFieldTypeName(f.String())

				if !c.Contains(name) {
					continue
				}

				if val, _ := c.Make(name); val != nil {
					f.Set(reflect.ValueOf(val))
				}
			}
		}

		break
	}

	return value, nil
}

// bindPiece will bind the given piece to the container.
func (c Container) bindPiece(p *Piece) error {
	if c == nil {
		return errors.New("Container is nil and cannot be used")
	}

	if c[p.Key] != nil && c[p.Key].Singleton {
		return fmt.Errorf("Identifier `%s` is a singleton and cannot be rebind", p.Key)
	}

	c[p.Key] = p

	return nil
}

// Bind will bind the identifier key with the given value. It can overwrite existing
// values with the same identifier.
func (c Container) Bind(key string, value interface{}) error {
	return c.bindPiece(&Piece{
		Key:       key,
		Singleton: false,
		Value:     value,
	})
}

// Singleton will bind identifier key with the given value if it don't exists.
// A singleton can only be bind once.
func (c Container) Singleton(key string, value interface{}) error {
	return c.bindPiece(&Piece{
		Key:       key,
		Singleton: true,
		Value:     value,
	})
}

// Contains will check if the given identifier
// key exists in the container or not.
func (c Container) Contains(key string) bool {
	return c[key] != nil
}

// All will return all identifier and values
// from the container.
func (c Container) All() Container {
	return c
}

// Remove value from the container by identifier key.
func (c Container) Remove(key string) {
	delete(c, key)
}
