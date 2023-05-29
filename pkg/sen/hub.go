package sen

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	autoInjectionTag = "*"
	optionalTag      = "optional"
	injectTag        = "inject"
)

// ErrComponentNotRegistered is returned when the expected component isn't registered
// so it couldn't be found by name.
var ErrComponentNotRegistered = errors.New("sen: the component is not registered")

// Hub is a container of components.
// It allows registering new components by names as well as
// injecting dependencies into a component via tags or types.
type Hub interface {
	// Register injects dependencies into a component and register the component into the depdenency container
	// for the next injection.
	Register(name string, component interface{}) error

	// Retrieve retrieves a component via name. It returns an error if there is any.
	Retrieve(name string) (interface{}, error)

	// Inject injects dependencies into a component.
	Inject(component interface{}) error
}

func newHub() Hub {
	hub := &defaultHub{
		dependencies: make(map[string]*dependency),
	}

	_ = hub.Register("hub", hub)
	return hub
}

type dependency struct {
	value        interface{}
	reflectValue reflect.Value
	reflectType  reflect.Type
}

type defaultHub struct {
	dependencies map[string]*dependency
}

func (hub *defaultHub) Register(name string, component interface{}) error {
	if err := hub.validateNamne(name); err != nil {
		return err
	}

	toAddDep := &dependency{
		value:        component,
		reflectType:  reflect.TypeOf(component),
		reflectValue: reflect.ValueOf(component),
	}

	if err := hub.inject(toAddDep); err != nil {
		return err
	}

	hub.dependencies[name] = toAddDep

	return nil
}

func (hub *defaultHub) Retrieve(name string) (interface{}, error) {
	loadedDep, found := hub.dependencies[name]
	if !found {
		return nil, ErrComponentNotRegistered
	}

	return loadedDep.value, nil
}

func (hub *defaultHub) Inject(component interface{}) error {
	toAddDep := &dependency{
		value:        component,
		reflectType:  reflect.TypeOf(component),
		reflectValue: reflect.ValueOf(component),
	}

	return hub.inject(toAddDep)
}

func (hub *defaultHub) validateNamne(name string) error {
	if _, found := hub.dependencies[name]; found {
		return fmt.Errorf("hub: %s is already registered", name)
	}

	if name == autoInjectionTag {
		return fmt.Errorf("hub: %s is revserved, please use a different name", autoInjectionTag)
	}

	return nil
}

func (hub *defaultHub) inject(dep *dependency) error {
	if !isStructPtr(dep.reflectType) {
		if hasInjectTag(dep) {
			return fmt.Errorf("hub: %s is not injectable, a pointer is expected", dep.reflectType)
		}

		return nil
	}

	for i := 0; i < dep.reflectValue.Elem().NumField(); i++ {
		fieldValue := dep.reflectValue.Elem().Field(i)
		fieldType := fieldValue.Type()
		structField := dep.reflectType.Elem().Field(i)
		fieldTag := structField.Tag
		tagValue, ok := fieldTag.Lookup(injectTag)
		if !ok {
			continue
		}

		loadedDep, err := hub.loadDepForTag(tagValue, fieldType)
		if err != nil {
			return err
		}

		if loadedDep == nil {
			// this is an optional field and there is no suitable dependency to inject.
			continue
		}

		if !loadedDep.reflectType.AssignableTo(fieldType) {
			return fmt.Errorf("hub: %s is not assignable from %s", fieldType, loadedDep.reflectType)
		}

		fieldValue.Set(loadedDep.reflectValue)
	}

	return nil
}

func (hub *defaultHub) loadDepForTag(tag string, t reflect.Type) (*dependency, error) {
	tagName, optional, err := parseTag(tag)
	if err != nil {
		return nil, err
	}

	if tag == autoInjectionTag {
		return hub.findByType(t, optional)
	}

	loadedDep, found := hub.dependencies[tagName]
	if !found && !optional {
		return nil, fmt.Errorf("hub: %s is not registered", tagName)
	}

	return loadedDep, nil
}

func (hub *defaultHub) findByType(t reflect.Type, optional bool) (*dependency, error) {
	var foundVal *dependency
	for _, v := range hub.dependencies {
		if v.reflectType.AssignableTo(t) {
			if foundVal != nil {
				return nil, fmt.Errorf("hub: there is a conflict when finding the dependency for %s", t.String())
			}

			foundVal = v
		}
	}

	if foundVal == nil && !optional {
		return nil, fmt.Errorf("hub: couldn't find the dependency for %s", t.String())
	}

	return foundVal, nil
}

func parseTag(tag string) (string, bool, error) {
	parts := strings.Split(tag, ",")
	switch len(parts) {
	case 2:
		if parts[1] != optionalTag {
			return "", false, fmt.Errorf("hub: %s is unexpected", parts[1])
		}

		return parts[0], parts[1] == optionalTag, nil
	case 1:
		return parts[0], false, nil
	case 0:
		return "", false, fmt.Errorf("hub: tag must not be empty")
	default:
		return "", false, fmt.Errorf("hub: unable to parse tag %s", tag)
	}
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func hasInjectTag(dep *dependency) bool {
	if dep.reflectType.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < dep.reflectType.NumField(); i++ {
		structField := dep.reflectType.Field(i)
		if _, ok := structField.Tag.Lookup(injectTag); ok {
			return true
		}
	}

	return false
}
