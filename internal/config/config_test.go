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

	ReadEnv()

	res1 := GetString("var1")
	test.AssertEqual(t, res1, "someString")

	res2, err := GetInt("var2")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 5)

	res3, err := GetFloat("var3")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 4.5)

	res4, err := GetBool("var4", false)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, true)

}

func TestConfig_FromPropertyFile(t *testing.T) {
	ReadPorpertyFile(testDataFolder + "test.properties")

	res1 := GetString("var1")
	test.AssertEqual(t, res1, "someOtherString")

	res2, err := GetInt("var2")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 4)

	res3, err := GetFloat("var3")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 3.5)

	res4, err := GetBool("var4", true)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, false)

}

func TestConfig_PropertyOverwrite(t *testing.T) {
	os.Setenv("var1", "someString")
	os.Setenv("var2", "5")
	os.Setenv("var3", "4.5")
	os.Setenv("var4", "true")

	ReadPorpertyFile(testDataFolder + "test.properties")
	ReadEnv()

	res1 := GetString("var1")
	test.AssertEqual(t, res1, "someString")

	res2, err := GetInt("var2")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res2, 5)

	res3, err := GetFloat("var3")
	test.AssertNotError(t, err)
	test.AssertEqual(t, res3, 4.5)

	res4, err := GetBool("var4", false)
	test.AssertNotError(t, err)
	test.AssertEqual(t, res4, true)

}
