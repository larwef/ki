package config

import (
	"github.com/larwef/ki/test"
	"os"
	"testing"
)

var testDataFolder = "../../test/testdata/"

func TestConfig_FromEnv(t *testing.T) {
	os.Setenv("var1", "someString")
	os.Setenv("var2", "5")
	os.Setenv("var3", "4.5")
	os.Setenv("var4", "true")

	conf := New(ReturnError)
	conf.ReadEnv()

	res1, err := conf.GetString("var1", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res1, "someString")

	res2, err := conf.GetInt("var2", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 5)

	res3, err := conf.GetFloat("var3", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 4.5)

	res4, err := conf.GetBool("var4", false, true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, true)

}

func TestConfig_FromPropertyFile(t *testing.T) {
	conf := New(ReturnError)
	conf.ReadPropertyFile(testDataFolder + "test.properties")

	res1, err := conf.GetString("var1", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res1, "someOtherString")

	res2, err := conf.GetInt("var2", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 4)

	res3, err := conf.GetFloat("var3", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 3.5)

	res4, err := conf.GetBool("var4", true, true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, false)

}

func TestConfig_PropertyOverwrite(t *testing.T) {
	os.Setenv("var1", "someString")
	os.Setenv("var2", "5")
	os.Setenv("var3", "4.5")
	os.Setenv("var4", "true")

	conf := New(ReturnError)
	conf.ReadPropertyFile(testDataFolder + "test.properties")
	conf.ReadEnv()

	res1, err := conf.GetString("var1", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res1, "someString")

	res2, err := conf.GetInt("var2", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 5)

	res3, err := conf.GetFloat("var3", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 4.5)

	res4, err := conf.GetBool("var4", false, true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, true)
}

func TestConfig_MissingProperty(t *testing.T) {
	conf := New(ReturnError)

	res1, err := conf.GetString("var1", true)
	test.AssertEqual(t, res1, "")
	test.AssertEqual(t, err, MissingPropertyError("var1"))
}

func TestConfig_GetString(t *testing.T) {
	var stringPropTest = []struct {
		prop           string
		required       bool
		expectedResult string
		expectedError  error
	}{
		{"var1", false, "someString", nil},
		{"var1", true, "someString", nil},
		{"notSet", false, "", nil},
		{"notSet", true, "", MissingPropertyError("notSet")},
	}

	os.Setenv("var1", "someString")

	conf := New(ReturnError)
	conf.ReadEnv()

	for _, tt := range stringPropTest {
		actual, err := conf.GetString(tt.prop, tt.required)
		test.AssertEqual(t, actual, tt.expectedResult)
		test.AssertEqual(t, err, tt.expectedError)
	}
}

func TestConfig_GetInt(t *testing.T) {
	var intPropTest = []struct {
		prop           string
		required       bool
		expectedResult int
		expectedError  error
	}{
		{"var1", false, 5, nil},
		{"var1", true, 5, nil},
		{"notSet", false, 0, nil},
		{"notSet", true, 0, MissingPropertyError("notSet")},
	}

	os.Setenv("var1", "5")

	conf := New(ReturnError)
	conf.ReadEnv()

	for _, tt := range intPropTest {
		actual, err := conf.GetInt(tt.prop, tt.required)
		test.AssertEqual(t, actual, tt.expectedResult)
		test.AssertEqual(t, err, tt.expectedError)
	}
}

func TestConfig_GetFloat(t *testing.T) {
	var floatPropTest = []struct {
		prop           string
		required       bool
		expectedResult float64
		expectedError  error
	}{
		{"var1", false, 5.0, nil},
		{"var1", true, 5.0, nil},
		{"var2", false, 5.5, nil},
		{"var2", true, 5.5, nil},
		{"notSet", false, 0.0, nil},
		{"notSet", true, 0.0, MissingPropertyError("notSet")},
	}

	os.Setenv("var1", "5")
	os.Setenv("var2", "5.5")

	conf := New(ReturnError)
	conf.ReadEnv()

	for _, tt := range floatPropTest {
		actual, err := conf.GetFloat(tt.prop, tt.required)
		test.AssertEqual(t, actual, tt.expectedResult)
		test.AssertEqual(t, err, tt.expectedError)
	}
}

func TestConfig_GetBool(t *testing.T) {
	var boolPropTest = []struct {
		prop           string
		defaul         bool
		required       bool
		expectedResult bool
		expectedError  error
	}{
		{"var1", false, false, false, nil},
		{"var1", true, false, false, nil},
		{"var1", false, true, false, nil},
		{"var1", true, true, false, nil},

		{"var2", false, false, true, nil},
		{"var2", true, false, true, nil},
		{"var2", false, true, true, nil},
		{"var2", true, true, true, nil},

		{"notSet", false, false, false, nil},
		{"notSet", true, false, true, nil},
		{"notSet", false, true, false, MissingPropertyError("notSet")},
		{"notSet", true, true, true, MissingPropertyError("notSet")},
	}

	os.Setenv("var1", "false")
	os.Setenv("var2", "true")

	conf := New(ReturnError)
	conf.ReadEnv()

	for _, tt := range boolPropTest {
		actual, err := conf.GetBool(tt.prop, tt.defaul, tt.required)
		test.AssertEqual(t, actual, tt.expectedResult)
		test.AssertEqual(t, err, tt.expectedError)
	}
}
