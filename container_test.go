package container

import (
	"fmt"
	"testing"

	"github.com/frozzare/go-assert"
)

type User struct {
	Name string
}

type Friend struct {
	User *User `inject`
}

func TestInstance(t *testing.T) {
	assert.NotNil(t, Instance())
}

func TestSimpleContainerValue(t *testing.T) {
	c := Instance()

	assert.False(t, c.Contains("Hello"))
	c.Bind("Hello", "World")

	assert.True(t, c.Contains("Hello"))
	v, err := c.Make("Hello")

	assert.Nil(t, err)
	assert.Equal(t, "World", v.(string))

	c.Bind("Hello", func() string {
		return "World"
	})
	v, _ = c.Make("Hello")

	assert.Equal(t, "World", v.(string))

	c.Remove("Hello")
	assert.False(t, c.Contains("Hello"))
}

func TestContainerFuncInject(t *testing.T) {
	c := Instance()

	assert.False(t, c.Contains("FuncInject"))
	c.Bind("FuncInject", func(text string, u *User) string {
		return fmt.Sprintf("%s %s", text, u.Name)
	})
	c.Bind("*container.User", &User{"Fredrik"})

	v, err := c.Make("FuncInject", "Hello")

	assert.Nil(t, err)
	assert.Equal(t, "Hello Fredrik", v.(string))

	c.Remove("FuncInject")
	assert.False(t, c.Contains("FuncInject"))
}

func TestContainerStructInject(t *testing.T) {
	c := Instance()

	assert.False(t, c.Contains("StructInject"))

	u := &User{"Fredrik"}
	c.Bind("StructInject", &Friend{u})
	c.Bind("*container.User", u)

	v, err := c.Make("StructInject")

	assert.Nil(t, err)
	assert.Equal(t, "Fredrik", v.(*Friend).User.Name)

	c.Remove("StructInject")
	assert.False(t, c.Contains("StructInject"))
}

func TestNonExistingKey(t *testing.T) {
	c := Instance()

	v, err := c.Make("Missing")

	assert.NotNil(t, err)
	assert.Nil(t, v)
}

func TestSingleton(t *testing.T) {
	c := Instance()

	assert.False(t, c.Contains("Singleton"))
	err := c.Singleton("Singleton", "Hello")
	assert.Nil(t, err)

	assert.True(t, c.Contains("Singleton"))
	v, err := c.Make("Singleton")

	assert.Nil(t, err)
	assert.Equal(t, "Hello", v.(string))

	err = c.Singleton("Singleton", "Hello2")
	assert.NotNil(t, err)
	v, err = c.Make("Singleton")

	assert.Nil(t, err)
	assert.Equal(t, "Hello", v.(string))
}

func TestAll(t *testing.T) {
	c := Instance()
	a := c.All()

	assert.True(t, a["All"] == nil)
	c.Bind("All", "Hello")

	a = c.All()
	assert.False(t, a["All"] == nil)
}
