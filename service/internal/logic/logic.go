// Package logic provides the logic for converting a CSV file into a JSON tree.
package logic

import (
	"context"
	"fmt"

	"github.com/mnabbasabadi/relex_convertor/service/internal/converter"
	"github.com/mnabbasabadi/relex_convertor/service/internal/extractdata"
	"github.com/mnabbasabadi/relex_convertor/service/shared/domain"
	"github.com/rs/zerolog"
)

var (
	_ Logic = new(controller)
	// ErrorInvalidPayload ...
	ErrorInvalidPayload = fmt.Errorf("invalid payload")
)

type (
	// Logic ...
	Logic interface {
		ConvertToTreeNode(ctx context.Context, records [][]string) (domain.Node, error)
	}

	controller struct {
		logger zerolog.Logger
		extractdata.Extractor
		converter.Converter
	}
)

// New ...
func New(logger zerolog.Logger) Logic {
	return &controller{
		logger:    logger,
		Extractor: extractdata.New(),
		Converter: converter.New(),
	}
}

// ConvertToTreeNode ...
func (c controller) ConvertToTreeNode(ctx context.Context, records [][]string) (domain.Node, error) {
	dataRecords, err := c.Extract(ctx, records)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to extract data")
		return domain.Node{}, fmt.Errorf("failed to extract data: %w %w", err, ErrorInvalidPayload)
	}
	root, err := c.BuildTree(ctx, dataRecords)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to build tree")
		return domain.Node{}, err
	}

	return root, nil
}
