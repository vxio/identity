// stub to get pkger to work
package identity

import (
	"github.com/markbates/pkger"
)

// Add in all includes that pkger should embed into the application here
var _ = pkger.Include("/configs/config.default.yml")
