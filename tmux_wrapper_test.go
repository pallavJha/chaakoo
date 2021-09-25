package chaakoo

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/pallavJha/chaakoo/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type TmuxWrapperTestCase struct {
	ID          int
	Error       string
	Ignore      bool
	Dimension   *Dimension
	SessionName string
	Windows     []*Window
	Commands    []*struct {
		Name     string
		Args     string
		Stdout   string
		Stderr   string
		Err      string
		ExitCode int
	}
}

func (c TmuxWrapperTestSuite) testTmuxWrapperApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var testCases []TmuxWrapperTestCase
	if err := viper.UnmarshalKey("configs", &testCases); err != nil {
		t.Log("unable to read from the config", err)
		t.Fail()
	}
	for _, testCase := range testCases {
		if testCase.Ignore {
			continue
		}
		t.Log("testing, id", testCase.ID)
		config := &Config{
			SessionName: testCase.SessionName,
			Windows:     testCase.Windows,
		}
		err := config.Validate()
		require.NoError(t, err)
		err = config.Parse()
		require.NoError(t, err)
		wrapper := NewTmuxWrapper(config, testCase.Dimension)

		mockCmdExecutor := mocks.NewMockICommandExecutor(ctrl)
		wrapper.executor = mockCmdExecutor
		for _, command := range testCase.Commands {
			command.Args = strings.TrimSpace(command.Args)
			arguments := strings.Split(command.Args, " ")
			var errorToReturn error
			if len(command.Err) > 0 {
				errorToReturn = errors.New(command.Err)
			}
			mockCmdExecutor.EXPECT().Execute(command.Name, arguments).Return(
				command.Stdout, command.Stderr, command.ExitCode, errorToReturn,
			)
		}

		err = wrapper.Apply()
		if len(testCase.Error) > 0 {
			require.Error(t, err)
			require.EqualError(t, err, testCase.Error)
		} else {
			require.NoError(t, err)
		}
	}
}
