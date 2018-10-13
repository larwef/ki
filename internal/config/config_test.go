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

	FromEnv()

	test.AssertEqual(t, GetString("var1"), "someString")
	test.AssertEqual(t, GetInt("var2"), 5)
	test.AssertEqual(t, GetFloat("var3"), 4.5)
	test.AssertEqual(t, GetBool("var4", false), true)
}

func TestConfig_FromPropertyFile(t *testing.T) {
	FromPorpertyFile(testDataFolder + "test.properties")

	test.AssertEqual(t, GetString("var1"), "someOtherString")
	test.AssertEqual(t, GetInt("var2"), 4)
	test.AssertEqual(t, GetFloat("var3"), 3.5)
	test.AssertEqual(t, GetBool("var4", true), false)
}

func TestConfig_PropertyOverwrite(t *testing.T) {
	os.Setenv("var1", "someString")
	os.Setenv("var2", "5")
	os.Setenv("var3", "4.5")
	os.Setenv("var4", "true")

	FromPorpertyFile(testDataFolder + "test.properties")
	FromEnv()

	test.AssertEqual(t, GetString("var1"), "someString")
	test.AssertEqual(t, GetInt("var2"), 5)
	test.AssertEqual(t, GetFloat("var3"), 4.5)
	test.AssertEqual(t, GetBool("var4", false), true)
}
