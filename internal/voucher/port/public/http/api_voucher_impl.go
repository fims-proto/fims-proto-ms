//  implementation of generated openapi
package http
import(
	"github/fims-proto/fims-proto-ms/internal/voucher/app"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) _AllVouchers(c *gin.Context){
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var results []VoucherQry
	for _, v := range vouchers {
		result, err := VoucherQry{
			
		}
	} 
}

func (h Handler) _Audit(c *gin.Context){
	return
}

func (h Handler) _Review(c *gin.Context){
	return
}

func (h Handler) _Update(c *gin.Context){
	return
}

func (h Handler) _Record(c *gin.Context){
	return
}

func (h Handler) _VoucherForUUID(c *gin.Context){
	return
}