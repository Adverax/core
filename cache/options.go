package cache

import (
	"time"
)

type features struct {
	resolved     bool
	expiration   bool
	prolongation bool
	capacity     bool
	size         bool
	daemon       bool
}

type Options[K comparable, V any] struct {
	feature[K, V]
	expiration     time.Duration
	daemonInterval time.Duration
	index          index[K, V]
	features       features
}

type Option[K comparable, V any] func(*Options[K, V])

func WithExpiration[K comparable, V any](expiration time.Duration) Option[K, V] {
	return func(options *Options[K, V]) {
		if !options.features.resolved {
			options.features.expiration = true
			return
		}

		options.expiration = expiration
		options.feature = &featureExpiration[K, V]{feature: options.feature}
	}
}

func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(options *Options[K, V]) {
		if !options.features.resolved {
			options.features.capacity = true
			return
		}

		options.feature = &featureCapacity[K, V]{
			feature:  options.feature,
			capacity: capacity,
			index:    options.index,
		}
	}
}

func WithProlongation[K comparable, V any]() Option[K, V] {
	return func(options *Options[K, V]) {
		if !options.features.resolved {
			options.features.prolongation = true
			return
		}

		if !options.features.expiration {
			options.feature = &featureExpiration[K, V]{feature: options.feature}
		}

		if options.features.expiration {
			options.feature = &featureExpirationProlongation[K, V]{
				feature:      options.feature,
				prolongation: options.expiration,
				index:        options.index,
			}
		} else if options.features.capacity {
			options.feature = &featureCapacityProlongation[K, V]{
				feature: options.feature,
				index:   options.index,
			}
		}
	}
}

func WithSize[K comparable, V any](maxSize int64, sizeOf func(item *Item[K, V]) int64) Option[K, V] {
	return func(options *Options[K, V]) {
		if !options.features.resolved {
			options.features.size = true
			return
		}

		options.feature = &featureSize[K, V]{
			feature: options.feature,
			index:   options.index,
			maxSize: maxSize,
			sizeOf:  sizeOf,
		}
	}
}

func WithDaemon[K comparable, V any](interval time.Duration) Option[K, V] {
	return func(options *Options[K, V]) {
		if !options.features.resolved {
			options.features.daemon = true
			return
		}

		options.daemonInterval = interval
	}
}
