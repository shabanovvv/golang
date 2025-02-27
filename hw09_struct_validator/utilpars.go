package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type UtilPars struct {
	rr RuleResult
}
type RuleResult struct {
	RuleLen    int
	RuleMin    int
	RuleMax    int
	RuleRegexp string
	RuleIn     []string
}

func (up *UtilPars) ParseTag(typeField reflect.Type, tagField string) error {
	parts := strings.Split(tagField, "|")
	for _, part := range parts {
		err := up.ResolveTag(typeField, part)
		if err != nil {
			return err
		}
	}

	return nil
}

func (up *UtilPars) ResolveTag(typeField reflect.Type, tagField string) error {
	switch {
	case strings.HasPrefix(tagField, "len"):
		ruleLen, err := up.ParsTagLen(typeField, tagField)
		if err != nil {
			return err
		}
		up.rr.RuleLen = ruleLen
	case strings.HasPrefix(tagField, "min"):
		ruleMin, err := up.ParsTagMin(tagField)
		if err != nil {
			return err
		}
		up.rr.RuleMin = ruleMin
	case strings.HasPrefix(tagField, "max"):
		ruleMax, err := up.ParsTagMax(tagField)
		if err != nil {
			return err
		}
		up.rr.RuleMax = ruleMax
	case strings.HasPrefix(tagField, "regexp"):
		ruleRegexp, err := up.ParsTagRegexp(tagField)
		if err != nil {
			return err
		}
		up.rr.RuleRegexp = ruleRegexp
	case strings.HasPrefix(tagField, "in"):
		ruleIn, err := up.ParsTagIn(tagField)
		if err != nil {
			return err
		}
		up.rr.RuleIn = ruleIn
	default:
		return nil
	}
	return nil
}

func (up *UtilPars) ParsTagLen(typeField reflect.Type, tagField string) (int, error) {
	values := strings.Split(tagField, ":")

	if len(values) != 2 {
		return 0, fmt.Errorf("invalid tag format; expected 'len:<length>', got '%s'", tagField)
	}

	expectLen, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return 0, err
	}
	if (typeField.Kind() == reflect.String ||
		typeField.Kind() == reflect.Slice && typeField.Elem().Kind() == reflect.String) &&
		expectLen > 0 {
		return expectLen, nil
	}

	return 0, fmt.Errorf("expected string type, got %s", typeField.Kind())
}

func (up *UtilPars) ParsTagMin(tagField string) (int, error) {
	v := strings.Split(tagField, ":")
	if len(v) != 2 {
		return 0, fmt.Errorf("неверный формат: '%s', ожидается 'max:value'", v)
	}
	if v[0] == "min" {
		valMin, err := strconv.Atoi(strings.TrimSpace(v[1]))
		if err != nil {
			return 0, err
		}
		return valMin, nil
	}
	return 0, errors.New("error parse tag min")
}

func (up *UtilPars) ParsTagMax(tagField string) (int, error) {
	v := strings.Split(tagField, ":")
	if len(v) != 2 {
		return 0, fmt.Errorf("неверный формат: '%s', ожидается 'max:value'", v)
	}
	if v[0] == "max" {
		valMax, err := strconv.Atoi(strings.TrimSpace(v[1]))
		if err != nil {
			return 0, err
		}
		return valMax, nil
	}
	return 0, errors.New("error parse tag max")
}

func (up *UtilPars) ParsTagRegexp(tagField string) (string, error) {
	values := strings.Split(tagField, ":")
	if len(values) != 2 {
		return "", fmt.Errorf("invalid tag format; expected 'len:<length>', got '%s'", tagField)
	}
	if values[0] == "regexp" {
		if _, err := regexp.Compile(values[1]); err != nil {
			return "", fmt.Errorf("invalid regular expression: %w", err)
		}
		return values[1], nil
	}
	return "", fmt.Errorf("expected 'regexp' as the first value, got '%s'", values[0])
}

func (up *UtilPars) ParsTagIn(tagField string) ([]string, error) {
	values := strings.Split(tagField, ":")
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid tag format; expected 'in:<value1,value2,...>', got '%s'", tagField)
	}
	if values[0] == "in" {
		validRoles := strings.Split(values[1], ",")
		if len(validRoles) != 2 {
			return nil, fmt.Errorf("invalid tag format; expected 'in:<value1,value2,...>', got '%s'", tagField)
		}
		return validRoles, nil
	}

	return nil, fmt.Errorf("expected 'in' as the first value, got '%s'", values[0])
}
