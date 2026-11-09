package errs

import "errors"

// ErrNotFound означает, что сущность не найдена. Репозитории оборачивают её
// через %w, контроллеры проверяют через errors.Is и выбирают код 404.
var ErrNotFound = errors.New("not found")
