package tracing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/tracing"
)

func TestCompare(t *testing.T) {
	workfactor := 10
	hasher := &tracing.TracedBCrypt{
		WorkFactor: workfactor,
	}

	expectedPassword := "hello world"
	expectedPasswordHash, err := hasher.Hash(context.TODO(), []byte(expectedPassword))
	assert.NoError(t, err)
	assert.NotNil(t, expectedPasswordHash)

	expectedTagsSuccess := map[string]interface{}{
		tracing.WorkFactorTagName: int(workfactor),
	}

	expectedTagsError := map[string]interface{}{
		tracing.WorkFactorTagName: int(workfactor),
		"error":                   true,
	}

	testCases := []struct {
		testDescription  string
		providedPassword string
		expectedTags     map[string]interface{}
		shouldError      bool
	}{
		{
			testDescription:  "should not return an error if hash of provided password matches hash of expected password",
			providedPassword: expectedPassword,
			expectedTags:     expectedTagsSuccess,
			shouldError:      false,
		},
		{
			testDescription:  "should return an error if hash of provided password does not match hash of expected password",
			providedPassword: "some invalid password",
			expectedTags:     expectedTagsError,
			shouldError:      true,
		},
	}

	for _, test := range testCases {
		t.Run(test.testDescription, func(t *testing.T) {
			hash, err := hasher.Hash(context.TODO(), []byte(test.providedPassword))
			assert.NoError(t, err)
			assert.NotNil(t, hash)

			mockedTracer.Reset()

			err = hasher.Compare(context.TODO(), expectedPasswordHash, []byte(test.providedPassword))
			if test.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			spans := mockedTracer.FinishedSpans()
			assert.Len(t, spans, 1)
			span := spans[0]

			assert.Equal(t, tracing.CompareOpName, span.OperationName)
			assert.Equal(t, test.expectedTags, span.Tags())
		})
	}
}

func TestHashCreatesSpanWithCorrectTags(t *testing.T) {
	validWorkFactor := 10
	invalidWorkFactor := 1000 // this is an invalid work factor that will cause the call to Hash to fail!
	password := []byte("bar")

	expectedTagsSuccess := map[string]interface{}{
		tracing.WorkFactorTagName: int(validWorkFactor),
	}

	expectedTagsError := map[string]interface{}{
		tracing.WorkFactorTagName: int(invalidWorkFactor),
		"error":                   true,
	}

	testCases := []struct {
		testDescription string
		expectedTags    map[string]interface{}
		workFactor      int
		shouldError     bool
	}{
		{
			testDescription: "tests expected tags are created when call to Hash succeeds",
			expectedTags:    expectedTagsSuccess,
			workFactor:      validWorkFactor,
			shouldError:     false,
		},
		{
			testDescription: "tests expected tags are created when call to Hash fails",
			expectedTags:    expectedTagsError,
			workFactor:      invalidWorkFactor,
			shouldError:     true,
		},
	}

	for _, test := range testCases {
		t.Run(test.testDescription, func(t *testing.T) {
			mockedTracer.Reset()
			hasher := &tracing.TracedBCrypt{
				WorkFactor: test.workFactor,
			}

			_, err := hasher.Hash(context.TODO(), password)

			if test.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			spans := mockedTracer.FinishedSpans()
			assert.Len(t, spans, 1)
			span := spans[0]

			assert.Equal(t, tracing.HashOpName, span.OperationName)
			assert.Equal(t, test.expectedTags, span.Tags())
		})
	}
}
