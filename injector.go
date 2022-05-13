package sen

import (
	"fmt"
	"reflect"
)

const (
	autoInjectionTag = "*"
	injectTag        = "inject"
)

func newInjector() *defaultInjector {
	return &defaultInjector{
		dependencies: make(map[string]*dependency),
	}
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

		if !loadedDep.reflectType.AssignableTo(fieldType) {
			return fmt.Errorf("injector: %s is not assignable from %s", fieldType, loadedDep.reflectType)
		}

		fieldValue.Set(loadedDep.reflectValue)
	}

	return nil
}

func (injector *defaultInjector) loadDepForTag(tag string, t reflect.Type) (*dependency, error) {
	if tag == autoInjectionTag {
		return injector.findByType(t)
	}

	loadedDep, found := injector.dependencies[tag]
	if !found {
		return nil, fmt.Errorf("injector: %s is not registered", tag)
	}

	return loadedDep, nil
}

func (injector *defaultInjector) findByType(t reflect.Type) (*dependency, error) {
	var foundVal *dependency
	for _, v := range injector.dependencies {
		if v.reflectType.AssignableTo(t) {
			if foundVal != nil {
				return nil, fmt.Errorf("injector: there is a conflict when finding the dependency for %s", t.String())
			}

			foundVal = v
		}
	}

	if foundVal == nil {
		return nil, fmt.Errorf("injector: couldn't find the dependency for %s", t.String())
	}

	return foundVal, nil
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
