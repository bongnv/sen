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

// Injector is a hub of components. It allows injecting components via tags or types.
type Injector interface {
	Register(name string, component interface{}) error
	Retrieve(name string) (interface{}, error)
	Inject(component interface{}) error
}

func newInjector() Injector {
	injector := &defaultInjector{
		dependencies: make(map[string]*dependency),
	}

	_ = injector.Register("injector", injector)
	return injector
}

type dependency struct {
	value        interface{}
	reflectValue reflect.Value
	reflectType  reflect.Type
}

type defaultInjector struct {
	dependencies map[string]*dependency
}

// Register injects dependencies into a component and register the component into the depdenency container
// for the next injection.
func (injector *defaultInjector) Register(name string, component interface{}) error {
	if err := injector.validateNamne(name); err != nil {
		return err
	}

	toAddDep := &dependency{
		value:        component,
		reflectType:  reflect.TypeOf(component),
		reflectValue: reflect.ValueOf(component),
	}

	if err := injector.inject(toAddDep); err != nil {
		return err
	}

	injector.dependencies[name] = toAddDep

	return nil
}

// Retrieve retrieves a component via name. It returns an error if there is any.
func (injector *defaultInjector) Retrieve(name string) (interface{}, error) {
	loadedDep, found := injector.dependencies[name]
	if !found {
		return nil, ErrComponentNotRegistered
	}

	return loadedDep.value, nil
}

// Inject injects dependencies into a component.
func (injector *defaultInjector) Inject(component interface{}) error {
	toAddDep := &dependency{
		value:        component,
		reflectType:  reflect.TypeOf(component),
		reflectValue: reflect.ValueOf(component),
	}

	return injector.inject(toAddDep)
}

func (injector *defaultInjector) validateNamne(name string) error {
	if _, found := injector.dependencies[name]; found {
		return fmt.Errorf("injector: %s is already registered", name)
	}

	if name == autoInjectionTag {
		return fmt.Errorf("injector: %s is revserved, please use a different name", autoInjectionTag)
	}

	return nil
}

func (injector *defaultInjector) inject(dep *dependency) error {
	if !isStructPtr(dep.reflectType) {
		if hasInjectTag(dep) {
			return fmt.Errorf("injector: %s is not injectable, a pointer is expected", dep.reflectType)
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

		loadedDep, err := injector.loadDepForTag(tagValue, fieldType)
		if err != nil {
			return err
		}

		if loadedDep == nil {
			// this is an optional field and there is no suitable dependency to inject.
			continue
		}

		if !loadedDep.reflectType.AssignableTo(fieldType) {
			return fmt.Errorf("injector: %s is not assignable from %s", fieldType, loadedDep.reflectType)
		}

		fieldValue.Set(loadedDep.reflectValue)
	}

	return nil
}

func (injector *defaultInjector) loadDepForTag(tag string, t reflect.Type) (*dependency, error) {
	tagName, optional, err := parseTag(tag)
	if err != nil {
		return nil, err
	}

	if tag == autoInjectionTag {
		return injector.findByType(t, optional)
	}

	loadedDep, found := injector.dependencies[tagName]
	if !found && !optional {
		return nil, fmt.Errorf("injector: %s is not registered", tagName)
	}

	return loadedDep, nil
}

func (injector *defaultInjector) findByType(t reflect.Type, optional bool) (*dependency, error) {
	var foundVal *dependency
	for _, v := range injector.dependencies {
		if v.reflectType.AssignableTo(t) {
			if foundVal != nil {
				return nil, fmt.Errorf("injector: there is a conflict when finding the dependency for %s", t.String())
			}

			foundVal = v
		}
	}

	if foundVal == nil && !optional {
		return nil, fmt.Errorf("injector: couldn't find the dependency for %s", t.String())
	}

	return foundVal, nil
}

func parseTag(tag string) (string, bool, error) {
	parts := strings.Split(tag, ",")
	switch len(parts) {
	case 2:
		if parts[1] != optionalTag {
			return "", false, fmt.Errorf("injector: %s is unexpected", parts[1])
		}

		return parts[0], parts[1] == optionalTag, nil
	case 1:
		return parts[0], false, nil
	case 0:
		return "", false, fmt.Errorf("injector: tag must not be empty")
	default:
		return "", false, fmt.Errorf("injector: unable to parse tag %s", tag)
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
