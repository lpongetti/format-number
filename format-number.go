package formatnumber

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type MaskObj struct {
	Prefix              string
	Suffix              string
	Mask                string
	MaskHasNegativeSign bool
	MaskHasPositiveSign bool
	Decimal             string
	Separator           string
	Integer             string
	Fraction            string
}

type ValObj struct {
	Value    float64
	Sign     string
	Integer  string
	Fraction string
	Result   string
}

func getIndex(mask string) int {
	for i, char := range mask {
		val := string(char)
		if unicode.IsDigit(char) || val == "-" || val == "+" || val == "#" {
			return i
		}
		i++
	}
	return -1
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func processMask(mask string) MaskObj {
	maskObj := MaskObj{}
	lenMask := len(mask)
	start := getIndex(mask)
	maskObj.Prefix = ""
	if start > 0 {
		maskObj.Prefix = mask[:start]
	}

	end := getIndex(reverse(mask))
	offset := lenMask - end
	if offset < lenMask {
		substr := mask[offset : offset+1]
		if substr == "." || substr == "," {
			offset++
		}
	}
	maskObj.Suffix = ""
	if end > 0 {
		maskObj.Suffix = mask[offset:]
	}

	maskObj.Mask = mask[start:offset]
	maskObj.MaskHasNegativeSign = strings.HasPrefix(maskObj.Mask, "-")
	maskObj.MaskHasPositiveSign = strings.HasPrefix(maskObj.Mask, "+")

	result := make([]string, 0)
	for _, char := range maskObj.Mask {
		val := string(char)
		if !unicode.IsDigit(char) && val != "-" && val != "+" && val != "#" {
			result = append(result, string(char))
		}
	}

	maskObj.Decimal = "."
	if len(result) > 0 {
		maskObj.Decimal = result[len(result)-1]
	}
	maskObj.Separator = ","
	if len(result) > 1 {
		maskObj.Separator = result[0]
	}

	splitResult := strings.Split(maskObj.Mask, maskObj.Decimal)
	maskObj.Integer = splitResult[0]
	if len(splitResult) > 1 {
		maskObj.Fraction = splitResult[1]
	}

	return maskObj
}

func processValue(value float64, maskObj MaskObj, enforceMaskSign bool) ValObj {
	isNegative := value < 0
	valObj := ValObj{
		Value: value,
	}

	if isNegative {
		valObj.Value = -valObj.Value
		valObj.Sign = "-"
	} else {
		valObj.Sign = ""
	}

	valStr := strconv.FormatFloat(valObj.Value, 'f', len(maskObj.Fraction), 64)
	valParts := strings.Split(valStr, ".")

	valObj.Integer = valParts[0]
	if len(valParts) > 1 {
		valObj.Fraction = valParts[1]
	} else {
		valObj.Fraction = ""
	}

	posTrailZero := strings.LastIndex(maskObj.Fraction, "0")
	if posTrailZero >= 0 && len(valObj.Fraction) <= posTrailZero {
		valObj.Fraction = fmt.Sprintf("%0*s", posTrailZero+1, valObj.Fraction)
	}

	addSeparators(&valObj, maskObj)

	if valObj.Result == "0" || valObj.Result == "" {
		isNegative = false
		valObj.Sign = ""
	}

	if !isNegative && maskObj.MaskHasPositiveSign {
		valObj.Sign = "+"
	} else if isNegative && maskObj.MaskHasPositiveSign {
		valObj.Sign = "-"
	} else if isNegative {
		if enforceMaskSign && !maskObj.MaskHasNegativeSign {
			valObj.Sign = ""
		} else {
			valObj.Sign = "-"
		}
	}

	return valObj
}

func addSeparators(valObj *ValObj, maskObj MaskObj) {
	valObj.Result = ""
	szSep := strings.Split(maskObj.Integer, maskObj.Separator)
	maskInteger := strings.Join(szSep, "")

	posLeadZero := strings.Index(maskInteger, "0")
	if posLeadZero > -1 {
		for len(valObj.Integer) < len(maskInteger)-posLeadZero {
			valObj.Integer = "0" + valObj.Integer
		}
	} else if valObj.Integer == "0" {
		valObj.Integer = ""
	}

	posSeparator := len(szSep[len(szSep)-1])
	if posSeparator > 0 {
		lenValInteger := len(valObj.Integer)
		offset := lenValInteger % posSeparator
		for i := 0; i < lenValInteger; i++ {
			valObj.Result += string(valObj.Integer[i])
			if (i-offset+1)%posSeparator == 0 && i < lenValInteger-posSeparator {
				valObj.Result += maskObj.Separator
			}
		}
	} else {
		valObj.Result = valObj.Integer
	}

	if maskObj.Fraction != "" && valObj.Fraction != "" {
		valObj.Result += maskObj.Decimal + valObj.Fraction
	}
}

func FormatWithOptions(mask string, value float64, enforceMaskSign bool) string {
	if mask == "" || strconv.FormatFloat(value, 'f', -1, 64) == "NaN" {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}

	maskObj := processMask(mask)
	valObj := processValue(value, maskObj, enforceMaskSign)
	return maskObj.Prefix + valObj.Sign + valObj.Result + maskObj.Suffix
}

func Format(mask string, value float64) string {
	if mask == "" {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}

	maskObj := processMask(mask)
	valObj := processValue(value, maskObj, false)
	return maskObj.Prefix + valObj.Sign + valObj.Result + maskObj.Suffix
}
