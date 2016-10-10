package libv2ray

import (
	"fmt"

	"v2ray.com/core"
)

/*CheckVersion int
This func will return libv2ray binding version.
*/
func CheckVersion() int {
	return 20
}

/*CheckVersionX string
This func will return libv2ray binding version and V2Ray version used.
*/
func CheckVersionX() string {
	return fmt.Sprintf("Libv2ray rev. %d, along with V2Ray %s", CheckVersion(), core.Version())
}
