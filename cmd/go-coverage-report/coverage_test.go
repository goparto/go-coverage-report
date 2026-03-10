package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cov, err := ParseCoverage("testdata/01-new-coverage.txt")
	require.NoError(t, err)

	assert.EqualValues(t, 102, cov.TotalStmt)
	assert.EqualValues(t, 92, cov.CoveredStmt)
	assert.EqualValues(t, 10, cov.MissedStmt)
	assert.InDelta(t, 90.196, cov.Percent(), 0.001)
}

func TestCoverage_ByPackage(t *testing.T) {
	cov, err := ParseCoverage("testdata/01-new-coverage.txt")
	require.NoError(t, err)

	pkgs := cov.ByPackage()
	assert.Len(t, pkgs, 1)

	pkgCov := pkgs["github.com/fgrosse/prioqueue"]
	assert.NotNil(t, pkgCov)
	assert.EqualValues(t, 102, pkgCov.TotalStmt)
	assert.EqualValues(t, 92, pkgCov.CoveredStmt)
	assert.EqualValues(t, 10, pkgCov.MissedStmt)
}

func TestIsMockFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"foo_mock.go", true},
		{"github.com/foo/bar/baz_mock.go", true},
		{"mock.go", false},
		{"foo_mock_test.go", false},
		{"foo.go", false},
		{"foo_mocked.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.expected, isMockFile(tt.path))
		})
	}
}

func TestCoverage_ExcludeMockFiles(t *testing.T) {
	cov := New([]*Profile{
		{FileName: "github.com/foo/bar/real.go", TotalStmt: 10, CoveredStmt: 8, MissedStmt: 2},
		{FileName: "github.com/foo/bar/real_mock.go", TotalStmt: 5, CoveredStmt: 5, MissedStmt: 0},
		{FileName: "github.com/foo/bar/other_mock.go", TotalStmt: 3, CoveredStmt: 1, MissedStmt: 2},
	})

	cov.ExcludeMockFiles()

	assert.Len(t, cov.Files, 1)
	assert.Contains(t, cov.Files, "github.com/foo/bar/real.go")
	assert.EqualValues(t, 10, cov.TotalStmt)
	assert.EqualValues(t, 8, cov.CoveredStmt)
	assert.EqualValues(t, 2, cov.MissedStmt)
}
