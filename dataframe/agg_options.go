package dataframe

// Interpolation selects how Median/Quantile interpolate between data points
// when the target rank falls between two values (numpy/pandas method set).
type Interpolation uint8

const (
	Linear Interpolation = iota // default: i + (j-i)*frac
	Lower                       // value at the lower rank
	Higher                      // value at the higher rank
	Nearest                     // nearest rank
	Midpoint                    // (i + j) / 2
)

type aggConfig struct {
	ddof      int
	ddofSet   bool
	interp    Interpolation
	interpSet bool
}

// AggOption configures a per-aggregator parameter.
type AggOption func(*aggConfig)

// WithDDoF sets the delta degrees of freedom for Std/Variance (divisor n-ddof).
func WithDDoF(d int) AggOption { return func(c *aggConfig) { c.ddof = d; c.ddofSet = true } }

// WithInterpolation sets the interpolation method for Median/Quantile.
func WithInterpolation(m Interpolation) AggOption {
	return func(c *aggConfig) { c.interp = m; c.interpSet = true }
}

func newAggConfig(opts []AggOption) aggConfig {
	c := aggConfig{ddof: 0, interp: Linear}
	for _, o := range opts {
		o(&c)
	}
	return c
}
