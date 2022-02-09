# fair9

## fair9/Retry

Idiomatic retry library that respects a retry budget.

### Why

To understand what this is trying to achieve, read https://jacquesbernier.medium.com/zero-is-a-good-default-number-of-retries-abe431941994

### Features

* Process level retry budget
* Default to 0 retries
* Lock free implementation
* No background go routines
* Context support

### Get Started

https://github.com/jacquesbernier/fair9/blob/c3dac0bd88b0db61cd63a02911dab7913ec69c55/retry/example_test.go#L1-L6


### Other libraries

Most other libraries rely on simple attempt count and backoff. This is not helpful in case large degradation.

* [avast/retry-go](https://github.com/avast/retry-go)
* [giantswarm/retry-go](https://github.com/giantswarm/retry-go)
* [sethgrid/pester](https://github.com/sethgrid/pester)
* [cenkalti/backoff](https://github.com/cenkalti/backoff)
* [rafaeljesus/retry-go](https://github.com/rafaeljesus/retry-go)
* [matryer/try](https://github.com/matryer/try)

