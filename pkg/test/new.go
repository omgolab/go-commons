package gctest

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// we may create a test case like as follows:
// NewTest ->
//
//	NewScenario ->
//		NewTestCaseGiven -> When -> Then -> Error -> FinalStep
var (
	scenarioPrefix  = []byte("Scenario: ")
	givenPrefix     = []byte("\tGiven: ")
	whenPrefix      = []byte(" When: ")
	thenPrefix      = []byte(" Then: ")
	returnsErrors   = []byte(" And returns errors.")
	returnsNoErrors = []byte(" And returns NO errors.")
)

type BeforeFn[A any] func(t *testing.T) (A, error)
type ExecFn[A, B any] func(t *testing.T, arg A) (B, error)
type AfterFn[B any] func(t *testing.T, arg B) error

type testCfg struct {
	t                       *testing.T
	beforeAllTestsFn        BeforeFn[any]
	afterAllTestsFn         AfterFn[any]
	commonBeforeEachTestsFn BeforeFn[any]
	commonAfterEachTestsFn  AfterFn[any]
	isParallel              bool
	activeTestCaseID        atomic.Int64
}

type testCase struct {
	id   int64
	cfg  *testCfg
	name []byte
}

func NewTest(t *testing.T, opts ...TestOption) TestCase {
	tc := &testCfg{
		t:                t,
		activeTestCaseID: atomic.Int64{},
	}

	// apply options
	for _, f := range opts {
		f(tc)
	}

	// enable parallel exec if available
	if tc.isParallel {
		tc.t.Parallel()
	}

	// execute beforeAllTestsFn if not nil
	if tc.beforeAllTestsFn != nil {
		_, err := tc.beforeAllTestsFn(t)
		assert.NoError(t, err)
	}

	// TODO: exec afterAllTestF

	return testCase{
		cfg: tc,
	}
}