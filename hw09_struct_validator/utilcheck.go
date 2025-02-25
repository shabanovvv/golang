package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type UtilCheck struct{}

func (uc UtilCheck) CheckTag(valueField reflect.Value, rr RuleResult) bool {
	b := true
	switch valueField.Kind() { //nolint:exhaustive
	case reflect.String:
		if rr.RuleLen > 0 {
			b = uc.CheckLenString(valueField.Interface().(string), rr.RuleLen)
		}
		if len(rr.RuleRegexp) > 0 {
			b = b && uc.CheckRegexp(valueField.Interface().(string), rr.RuleRegexp)
		}
		if rr.RuleIn != nil {
			b = b && uc.CheckInString(valueField.String(), rr.RuleIn)
		}
	case reflect.Int:
		if rr.RuleMin > 0 {
			b = uc.CheckMin(valueField.Interface().(int), rr.RuleMin)
		}
		if rr.RuleMax > 0 {
			b = b && uc.CheckMax(valueField.Interface().(int), rr.RuleMax)
		}
		if rr.RuleIn != nil {
			b = b && uc.CheckInInt(valueField.Interface().(int), rr.RuleIn)
		}
	case reflect.Slice:
		if rr.RuleLen > 0 {
			b = uc.CheckLenSlice(valueField.Interface().([]string), rr.RuleLen)
		}
	default:
		return true
	}

	return b
}

func (uc UtilCheck) CheckLenString(valueField string, ruleLen int) bool {
	return ruleLen > 0 && len(valueField) == ruleLen
}

func (uc UtilCheck) CheckLenSlice(valueField []string, ruleLen int) bool {
	sliceValue := reflect.ValueOf(valueField)
	for i := 0; i < sliceValue.Len(); i++ {
		if ruleLen > 0 && len(sliceValue.Index(i).Interface().(string)) != ruleLen {
			return false
		}
	}
	return true
}

func (uc UtilCheck) CheckMin(valueField, ruleMin int) bool {
	return valueField >= ruleMin
}

func (uc UtilCheck) CheckMax(valueField, ruleMax int) bool {
	return valueField <= ruleMax
}

func (uc UtilCheck) CheckRegexp(valueField, ruleRegexp string) bool {
	re, err := regexp.Compile(ruleRegexp)
	if err != nil {
		return false
	}
	return re.MatchString(valueField)
}

func (uc UtilCheck) CheckInString(valueField string, ruleIn []string) bool {
	valueStr := fmt.Sprintf("%v", valueField)
	for _, valid := range ruleIn {
		if !strings.Contains(valueStr, valid) {
			return false
		}
	}
	return true
}

func (uc UtilCheck) CheckInInt(valueField int, ruleIn []string) bool {
	ruleMin, err := strconv.Atoi(ruleIn[0])
	if err != nil {
		return false
	}
	ruleMax, err := strconv.Atoi(ruleIn[1])
	if err != nil {
		return false
	}
	return valueField >= ruleMin && valueField <= ruleMax
}
