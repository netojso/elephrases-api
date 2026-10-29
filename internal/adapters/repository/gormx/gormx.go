// Package gormx reúne utilitários compartilhados pelos repositórios GORM.
package gormx

import (
	"errors"

	"github.com/netojso/elephrases-api/pkg"
	"gorm.io/gorm"
)

// MapError traduz erros do GORM para erros de domínio.
func MapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pkg.ErrNotFound
	}
	return err
}

// DeleteResult retorna pkg.ErrNotFound quando nenhuma linha foi afetada.
func DeleteResult(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return pkg.ErrNotFound
	}
	return nil
}
