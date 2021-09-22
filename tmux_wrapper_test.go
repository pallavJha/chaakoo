package chaakoo

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

type TmuxWrapperTestCase struct {
	Id          int
	Dimension   *Dimension
	SessionName string
	Windows     []*Window
	Commands    []struct {
		Name string
		Args string
	}
}

func (c TmuxWrapperTestSuite) testTmuxWrapperApply(t *testing.T) {
	var testCases []TmuxWrapperTestCase
	if err := viper.UnmarshalKey("configs", &testCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for _, testCase := range testCases {
		t.Log("testing, id", testCase.Id)
		config := &Config{
			SessionName: testCase.SessionName,
			Windows:     testCase.Windows,
		}
		err := config.Validate()
		require.NoError(t, err)
		err = config.Parse()
		require.NoError(t, err)
		wrapper := NewTmuxWrapper(config, testCase.Dimension)
		err = wrapper.Apply()
		require.NoError(t, err)
	}
}
