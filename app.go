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

// FileChangeDTO is a JSON-friendly DTO for the frontend.
type FileChangeDTO struct {
	Path   string `json:"path"`
	Status string `json:"status"`
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
		dtos[i] = FileChangeDTO{Path: c.Path, Status: c.Status}
	}
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].Path < dtos[j].Path
	})
	return dtos, nil
}

// LineDTO is a JSON-friendly DTO for a single line in the diff.
type LineDTO struct {
	Type    int    `json:"type"`    // 0=context, 1=added, 2=removed
	Content string `json:"content"`
	OldNum  int    `json:"oldNum"`
	NewNum  int    `json:"newNum"`
}

// HunkDTO is a JSON-friendly DTO for a hunk.
type HunkDTO struct {
	OldStart int       `json:"oldStart"`
	OldCount int       `json:"oldCount"`
	NewStart int       `json:"newStart"`
	NewCount int       `json:"newCount"`
	Header   string    `json:"header"`
	Lines    []LineDTO `json:"lines"`
}

// PatchDTO is a JSON-friendly DTO for a parsed diff file.
type PatchDTO struct {
	OldPath string    `json:"oldPath"`
	NewPath string    `json:"newPath"`
	Hunks   []HunkDTO `json:"hunks"`
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
		OldPath: df.OldPath,
		NewPath: df.NewPath,
	}
	for _, h := range df.Hunks {
		hunkDTO := HunkDTO{
			OldStart: h.OldStart,
			OldCount: h.OldCount,
			NewStart: h.NewStart,
			NewCount: h.NewCount,
			Header:   h.Header,
		}
		for _, l := range h.Lines {
			hunkDTO.Lines = append(hunkDTO.Lines, LineDTO{
				Type:    int(l.Type),
				Content: l.Content,
				OldNum:  l.OldNum,
				NewNum:  l.NewNum,
			})
		}
		dto.Hunks = append(dto.Hunks, hunkDTO)
	}
	return dto
}
