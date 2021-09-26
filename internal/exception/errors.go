package exception

var (
	_ error = BadConfigs{}
	_ error = InvalidInput{}
)

// InvalidInput is a generic struct returned for otherwise invalid input.
type InvalidInput struct {
	InputName        string
	AcceptableInputs []string
}

func (ii InvalidInput) Error() string {
	return "Invalid input %s was given, expected something in %s"
}

// UserError for invalidinput should detail what they messed up but potentially not leak internal options.
func (ii InvalidInput) UserError() string {
	return "Invalid request of '%s' please check your request for validity."
}

// BadConfigs is an error that should only be thrown during standup or reloads.
type BadConfigs struct {
	InputName string
}

func (bc BadConfigs) Error() string {
	return "Invalid config input %s was given check that things are where they should be prior to restart."
}

// UserError for badconfigs is empty as it should not be exposed to users.
func (bc BadConfigs) UserError() string {
	return ""
}
