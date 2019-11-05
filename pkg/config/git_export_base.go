package config

import "path/filepath"

type GitExportBase struct {
	*GitExport
	StageDependencies *StageDependencies

	raw *rawGit
}

func (c *ExportBase) GitMappingAdd() string {
	if c.Add == "/" {
		return ""
	}
	return filepath.ToSlash(c.Add)
}

func (c *ExportBase) GitMappingTo() string {
	return filepath.ToSlash(c.To)
}

func (c *ExportBase) GitMappingIncludePaths() []string {
	return gitMappingPaths(c.IncludePaths)
}

func (c *ExportBase) GitMappingExcludePath() []string {
	return gitMappingPaths(c.ExcludePaths)
}

func (c *GitExportBase) GitMappingStageDependencies() *StageDependencies {
	s := &StageDependencies{}
	s.Install = gitMappingPaths(c.StageDependencies.Install)
	s.BeforeSetup = gitMappingPaths(c.StageDependencies.BeforeSetup)
	s.Setup = gitMappingPaths(c.StageDependencies.Setup)
	return s
}

func gitMappingPaths(paths []string) []string {
	var newPaths []string
	for _, path := range paths {
		newPaths = append(newPaths, filepath.ToSlash(path))
	}

	return newPaths
}

func (c *GitExportBase) validate() error {
	return nil
}
