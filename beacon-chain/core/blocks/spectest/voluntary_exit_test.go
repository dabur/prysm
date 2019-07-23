package spectest

import (
	"io/ioutil"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"gopkg.in/d4l3k/messagediff.v1"
)

func runVoluntaryExitTest(t *testing.T, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not load file %v", err)
	}

	test := &BlockOperationTest{}
	if err := testutil.UnmarshalYaml(file, test); err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}

	if err := spectest.SetConfig(test.Config); err != nil {
		t.Fatal(err)
	}

	if len(test.TestCases) == 0 {
		t.Fatal("No tests!")
	}

	for _, tt := range test.TestCases {
		t.Run(tt.Description, func(t *testing.T) {
			helpers.ClearAllCaches()

			body := &ethpb.BeaconBlockBody{VoluntaryExits: []*ethpb.VoluntaryExit{tt.VoluntaryExit}}

			postState, err := blocks.ProcessVoluntaryExits(tt.Pre, body, true)
			// Note: This doesn't test anything worthwhile. It essentially tests
			// that *any* error has occurred, not any specific error.
			if tt.Post == nil {
				if err == nil {
					t.Fatal("Did not fail when expected")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !proto.Equal(postState, tt.Post) {
				diff, _ := messagediff.PrettyDiff(postState, tt.Post)
				t.Log(diff)
				t.Fatal("Post state does not match expected")
			}
		})
	}
}

var exitPrefix = "tests/operations/voluntary_exit/"

func TestVoluntaryExitMinimal(t *testing.T) {
	filepath, err := bazel.Runfile(exitPrefix + "voluntary_exit_mainnet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	runVoluntaryExitTest(t, filepath)
}

func TestVoluntaryExitMainnet(t *testing.T) {
	filepath, err := bazel.Runfile(exitPrefix + "voluntary_exit_mainnet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	runVoluntaryExitTest(t, filepath)
}