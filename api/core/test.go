package core

import (
    "fmt"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/eriq-augustine/autograder/config"
    "github.com/eriq-augustine/autograder/grader"
    "github.com/eriq-augustine/autograder/util"
)

var server *httptest.Server;
var serverURL string;

func startTestServer(routes *[]*Route) {
    if (server != nil) {
        panic("Test server already started.");
    }

    server = httptest.NewServer(GetRouteServer(routes));
    serverURL = server.URL;
}

func stopTestServer() {
    if (server != nil) {
        server.Close();

        server = nil;
        serverURL = "";
    }
}

// Common setup for all API tests.
func APITestingMain(suite *testing.M, routes *[]*Route) {
    config.EnableTestingMode(false, true);
    config.NO_AUTH.Set(false);

    err := grader.LoadCourses();
    if (err != nil) {
        fmt.Printf("Failed to load test courses: '%v'.", err);
        os.Exit(1);
    }

    startTestServer(routes);
    defer stopTestServer();

    os.Exit(suite.Run())
}

// Make a request to the test server using fields for
// a standard test request plus whatever other fields are specified.
// Provided fields will override base fields.
func SendTestAPIRequest(test *testing.T, endpoint string, fields map[string]any) *APIResponse {
    return SendTestAPIRequestFull(test, endpoint, fields, nil);
}

func SendTestAPIRequestFull(test *testing.T, endpoint string, fields map[string]any, paths []string) *APIResponse {
    url := serverURL + endpoint;

    content := map[string]any{
        "course-id": "COURSE101",
        "user-email": "admin@test.com",
        "user-pass": util.Sha256HexFromString("admin"),
    };

    for key, value := range fields {
        content[key] = value;
    }

    form := map[string]string{
        API_REQUEST_CONTENT_KEY: util.MustToJSON(content),
    };

    var responseText string;
    var err error;

    if (len(paths) == 0) {
        responseText, err = util.PostNoCheck(url, form);
    } else {
        responseText, err = util.PostFiles(url, form, paths, false);
    }

    if (err != nil) {
        test.Fatalf("API POST returned an error: '%v'.", err);
    }

    var response APIResponse;
    err = util.JSONFromString(responseText, &response);
    if (err != nil) {
        test.Fatalf("Could not unmarshal JSON response '%s': '%v'.", responseText, err);
    }

    return &response;
}
