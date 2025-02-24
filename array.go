package vjson

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// ArrayField is the type for validating arrays in a JSON
type ArrayField struct {
	name     string
	required bool
	items    Field

	minLength           int
	minLengthValidation bool

	maxLength           int
	maxLengthValidation bool
}

// To Force Implementing Field interface by ArrayField
var _ Field = (*ArrayField)(nil)

// GetName returns name of the field
func (a *ArrayField) GetName() string {
	return a.name
}

// GetType returns the Fields type
func (a *ArrayField) GetType() string {
	return "array"
}

// GetRequired returns true if field is required
func (a *ArrayField) GetRequired() bool {
	return a.required
}

// Validate is used for validating a value. it returns an error if the value is invalid.
func (a *ArrayField) Validate(v interface{}) error {
	if v == nil {
		if !a.required {
			return nil
		}
		return errors.Errorf("Value for %s field is required", a.name)
	}

	values, ok := v.([]interface{})
	if !ok {
		return errors.Errorf("Value of %s should be array", a.name)
	}

	var result error
	if a.minLengthValidation {
		if len(values) < a.minLength {
			result = multierror.Append(result, errors.Errorf("length of %s array should be at least %d", a.name, a.minLength))
		}
	}

	if a.maxLengthValidation {
		if len(values) > a.maxLength {
			result = multierror.Append(result, errors.Errorf("length of %s array should be at most %d", a.name, a.maxLength))
		}
	}

	for _, value := range values {
		err := a.items.Validate(value)
		if err != nil {
			result = multierror.Append(result, errors.Wrapf(err, "%v item is invalid in %s array", value, a.name))
		}
	}
	return result
}

// Required is called to make a field required in a JSON
func (a *ArrayField) Required() *ArrayField {
	a.required = true
	return a
}

// MinLength is called to set minimum length for an array field in a JSON
func (a *ArrayField) MinLength(length int) *ArrayField {
	a.minLength = length
	a.minLengthValidation = true
	return a
}

// MaxLength is called to set maximum length for an array field in a JSON
func (a *ArrayField) MaxLength(length int) *ArrayField {
	a.maxLength = length
	a.maxLengthValidation = true
	return a
}

// Array is the constructor of an array field.
func Array(name string, itemField Field) *ArrayField {
	return &ArrayField{
		name:     name,
		required: false,
		items:    itemField,
	}
}
