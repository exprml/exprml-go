package exprml_test

import (
	"fmt"
	"github.com/exprml/exprml-go"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	type Testcase struct {
		YamlInput string
	}
	testcases := map[string]*Testcase{}
	err := filepath.Walk("testdata/parser/error", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".in.yaml") {
			key := strings.TrimSuffix(path, ".in.yaml")
			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("fail to read file: %+v", err)
			}
			if _, ok := testcases[key]; !ok {
				testcases[key] = &Testcase{}
			}
			testcases[key].YamlInput = string(b)
			return nil
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	testcaseKeys := []string{}
	for key := range testcases {
		testcaseKeys = append(testcaseKeys, key)
	}
	slices.Sort(testcaseKeys)
	for _, name := range testcaseKeys {
		testcase := testcases[name]
		t.Run(name, func(t *testing.T) {
			decodeResult := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: testcase.YamlInput})
			require.False(t, decodeResult.IsError)

			parseResult := exprml.NewParser().Parse(&pb.ParseInput{Value: decodeResult.Value})
			require.True(t, parseResult.IsError)
		})
	}
}
