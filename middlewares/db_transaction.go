package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const txKey = "txKey"

// TransactionMiddleware is a custom middleware to handle transactions
func TransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		c.Set(txKey, tx)
		defer func() {
			if c.Writer.Status() >= 400 {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
		c.Next()
	}
}
