package sen

// NewInjector exports newInjector for testing purpose only.
var NewInjector = newInjector

// GetComponent returns a registered component via name.
// The function is created for testing purpose only.
func GetComponent(app *Application, name string) (interface{}, error) {
	dep, err := app.defaultInjector.loadDepForTag(name, nil)
	if err != nil {
		return nil, err
	}

	return dep.value, nil
}
