package formatnumber

import "testing"

func TestFormatPositiveDecimal(t *testing.T) {
	mask := "#,##0.00"
	value := 12345.67
	expected := "12,345.67"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatNegativeDecimal(t *testing.T) {
	mask := "#,##0.00"
	value := -12345.67
	expected := "-12,345.67"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatPositiveDecimalWithPrecision(t *testing.T) {
	mask := "#,##0.000"
	value := 12345.67
	expected := "12,345.670"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatNegativeDecimalWithPrecision(t *testing.T) {
	mask := "#,##0.000"
	value := -12345.67
	expected := "-12,345.670"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatWithCurrencySymbol(t *testing.T) {
	mask := "€ #,##0.00"
	value := 12345.67
	expected := "€ 12,345.67"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatWithCustomPrefixAndSuffix(t *testing.T) {
	mask := "Custom Prefix #,##0.00 Custom Suffix"
	value := 12345.67
	expected := "Custom Prefix 12,345.67 Custom Suffix"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatWithZeroValue(t *testing.T) {
	mask := "#,##0.00"
	value := 0.0
	expected := "0.00"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestFormatMaxFloat64(t *testing.T) {
	mask := "#,##0.00"
	value := 1797693134862315.0 // maximum float64 value
	expected := "1,797,693,134,862,315.00"

	result := Format(mask, value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
