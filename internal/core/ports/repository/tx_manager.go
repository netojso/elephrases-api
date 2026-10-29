package portrepository

import "context"

// TxManager executa fn dentro de uma transação. O ctx recebido por fn
// carrega a transação; repositórios chamados com ele participam dela.
// Se fn retornar erro, a transação é desfeita.
type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
