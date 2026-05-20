package main

import (
	"context"
	"fmt"
	"sort"

	"review-diff/changelist"
	"review-diff/diffparser"
	"review-diff/gitrunner"
	"review-diff/patchloader"
	"review-diff/session"
)

type App struct {
	ctx    context.Context
	sess   *session.Session
	runner *gitrunner.Runner
}

func NewApp(sess *session.Session) *App {
	return &App{
		sess: sess,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.sess != nil {
		a.runner = gitrunner.New(a.sess.Repo)
	}
}

func validateRefs(ctx context.Context, sess *session.Session) error {
	runner := gitrunner.New(sess.Repo)
	_, err := runner.Run(ctx, "merge-base", sess.Base, sess.Head)
	return err
}

type FileChangeDTO struct {
	OldPath string `json:"oldPath"`
	Path    string `json:"path"`
	Status  string `json:"status"`
}

func (a *App) ListChanges() ([]FileChangeDTO, error) {
	if a.sess == nil {
		return nil, fmt.Errorf("no session: use -C, --base, --head flags")
	}
	changes, err := changelist.ListChanges(a.ctx, a.runner, a.sess.Base, a.sess.Head)
	if err != nil {
		return nil, err
	}
	dtos := make([]FileChangeDTO, len(changes))
	for i, c := range changes {
		dtos[i] = FileChangeDTO{OldPath: c.OldPath, Path: c.NewPath, Status: c.Status}
	}
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].Path < dtos[j].Path
	})
	return dtos, nil
}

type LineDTO struct {
	Type    int    `json:"type"`    // 0=context, 1=added, 2=removed
	Content string `json:"content"`
	OldNum  int    `json:"oldNum"`
	NewNum  int    `json:"newNum"`
}

type HunkDTO struct {
	OldStart int       `json:"oldStart"`
	OldCount int       `json:"oldCount"`
	NewStart int       `json:"newStart"`
	NewCount int       `json:"newCount"`
	Header   string    `json:"header"`
	Lines    []LineDTO `json:"lines"`
}

type PatchDTO struct {
	OldPath     string    `json:"oldPath"`
	NewPath     string    `json:"newPath"`
	Hunks       []HunkDTO `json:"hunks"`
	IsBinary    bool      `json:"isBinary"`
	IsSubmodule bool      `json:"isSubmodule"`
}

func (a *App) LoadPatch(filePath string) (*PatchDTO, error) {
	if a.sess == nil {
		return nil, fmt.Errorf("no session: use -C, --base, --head flags")
	}
	df, err := patchloader.LoadPatch(a.ctx, a.runner, a.sess.Base, a.sess.Head, filePath)
	if err != nil {
		return nil, err
	}
	return diffFileToDTO(df), nil
}

func diffFileToDTO(df *diffparser.DiffFile) *PatchDTO {
	dto := &PatchDTO{
		OldPath:     df.OldPath,
		NewPath:     df.NewPath,
		IsBinary:    df.IsBinary,
		IsSubmodule: df.IsSubmodule,
	}
	for _, hunk := range df.Hunks {
		hunkDTO := HunkDTO{
			OldStart: hunk.OldStart,
			OldCount: hunk.OldCount,
			NewStart: hunk.NewStart,
			NewCount: hunk.NewCount,
			Header:   hunk.Header,
		}
		for _, line := range hunk.Lines {
			hunkDTO.Lines = append(hunkDTO.Lines, LineDTO{
				Type:    int(line.Type),
				Content: line.Content,
				OldNum:  line.OldNum,
				NewNum:  line.NewNum,
			})
		}
		dto.Hunks = append(dto.Hunks, hunkDTO)
	}
	return dto
}
