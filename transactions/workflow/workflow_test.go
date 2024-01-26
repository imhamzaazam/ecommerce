package workflow

import (
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"testing"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	//Cannot figure out how to import CreateTransactionActivty
	//env.OnActivity(CreateTransaction, mock.Anything, "World").Return("Hello World!", nil)

	//env.ExecuteWorkflow(CreateTransactionWorkflow, "World")
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	//var greeting string
	//require.NoError(t, env.GetWorkflowResult(&greeting))
	//require.Equal(t, "Hello World!", greeting)
}
