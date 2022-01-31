// Copyright 2022 Dhi Aurrahman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ratelimit

import (
	"fmt"
	"os"

	"github.com/envoyproxy/ratelimit/src/settings"

	settingsv1 "github.com/dio/rundown/generated/ratelimit/settings/v1"
)

func NewSettings(s *settingsv1.Settings) settings.Settings { //nolint:gocyclo
	c := settings.NewSettings()

	// I know this looks tedious, we could have a generator somehow later, but this works.

	// Server listen address config
	if s.Host != nil {
		c.Host = s.Host.Value
	}
	if s.Port != nil {
		c.Port = int(s.Port.Value)
	}
	if s.GrpcHost != nil {
		c.GrpcHost = s.GrpcHost.Value
	}
	if s.GrpcPort != nil {
		c.GrpcPort = int(s.GrpcPort.Value)
	}
	if s.DebugHost != nil {
		c.DebugHost = s.DebugHost.Value
	}

	// gRPC server settings.
	// GrpcMaxConnectionAge is a duration for the maximum amount of time a connection may exist before
	// it will be closed by sending a GoAway. A random jitter of +/-10% will be added to
	// MaxConnectionAge to spread out connection storms.
	if s.GrpcMaxConnectionAge != nil {
		c.GrpcMaxConnectionAge = s.GrpcMaxConnectionAge.AsDuration()
	}
	// GrpcMaxConnectionAgeGrace is an additive period after MaxConnectionAge after which the
	// connection will be forcibly closed.
	if s.GrpcMaxConnectionAgeGrace != nil {
		c.GrpcMaxConnectionAgeGrace = s.GrpcMaxConnectionAgeGrace.AsDuration()
	}

	// Logging settings.
	if s.LogLevel != nil {
		c.LogLevel = s.LogLevel.Value // TODO(dio): Enumerate this.
	}
	if s.LogFormat != nil {
		c.LogFormat = s.LogFormat.Value // TODO(dio): Enumerate this.
	}

	// Statsd-related settings. The github.com/envoyproxy/ratelimit uses github.com/lyft/gostats for
	// setting up statsd, the library expects environment USE_STATSD and STATSD_* variables, namely:
	//   - USE_STATSD, defaults to true.
	//   - STATSD_HOST, defaults to "localhost".
	//   - STATSD_PORT, defaults to 8125.
	// Reference: https://github.com/lyft/gostats/blob/aa005a717918424bef89063bdac8772d976782f8/settings.go#L28-L34.
	//
	// Note: STATSD_PROTOCOL, defaults to "tcp", is not exposed by github.com/envoyproxy/ratelimit.
	// See https://github.com/envoyproxy/ratelimit/blob/8d6488ead8618ce49a492858321dae946f2d97bc/src/settings/settings.go#L34-L36.
	//
	// Hence, we also do setenv here.
	if s.UseStatsd != nil {
		c.UseStatsd = s.UseStatsd.Value
		_ = os.Setenv("USE_STATSD", fmt.Sprintf("%t", c.UseStatsd))
	}
	if s.StatsdHost != nil {
		c.StatsdHost = s.StatsdHost.Value
		_ = os.Setenv("STATSD_HOST", c.StatsdHost)
	}
	if s.StatsdPort != nil {
		c.StatsdPort = int(s.StatsdPort.Value)
		_ = os.Setenv("STATSD_PORT", fmt.Sprintf("%d", c.StatsdPort))
	}
	// While this ExtraTags field is also considered as a setting for statsd, this one is not being
	// passed to github.com/lyft/gostats through environment variables.
	if s.ExtraTags != nil {
		c.ExtraTags = s.ExtraTags
	}

	// Settings for rate limit configuration.
	if s.RuntimePath != nil {
		c.RuntimePath = s.RuntimePath.Value
	}
	if s.RuntimeSubdirectory != nil {
		c.RuntimeSubdirectory = s.RuntimeSubdirectory.Value
	}
	if s.RuntimeIgnoreDotFiles != nil {
		c.RuntimeIgnoreDotFiles = s.RuntimeIgnoreDotFiles.Value
	}
	if s.RuntimeWatchRoot != nil {
		c.RuntimeWatchRoot = s.RuntimeWatchRoot.Value
	}

	// Settings for all cache types.
	if s.ExpirationJitterMaxSeconds != nil {
		c.ExpirationJitterMaxSeconds = s.ExpirationJitterMaxSeconds.Value
	}
	if s.LocalCacheSizeInBytes != nil {
		c.LocalCacheSizeInBytes = int(s.LocalCacheSizeInBytes.Value)
	}
	if s.NearLimitRatio != nil {
		c.NearLimitRatio = s.NearLimitRatio.Value
	}
	if s.CacheKeyPrefix != nil {
		c.CacheKeyPrefix = s.CacheKeyPrefix.Value
	}
	if s.BackendType != nil {
		c.BackendType = s.BackendType.Value // TODO(dio): Enumerate this.
	}

	// Settings for optional returning of custom headers.
	if s.RateLimitResponseHeadersEnabled != nil {
		c.RateLimitResponseHeadersEnabled = s.RateLimitResponseHeadersEnabled.Value
	}
	if s.HeaderRatelimitLimit != nil {
		c.HeaderRatelimitLimit = s.HeaderRatelimitLimit.Value
	}
	if s.HeaderRatelimitRemaining != nil {
		c.HeaderRatelimitRemaining = s.HeaderRatelimitRemaining.Value
	}
	if s.HeaderRatelimitReset != nil {
		c.HeaderRatelimitLimit = s.HeaderRatelimitLimit.Value
	}

	// Redis settings.
	if s.RedisSocketType != nil {
		c.RedisSocketType = s.RedisSocketType.Value // TODO(dio): Enumerate this.
	}
	if s.RedisType != nil {
		c.RedisType = s.RedisType.Value
	}
	if s.RedisUrl != nil {
		c.RedisUrl = s.RedisUrl.Value
	}
	if s.RedisPoolSize != nil {
		c.RedisPoolSize = int(s.RedisPoolSize.Value)
	}
	if s.RedisAuth != nil {
		c.RedisAuth = s.RedisAuth.Value
	}
	if s.RedisTls != nil {
		c.RedisTls = s.RedisTls.Value
	}

	// RedisPipelineWindow sets the duration after which internal pipelines will be flushed.
	// If window is zero then implicit pipelining will be disabled. Radix use 150us for the
	// default value, see https://github.com/mediocregopher/radix/blob/v3.5.1/pool.go#L278.
	if s.RedisPipelineWindow != nil {
		c.RedisPipelineWindow = s.RedisPipelineWindow.AsDuration()
	}
	// RedisPipelineLimit sets maximum number of commands that can be pipelined before flushing.
	// If limit is zero then no limit will be used and pipelines will only be limited by the specified time window.
	if s.RedisPipelineLimit != nil {
		c.RedisPipelineLimit = int(s.RedisPipelineLimit.Value)
	}
	if s.RedisPerSecond != nil {
		c.RedisPerSecond = s.RedisPerSecond.Value
	}
	if s.RedisPerSecondSocketType != nil {
		c.RedisPerSecondSocketType = s.RedisPerSecondSocketType.Value // TODO(dio): Enumerate this.
	}
	if s.RedisPerSecondType != nil {
		c.RedisPerSecondType = s.RedisPerSecondType.Value // TODO(dio): Enumerate this.
	}
	if s.RedisPerSecondUrl != nil {
		c.RedisPerSecondUrl = s.RedisPerSecondUrl.Value
	}
	if s.RedisPerSecondPoolSize != nil {
		c.RedisPerSecondPoolSize = int(s.RedisPerSecondPoolSize.Value)
	}
	if s.RedisPerSecondAuth != nil {
		c.RedisPerSecondAuth = s.RedisPerSecondAuth.Value
	}
	if s.RedisPerSecondTls != nil {
		c.RedisPerSecondTls = s.RedisPerSecondTls.Value
	}
	// RedisPerSecondPipelineWindow sets the duration after which internal pipelines will be flushed for per second redis.
	// See comments of RedisPipelineWindow for details.
	if s.RedisPerSecondPipelineWindow != nil {
		c.RedisPerSecondPipelineWindow = s.RedisPerSecondPipelineWindow.AsDuration()
	}
	if s.RedisPerSecondPipelineLimit != nil {
		c.RedisPerSecondPipelineLimit = int(s.RedisPerSecondPipelineLimit.Value)
	}
	if s.RedisHealthCheckActiveConnection != nil {
		c.RedisHealthCheckActiveConnection = s.RedisHealthCheckActiveConnection.Value
	}

	// Memcache settings.
	if s.MemcacheHostPort != nil {
		c.MemcacheHostPort = s.MemcacheHostPort
	}
	// MemcacheMaxIdleConns sets the maximum number of idle TCP connections per memcached node.
	// The default is 2 as that is the default of the underlying library. This is the maximum
	// number of connections to memcache kept idle in pool, if a connection is needed but none
	// are idle a new connection is opened, used and closed and can be left in a time-wait state
	// which can result in high CPU usage.
	if s.MemcacheMaxIdleConns != nil {
		c.MemcacheMaxIdleConns = int(s.MemcacheMaxIdleConns.Value)
	}
	if s.MemcacheSrv != nil {
		c.MemcacheSrv = s.MemcacheSrv.Value
	}
	if s.MemcacheSrvRefresh != nil {
		c.MemcacheSrvRefresh = s.MemcacheSrvRefresh.AsDuration()
	}

	// Should the ratelimiting be running in Global shadow-mode, ie. never report a ratelimit status,
	// unless a rate was provided from envoy as an override.
	if s.GlobalShadownMode != nil {
		c.GlobalShadowMode = s.GlobalShadownMode.Value
	}
	return c
}

// Merge github.com/envoyproxy/ratelimit/src/settings with the one we have from YAML-based loading.
