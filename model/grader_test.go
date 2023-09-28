package model

import (
    "fmt"
    "testing"
    "time"

    "github.com/eriq-augustine/autograder/util"
)

// Test some dates to make sure that are marshed from JSON correctly.
func TestDates(test *testing.T) {
    testCases := []dateTestCase{
        dateTestCase{"2023-09-28T04:00:20+00:00", getTime("2023-09-28T04:00:20+00:00")},
        dateTestCase{"2023-09-28T04:00:20Z", getTime("2023-09-28T04:00:20Z")},
    };

    for i, testCase := range testCases {
        jsonString := fmt.Sprintf(`{"time": "%s"}`, testCase.Input);

        actual := make(map[string]time.Time);
        err := util.JSONFromString(jsonString, &actual);
        if (err != nil) {
            test.Fatal(err);
        }

        if (actual["time"] != testCase.Expected) {
            test.Errorf("Date case %d does not match. Expected '%s', Got '%s'.", i, testCase.Expected, actual["time"]);
        }
    }
}

func TestSummaryMarshalling(test *testing.T) {
    for i, testCase := range testSummaries {
        summary := SubmissionSummary{}
        util.JSONFromString(testCase.JSON, &summary);

        // TODO(eriq): Expand to also use a direct equality check.
        if (summary.String() != testCase.Summary.String()) {
            test.Errorf("Summaries %d do not match:\n--- Exepcted ---\n%s\n--- Actual ---\n%s\n---\n", i, testCase.Summary, summary);
        }
    }
}

type dateTestCase struct {
    Input string
    Expected time.Time
}

func getTime(value string) time.Time {
    result, err := time.Parse(time.RFC3339, value);
    if (err != nil) {
        panic(err);
    }

    return result;
}

func getSimpleTime(value string) time.Time {
    result, err := time.Parse(time.DateTime, value);
    if (err != nil) {
        panic(err);
    }

    return result;
}

type summaryTestCase struct {
    Summary SubmissionSummary
    JSON string
}

var testSummaries []summaryTestCase = []summaryTestCase{
    summaryTestCase{
        Summary: SubmissionSummary{
            ID: "1",
            Message: "Test 1!",
            MaxPoints: 1,
            Score: 1,
            GradingStartTime: getSimpleTime("2023-09-28 04:00:20"),
        },
        JSON: `
            {
                "id": "1",
                "message": "Test 1!",
                "max_points": 1,
                "score": 1,
                "grading_start_time": "2023-09-28T04:00:20+00:00"
            }
        `,
    },
}
