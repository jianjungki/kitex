/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ttstream

// ClientHandlerOption define client handler options
type ClientHandlerOption func(cp *clientTransHandler)

// Deprecated: use ClientHandlerOption instead
type ClientProviderOption = ClientHandlerOption

// WithClientMetaFrameHandler register TTHeader Streaming meta frame handler
func WithClientMetaFrameHandler(handler MetaFrameHandler) ClientHandlerOption {
	return func(cp *clientTransHandler) {
		cp.metaHandler = handler
	}
}

// WithClientHeaderFrameHandler register TTHeader Streaming header frame handler
func WithClientHeaderFrameHandler(handler HeaderFrameWriteHandler) ClientHandlerOption {
	return func(cp *clientTransHandler) {
		cp.headerHandler = handler
	}
}

// WithClientLongConnPool using long connection pool for client
func WithClientLongConnPool(config LongConnConfig) ClientHandlerOption {
	return func(cp *clientTransHandler) {
		cp.transPool = newLongConnTransPool(config)
	}
}

// WithClientShortConnPool using short connection pool for client
func WithClientShortConnPool() ClientHandlerOption {
	return func(cp *clientTransHandler) {
		cp.transPool = newShortConnTransPool()
	}
}

// WithClientMuxConnPool using mux connection pool for client
func WithClientMuxConnPool(config MuxConnConfig) ClientHandlerOption {
	return func(cp *clientTransHandler) {
		cp.transPool = newMuxConnTransPool(config)
	}
}
