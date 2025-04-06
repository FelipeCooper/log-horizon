package mlog

// OrderDirection define a direção da ordenação
type OrderDirection string

const (
	// OrderAsc representa ordenação ascendente
	OrderAsc OrderDirection = "asc"

	// OrderDesc representa ordenação descendente
	OrderDesc OrderDirection = "desc"
)

// OrderField define os campos disponíveis para ordenação
type OrderField string

const (
	// OrderByTimestamp ordena por timestamp
	OrderByTimestamp OrderField = "timestamp"

	// OrderByLevel ordena por nível
	OrderByLevel OrderField = "level"
)

// OrderOption define uma opção de ordenação para consulta de logs
type OrderOption func(*OrderOptions)

// OrderOptions contém todas as opções aplicáveis para ordenar logs
type OrderOptions struct {
	// Field é o campo pelo qual ordenar
	Field OrderField

	// Direction é a direção da ordenação
	Direction OrderDirection
}

// DefaultOrderOptions retorna as opções de ordenação padrão
func DefaultOrderOptions() OrderOptions {
	return OrderOptions{
		Field:     OrderByTimestamp,
		Direction: OrderDesc,
	}
}

// WithOrderField define o campo para a ordenação
func WithOrderField(field OrderField) OrderOption {
	return func(o *OrderOptions) {
		o.Field = field
	}
}

// WithOrderDirection define a direção da ordenação
func WithOrderDirection(direction OrderDirection) OrderOption {
	return func(o *OrderOptions) {
		o.Direction = direction
	}
}

// NewOrderOptions cria OrderOptions com as opções fornecidas
func NewOrderOptions(opts ...OrderOption) OrderOptions {
	options := DefaultOrderOptions()
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
