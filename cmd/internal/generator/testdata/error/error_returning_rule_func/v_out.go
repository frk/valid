// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

func (v Validator) Validate() error {
	if ok, err := mypkg.RuleWithErr1(v.F1); err != nil {
		return err
	} else if !ok {
		return errors.New("F1 is not valid")
	}
	if v.F2 != nil {
		if ok, err := mypkg.RuleWithErr1(*v.F2); err != nil {
			return err
		} else if !ok {
			return errors.New("F2 is not valid")
		}
	}
	if v.F3 == "" {
		return errors.New("F3 is required")
	} else if ok, err := mypkg.RuleWithErr1(v.F3); err != nil {
		return err
	} else if !ok {
		return errors.New("F3 is not valid")
	}
	if v.F4 == nil || *v.F4 == "" {
		return errors.New("F4 is required")
	} else if ok, err := mypkg.RuleWithErr1(*v.F4); err != nil {
		return err
	} else if !ok {
		return errors.New("F4 is not valid")
	}
	if v.F5 == "" {
		return errors.New("F5 is required")
	} else if ok, err := mypkg.RuleWithErr1(v.F5); err != nil {
		return err
	} else if !ok {
		return errors.New("F5 is not valid")
	} else if !valid.Email(v.F5) {
		return errors.New("F5 must be a valid email address")
	} else if ok, err := mypkg.RuleWithErr2(v.F5, 3, 13); err != nil {
		return err
	} else if !ok {
		return errors.New("F5 is not valid")
	}
	if v.F6 == nil || *v.F6 == "" {
		return errors.New("F6 is required")
	} else if !valid.Email(*v.F6) {
		return errors.New("F6 must be a valid email address")
	} else if ok, err := mypkg.RuleWithErr1(*v.F6); err != nil {
		return err
	} else if !ok {
		return errors.New("F6 is not valid")
	} else if ok, err := mypkg.RuleWithErr2(*v.F6, 3, 8); err != nil {
		return err
	} else if !ok {
		return errors.New("F6 is not valid")
	}
	for _, e := range v.F7 {
		if ok, err := mypkg.RuleWithErr2(e, 3, 8); err != nil {
			return err
		} else if !ok {
			return errors.New("F7 is not valid")
		}
	}
	for _, e := range v.F8 {
		if e == nil || *e == "" {
			return errors.New("F8 is required")
		} else if ok, err := mypkg.RuleWithErr2(*e, 3, 8); err != nil {
			return err
		} else if !ok {
			return errors.New("F8 is not valid")
		} else if ok, err := mypkg.RuleWithErr1(*e); err != nil {
			return err
		} else if !ok {
			return errors.New("F8 is not valid")
		}
	}
	return nil
}