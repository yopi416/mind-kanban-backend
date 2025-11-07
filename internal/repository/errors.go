package repository

import "errors"

var ErrOptimisticLock = errors.New("optimistic lock conflict")
