package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
	"go.uber.org/zap"
)

type (
	CodeService interface {
		GenerateTransactionCode(ctx context.Context, transactionType model.EntityType, teamId uint) (string, error)
		GenerateProductCode(ctx context.Context, teamId uint) (string, error)
		GenerateLotCode(ctx context.Context, productId uint) (string, error)
	}

	codeService struct {
		logger             logger.AppLogger
		sequenceRepository repository.AutoIncrementSequence
	}
)

// GenerateTransactionCode generates a unique transaction code
func (c *codeService) GenerateTransactionCode(ctx context.Context, transactionType model.EntityType, teamId uint) (string, error) {

	sequence, err := c.sequenceRepository.GetNextSequence(ctx, nil, transactionType, teamId, 0)
	if err != nil {
		c.logger.Error("Failed to get next sequence for transaction code", zap.Error(err))
		return "", err
	}
	now := time.Now()
	dateStr := now.Format("060102")

	typePrefix := strings.ToUpper(string(transactionType))

	code := fmt.Sprintf("%s-%s%01d%03d", typePrefix, dateStr, teamId, sequence)

	c.logger.Info("Generated transaction code", zap.String("code", code), zap.Uint("sequence", sequence))

	return code, nil
}

// GenerateProductCode generates a unique product code
func (c *codeService) GenerateProductCode(ctx context.Context, teamId uint) (string, error) {
	sequence, err := c.sequenceRepository.GetNextSequence(ctx, nil, model.EntityTypeProduct, teamId, 0)
	if err != nil {
		c.logger.Error("Failed to get next sequence for product code", zap.Error(err))
		return "", err
	}

	code := fmt.Sprintf("PN-%02d%04d", teamId, sequence)
	c.logger.Info("Generated product code", zap.String("code", code), zap.Uint("sequence", sequence))

	return code, nil
}

// GenerateLotCode generates a unique lot code for a specific product
func (c *codeService) GenerateLotCode(ctx context.Context, productId uint) (string, error) {
	sequence, err := c.sequenceRepository.GetNextSequence(ctx, nil, model.EntityTypeLot, 0, productId)
	if err != nil {
		c.logger.Error("Failed to get next sequence for lot code", zap.Error(err))
		return "", err
	}

	code := fmt.Sprintf("LN-%02d%04d", productId, sequence)
	c.logger.Info("Generated lot code", zap.String("code", code), zap.Uint("sequence", sequence))

	return code, nil
}

func NewCodeService(
	logger logger.AppLogger,
	sequenceRepository repository.AutoIncrementSequence,
) CodeService {
	return &codeService{
		logger:             logger,
		sequenceRepository: sequenceRepository,
	}
}
