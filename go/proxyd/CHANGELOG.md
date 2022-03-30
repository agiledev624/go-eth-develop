# @eth-optimism/proxyd

## 3.8.2

### Patch Changes

- ae18cea1: Don't hit Redis when the out of service interval is zero

## 3.8.1

### Patch Changes

- acf7dbd5: Update to go-ethereum v1.10.16

## 3.8.0

### Minor Changes

- 527448bb: Handle nil responses better

## 3.7.0

### Minor Changes

- 3c2926b1: Add debug cache status header to proxyd responses

## 3.6.0

### Minor Changes

- 096c5f20: proxyd: Allow cached RPCs to be evicted by redis
- 71d64834: Add caching for block-dependent RPCs
- fd2e1523: proxyd: Cache block-dependent RPCs
- 1760613c: Add integration tests and batching

## 3.5.0

### Minor Changes

- 025a3c0d: Add request/response payload size metrics to proxyd
- daf8db0b: cache immutable RPC responses in proxyd
- 8aa89bf3: Add X-Forwarded-For header when proxying RPCs on proxyd

## 3.4.1

### Patch Changes

- 415164e1: Force proxyd build

## 3.4.0

### Minor Changes

- 4b56ed84: Various proxyd fixes

## 3.3.0

### Minor Changes

- 7b7ffd2e: Allows string RPC ids on proxyd

## 3.2.0

### Minor Changes

- 73484138: Adds ability to specify env vars in config

## 3.1.2

### Patch Changes

- 1b79aa62: Release proxyd

## 3.1.1

### Patch Changes

- b8802054: Trigger release of proxyd
- 34fcb277: Bump proxyd to test release build workflow

## 3.1.0

### Minor Changes

- da6138fd: Updated metrics, support local rate limiter

### Patch Changes

- 6c7f483b: Add support for additional SSL certificates in Docker container

## 3.0.0

### Major Changes

- abe231bf: Make endpoints match Geth, better logging

## 2.0.0

### Major Changes

- 6c50098b: Update metrics, support WS
- f827dbda: Brings back the ability to selectively route RPC methods to backend groups

### Minor Changes

- 8cc824e5: Updates proxyd to include additional error metrics.
- 9ba4c5e0: Update metrics, support authenticated endpoints
- 78d0f3f0: Put special errors in a dedicated metric, pass along the content-type header

### Patch Changes

- 6e6a55b1: Canary release

## 1.0.2

### Patch Changes

- b9d2fbee: Trigger releases

## 1.0.1

### Patch Changes

- 893623c9: Trigger patch releases for dockerhub

## 1.0.0

### Major Changes

- 28aabc41: Initial release of RPC proxy daemon
